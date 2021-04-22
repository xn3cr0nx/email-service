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
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var (
	port, redisDb, concurrency, otel_exporter_jaeger_agent_port, otel_exporter_prometheus_port              int
	debug, rest, k, otel_exporter_jaeger_enable, otel_exporter_prometheus_enable                            bool
	name, templateDir, redisAddress, redisPassword, kafkaGroup, kafkaTopic, otel_exporter_jaeger_agent_host string
	kafkaAddress                                                                                            []string
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

	// Adds root flags and persistent flags
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Sets logging level to Debug")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 6066, "Bind http server to port")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "Mailer", "Set service name")
	rootCmd.PersistentFlags().BoolVar(&rest, "rest", false, "Enable exposed REST API to interact with the service")
	rootCmd.PersistentFlags().StringVar(&templateDir, "template_dir", "templates/", "Define templates folder path")
	rootCmd.PersistentFlags().BoolVar(&k, "kafka", false, "Set kafka as broker backend")
	rootCmd.PersistentFlags().StringArrayVar(&kafkaAddress, "kafka_address", []string{"localhost:9092"}, "Set kafka addresses")
	rootCmd.PersistentFlags().StringVar(&kafkaGroup, "kafka_group", "my-group", "Set kafka group name")
	rootCmd.PersistentFlags().StringVar(&kafkaTopic, "kafka_topic", "emails", "Set kafka partition name")
	rootCmd.PersistentFlags().StringVar(&redisAddress, "redis_address", "localhost:6379", "Set host address for redis backend")
	rootCmd.PersistentFlags().StringVar(&redisPassword, "redis_password", "", "Set password for redis backend")
	rootCmd.PersistentFlags().IntVar(&redisDb, "redis_db", 2, "Set redis database number")
	rootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 10, "Set number of concurrent workers for redis backend")
	rootCmd.PersistentFlags().BoolVar(&otel_exporter_jaeger_enable, "otel_exporter_jaeger_enable", false, "Enable OpenTelemetry based jager tracing")
	rootCmd.PersistentFlags().StringVar(&otel_exporter_jaeger_agent_host, "otel_exporter_jaeger_agent_host", "jaeger", "Override Jaeger agent hostname")
	rootCmd.PersistentFlags().IntVar(&otel_exporter_jaeger_agent_port, "otel_exporter_jaeger_agent_port", 14268, "Override Jaeger agent port")
	rootCmd.PersistentFlags().BoolVar(&otel_exporter_prometheus_enable, "otel_exporter_prometheus_enable", false, "Enable OpenTelemetry based prometheus metrics")
	rootCmd.PersistentFlags().IntVar(&otel_exporter_prometheus_port, "otel_exporter_prometheus_port", 14268, "Override Prometheus exposed port")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("debug", false)
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("rest", false)
	viper.BindPFlag("rest", rootCmd.PersistentFlags().Lookup("rest"))
	viper.SetDefault("port", 6066)
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("name", "Mailer")
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
	viper.SetDefault("template_dir", "templates/")
	viper.BindPFlag("template_dir", rootCmd.PersistentFlags().Lookup("template_dir"))
	viper.SetDefault("kafka", false)
	viper.BindPFlag("kafka", rootCmd.PersistentFlags().Lookup("kafka"))
	viper.SetDefault("kafka_address", []string{"localhost:9092"})
	viper.BindPFlag("kafka_address", rootCmd.PersistentFlags().Lookup("kafka_address"))
	viper.SetDefault("kafka_group", "my-group")
	viper.BindPFlag("kafka_group", rootCmd.PersistentFlags().Lookup("kafka_group"))
	viper.SetDefault("kafka_topic", "emails")
	viper.BindPFlag("kafka_topic", rootCmd.PersistentFlags().Lookup("kafka_topic"))
	viper.SetDefault("redis_address", "localhost:6379")
	viper.BindPFlag("redis_address", rootCmd.PersistentFlags().Lookup("redis_address"))
	viper.SetDefault("redis_password", "")
	viper.BindPFlag("redis_password", rootCmd.PersistentFlags().Lookup("redis_password"))
	viper.SetDefault("redis_db", 2)
	viper.BindPFlag("redis_db", rootCmd.PersistentFlags().Lookup("redis_db"))
	viper.SetDefault("concurrency", 10)
	viper.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))
	viper.SetDefault("otel_exporter_jaeger_enable", false)
	viper.BindPFlag("otel_exporter_jaeger_enable", rootCmd.PersistentFlags().Lookup("otel_exporter_jaeger_enable"))
	viper.SetDefault("otel_exporter_jaeger_agent_host", "jaeger")
	viper.BindPFlag("otel_exporter_jaeger_agent_host", rootCmd.PersistentFlags().Lookup("otel_exporter_jaeger_agent_host"))
	viper.SetDefault("otel_exporter_jaeger_agent_port", 14268)
	viper.BindPFlag("otel_exporter_jaeger_agent_port", rootCmd.PersistentFlags().Lookup("otel_exporter_jaeger_agent_port"))
	viper.SetDefault("otel_exporter_prometheus_enable", false)
	viper.BindPFlag("otel_exporter_prometheus_enable", rootCmd.PersistentFlags().Lookup("otel_exporter_prometheus_enable"))
	viper.SetDefault("otel_exporter_prometheus_port", 9464)
	viper.BindPFlag("otel_exporter_prometheus_port", rootCmd.PersistentFlags().Lookup("otel_exporter_prometheus_port"))

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
}

