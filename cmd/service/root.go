package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xn3cr0nx/email-service/internal/environment"
	"github.com/xn3cr0nx/email-service/internal/mailer/postmark"
	"github.com/xn3cr0nx/email-service/internal/server"
	"github.com/xn3cr0nx/email-service/internal/task"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"github.com/xn3cr0nx/email-service/pkg/meter"
	"github.com/xn3cr0nx/email-service/pkg/tracer"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sys/unix"
)

var rootCmd = &cobra.Command{
	Use:   "mailer",
	Short: "Email service",
	Long:  `Postmark based email service offering both event based execution and REST API`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Setup()
	},
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("mailer")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if value, ok := os.LookupEnv("CONFIG_FILE"); ok {
		viper.SetConfigFile(value)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/mailer/")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	viper.ReadInConfig()
	f := viper.ConfigFileUsed()
	if f != "" {
		fmt.Printf("Found configuration file: %s \n", f)
	}

	env := environment.New()
	if err := parseFlags(env); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	logger.Info("Email Service", "Starting", logger.Params{"timestamp": time.Now()})

	env := environment.Get()

	if err := validateConfig(env); err != nil {
		logger.Error("Email Service", fmt.Errorf("Configuration error: %w", err), logger.Params{})
		os.Exit(-1)
	}

	ctx := context.Background()

	var tr trace.Tracer
	if env.OtelExporterJaegerEnable {
		var err error
		tr, err = tracer.NewTracerProvider(ctx, env.ServiceName)
		if err != nil {
			logger.Error("Email Service", fmt.Errorf("Cannot initialize tracer: %w", err), logger.Params{})
			os.Exit(-1)
		}
	}

	var mt metric.Meter
	if env.OtelExporterPrometheusEnable {
		var err error
		mt, err = meter.NewMeter(&meter.Config{Name: env.ServiceName, Port: env.OtelExporterPrometheusPort})
		if err != nil {
			logger.Error("Email Service", fmt.Errorf("Cannot initialize tracer: %w", err), logger.Params{})
			os.Exit(-1)
		}
	}

	// initialize template cache
	templateDir := env.TemplateDir
	if templateDir == "" {
		templateDir = "templates/"
	}
	_, err := template.NewTemplateCache(&templateDir)
	if err != nil {
		logger.Error("Email Service", fmt.Errorf("Cannot initialize template cache: %w", err), logger.Params{})
		os.Exit(-1)
	}

	mailer := postmark.NewClient(viper.GetString("postmark.server"), viper.GetString("postmark.account"))

	if env.Rest {
		s := server.NewServer(env.Port, mailer, tr, mt)
		go s.Listen()
	}

	if env.Backend == "asynq" {
		redisAddress := fmt.Sprintf("%s:%d", env.RedisHost, env.RedisPort)
		server := asynq.NewServer(
			asynq.RedisClientOpt{
				Addr:     redisAddress,
				Password: env.RedisPassword,
				DB:       env.RedisDB,
			},
			asynq.Config{
				Concurrency: env.Concurrency,
				Queues: map[string]int{
					string(env.Queue): 10,
				},
			})
		defer server.Stop()

		h := task.NewEmailHandler(mailer, tr, mt)
		if err := server.Run(h); err != nil {
			logger.Error("Email Service", fmt.Errorf("cannot start queue server %v", err), logger.Params{})
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
		<-sigs // wait for termination signal
		server.Stop()
	} else if env.Backend == "nats" {
		natsAddress := fmt.Sprintf("nats://%s:%d", env.NatsHost, env.NatsPort)
		nc, err := nats.Connect(natsAddress)
		if err != nil {
			logger.Error("Email Service", fmt.Errorf("cannot connect NATS broker: %v", err), logger.Params{})
			os.Exit(-1)
		}

		subscriber := task.NewNatsEmailConsumer(nc, mailer, tr, mt)
		// blocking consumer reading messages
		if err := subscriber.Run(ctx); err != nil {
			logger.Error("Email service", err, logger.Params{})
			os.Exit(-1)
		}
	} else if env.Backend == "kafka" {
		// initialize a new reader with the brokers and topic
		// the groupID identifies the consumer and prevents
		// it from receiving duplicate messages
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     env.KafkaAddresses,
			Topic:       env.KafkaTopic,
			StartOffset: kafka.FirstOffset,
			GroupID:     env.KafkaGroup,
			Logger:      logger.Log,
		})
		consumer := task.NewKafkaEmailConsumer(r, mailer, tr, mt)
		defer r.Close()
		// blocking consumer reading messages
		if err := consumer.Run(ctx); err != nil {
			logger.Error("Email service", err, logger.Params{})
			os.Exit(-1)
		}
	} else {
		logger.Error("Email service", errors.New("not valid backend configured"), logger.Params{"env": env.Backend})
		os.Exit(-1)
	}
}

var (
	errMissingPort            = errors.New("missing server port")
	errMissingRedisAddress    = errors.New("missing redis address")
	errRedisDbOutOfRange      = errors.New("redis db out of range. allowed range 0-15")
	errQueueConcurrencyNotSet = errors.New("queue concurrent workers number not defined. min 1")
	errMissingNatsAddress     = errors.New("missing nats address")
)

func validateConfig(env *environment.Env) error {
	if env.Rest {
		if env.Port == 0 {
			return errMissingPort
		}
	}

	if env.Backend == "asynq" {
		if env.RedisHost == "" || env.RedisPort == 0 {
			return errMissingRedisAddress
		}
		if env.RedisDB > 15 {
			return errRedisDbOutOfRange
		}
		if env.Concurrency == 0 {
			return errQueueConcurrencyNotSet
		}
	}

	if env.Backend == "nats" {
		if env.NatsHost == "" || env.NatsPort == 0 {
			return errMissingNatsAddress
		}
	}

	return nil
}
