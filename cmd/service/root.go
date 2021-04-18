package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xn3cr0nx/email-service/internal/mailer/postmark"
	"github.com/xn3cr0nx/email-service/internal/server"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
)

var (
	port        int
	debug, rest bool
	target      string
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
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Sets logging level to Debug")
	rootCmd.Flags().IntVar(&port, "port", 6066, "Bind http server to port")
	rootCmd.PersistentFlags().BoolVar(&rest, "rest", false, "Enable exposed REST API to interact with the service")
	rootCmd.PersistentFlags().StringVar(&target, "target", "", "Sets if mailer should run once with a specific target")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("debug", false)
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	viper.SetDefault("rest", false)
	viper.BindPFlag("rest", rootCmd.PersistentFlags().Lookup("rest"))

	viper.SetDefault("port", 6066)
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	viper.SetDefault("target", "")
	viper.BindPFlag("target", rootCmd.PersistentFlags().Lookup("target"))

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

	// initialize template cache
	templateDir := viper.GetString("template_dir")
	if templateDir == "" {
		templateDir = "templates/"
	}
	_, err := template.NewTemplateCache(&templateDir)
	if err != nil {
		panic(fmt.Errorf("Cannot initialize template cache: %w", err))
	}

	mailer := postmark.NewClient(viper.GetString("postmark.server"), viper.GetString("postmark.account"))

	if viper.GetBool("server.enabled") || viper.GetBool("rest") {
		s := server.NewServer(viper.GetInt("server.port"), mailer)
		s.Listen()
	}
}