func run(cmd *cobra.Command, args []string) {
	logger.Info("Email Service", "Starting", logger.Params{"timestamp": time.Now()})

	if err := validateConfig(); err != nil {
		logger.Error("Email Service", fmt.Errorf("Configuration error: %w", err), logger.Params{})
		os.Exit(-1)
	}

	ctx := context.Background()

	var tr *trace.Tracer
	if viper.GetBool("otel_exporter_jaeger_enable") {
		var err error
		tr, err = tracer.NewTracer(&tracer.Config{
			Name:     viper.GetString("name"),
			Host:     viper.GetString("otel_exporter_jaeger_agent_host"),
			Port:     viper.GetInt("otel_exporter_jaeger_agent_port"),
			Exporter: tracer.Jaeger,
		})
		if err != nil {
			logger.Error("Email Service", fmt.Errorf("Cannot initialize tracer: %w", err), logger.Params{})
			os.Exit(-1)
		}
	}

	var mt *metric.Meter
	if viper.GetBool("otel_exporter_prometheus_enable") {
		var err error
		mt, err = meter.NewMeter(&meter.Config{Name: viper.GetString("name"), Port: viper.GetInt("otel_exporter_prometheus_port")})
		if err != nil {
			logger.Error("Email Service", fmt.Errorf("Cannot initialize tracer: %w", err), logger.Params{})
			os.Exit(-1)
		}
	}

	// initialize template cache
	templateDir := viper.GetString("template_dir")
	if templateDir == "" {
		templateDir = "templates/"
	}
	_, err := template.NewTemplateCache(&templateDir)
	if err != nil {
		logger.Error("Email Service", fmt.Errorf("Cannot initialize template cache: %w", err), logger.Params{})
		os.Exit(-1)
	}

	mailer := postmark.NewClient(viper.GetString("postmark.server"), viper.GetString("postmark.account"))

	if viper.GetBool("rest") {
		s := server.NewServer(viper.GetInt("port"), mailer, tr, mt)
		go s.Listen()
	}

	if viper.GetBool("kafka") {
		// initialize a new reader with the brokers and topic
		// the groupID identifies the consumer and prevents
		// it from receiving duplicate messages
		r := kafka.NewReader(kafka.ReaderConfig{
			// Brokers:     []string{"localhost:19092", "localhost:29092", "localhost:39092"},
			Brokers:     viper.GetStringSlice("kafka_address"),
			Topic:       viper.GetString("kafka_topic"),
			StartOffset: kafka.FirstOffset,
			GroupID:     viper.GetString("kafka_group"),
			Logger:      logger.Log,
		})
		consumer := task.NewKafkaEmailConsumer(r, mailer, tr, mt)
		// blocking consumer reading messages
		consumer.Run(ctx)

	} else {
		server := asynq.NewServer(
			asynq.RedisClientOpt{
				Addr:     viper.GetString("address"),
				Password: viper.GetString("password"),
				DB:       viper.GetInt("db"),
			},
			asynq.Config{
				Concurrency: viper.GetInt("concurrency"),
			})

		h := task.NewEmailHandler(mailer, tr, mt)
		if err := server.Run(h); err != nil {
			logger.Error("Email Service", fmt.Errorf("Cannot start queue server %v", err), logger.Params{})
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
		<-sigs // wait for termination signal
		server.Stop()
	}
}

var (
	errMissingPort            = errors.New("missing server port")
	errMissingRedisAddress    = errors.New("missing redis address")
	errRedisDbOutOfRange      = errors.New("redis db out of range. allowed range 0-15")
	errQueueConcurrencyNotSet = errors.New("queue concurrent workers number not defined. min 1")
)

func validateConfig() error {
	if viper.GetBool("rest") {
		if viper.GetInt("port") == 0 {
			return errMissingPort
		}
	}

	if !viper.GetBool("kafka") {
		if viper.GetString("address") == "" {
			return errMissingRedisAddress
		}
		if viper.GetInt("db") > 15 {
			return errRedisDbOutOfRange
		}
		if viper.GetInt("concurrency") == 0 {
			return errQueueConcurrencyNotSet
		}
	}

	return nil
}
