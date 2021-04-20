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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xn3cr0nx/email-service/internal/mailer/postmark"
	"github.com/xn3cr0nx/email-service/internal/server"
	"github.com/xn3cr0nx/email-service/internal/task"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"golang.org/x/sys/unix"
)

var (
	port, db, concurrency          int
	debug, rest, kafka             bool
	templateDir, address, password string
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
	rootCmd.PersistentFlags().BoolVar(&rest, "rest", false, "Enable exposed REST API to interact with the service")
	rootCmd.PersistentFlags().StringVar(&templateDir, "template_dir", "templates/", "Define templates folder path")
	rootCmd.PersistentFlags().BoolVar(&kafka, "kafka", false, "Set broken mode using kafka")
	rootCmd.PersistentFlags().StringVar(&address, "address", "localhost:6379", "Set host address for used backend")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Set password for used backend")
	rootCmd.PersistentFlags().IntVar(&db, "db", 2, "Set redis database number")
	rootCmd.PersistentFlags().IntVar(&concurrency, "concurrency", 10, "Set number of concurrent workers")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("debug", false)
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("rest", false)
	viper.BindPFlag("rest", rootCmd.PersistentFlags().Lookup("rest"))
	viper.SetDefault("port", 6066)
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("template_dir", "templates/")
	viper.BindPFlag("template_dir", rootCmd.PersistentFlags().Lookup("template_dir"))
	viper.SetDefault("kafka", false)
	viper.BindPFlag("kafka", rootCmd.PersistentFlags().Lookup("kafka"))
	viper.SetDefault("address", "localhost:6379")
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	viper.SetDefault("password", "")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.SetDefault("db", 2)
	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
	viper.SetDefault("concurrency", 10)
	viper.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))

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
		s := server.NewServer(viper.GetInt("port"), mailer)
		go s.Listen()
	}

	if viper.GetBool("kafka") {
		// TODO: dispatch Kafka consumer

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

		h := task.WelcomeEmailHandler{Mailer: mailer}
		mux := asynq.NewServeMux()

		// Define custom middleware to log processing time and log catched error
		mux.Use(func(h asynq.Handler) asynq.Handler {
			return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
				start := time.Now()
				logger.Info("Email Service Queue", fmt.Sprintf("Start processing"), logger.Params{"type": t.Type})
				if err := h.ProcessTask(ctx, t); err != nil {
					logger.Error("Email Service Queue", err, logger.Params{"type": t.Type})
					return err
				}
				logger.Info("Email Service Queue", fmt.Sprintf("Finished processing. Elapsed Time = %v", time.Since(start)), logger.Params{"type": t.Type})
				return nil
			})
		})
		mux.Handle(template.WelcomeEmail, h)

		if err := server.Run(mux); err != nil {
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
