package main

import (
	"github.com/spf13/viper"
	"github.com/xn3cr0nx/email-service/internal/environment"
)

func parseFlags(env *environment.Env) (err error) {
	viper.SetDefault("env", "development")
	viper.SetDefault("rest", true)
	viper.SetDefault("debug", false)
	viper.SetDefault("name", "Mailer")
	viper.SetDefault("template_dir", "templates/")
	viper.SetDefault("sender", "info@test.com")
	viper.SetDefault("frontend_host", "https://frontend.com")
	viper.SetDefault("concurrency", 10)
	viper.SetDefault("queue", "emails")
	viper.SetDefault("http.host", "localhost")
	viper.SetDefault("http.port", 8080)
	viper.SetDefault("provider", "postmark")
	viper.SetDefault("backend", "nats")
	viper.SetDefault("postmark.server", "")
	viper.SetDefault("postmark.account", "")
	viper.SetDefault("sendgrid.api_key", "")
	viper.SetDefault("mailgun.domain", "")
	viper.SetDefault("mailgun.api_key", "")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("asynq.db", 2)
	viper.SetDefault("kafka.addresses", []string{"localhost:6789"})
	viper.SetDefault("kafka.topic", "emails")
	viper.SetDefault("kafka.group", "my-group")
	viper.SetDefault("nats.host", "localhost")
	viper.SetDefault("nats.port", 4222)
	viper.SetDefault("nats.subject", "emails")
	viper.SetDefault("otel.jaeger.enable", false)
	viper.SetDefault("otel.jaeger.host", "jaeger")
	viper.SetDefault("otel.jaeger.port", 14268)
	viper.SetDefault("otel.prometheus.enable", false)
	viper.SetDefault("otel.prometheus.port", 9464)

	rootCmd.Flags().StringVarP(&env.Env, "env", "e", viper.GetString("env"), "Sets environment")
	rootCmd.Flags().BoolVar(&env.Rest, "rest", viper.GetBool("rest"), "Enabled REST API")
	rootCmd.Flags().BoolVarP(&env.Debug, "debug", "d", viper.GetBool("debug"), "Sets logging level to Debug")
	rootCmd.Flags().StringVarP(&env.ServiceName, "name", "n", viper.GetString("name"), "Set service name")
	rootCmd.Flags().StringVar(&env.TemplateDir, "template_dir", viper.GetString("template_dir"), "Define templates folder path")
	rootCmd.Flags().StringVar(&env.Sender, "sender", viper.GetString("sender"), "Set emails sender")
	rootCmd.Flags().StringVar(&env.FrontendHost, "frontend_host", viper.GetString("frontend_host"), "Set frontend host")
	rootCmd.Flags().IntVar(&env.Concurrency, "concurrency", viper.GetInt("concurrency"), "Define templates folder path")
	rootCmd.Flags().StringVar(&env.Queue, "queue", viper.GetString("queue"), "Set queue broker name")
	rootCmd.Flags().StringVarP(&env.Host, "host", "s", viper.GetString("http.host"), "bind http server to host")
	rootCmd.Flags().IntVarP(&env.Port, "port", "p", viper.GetInt("http.port"), "Bind http server to port")
	rootCmd.Flags().StringVar(&env.Provider, "provider", viper.GetString("provider"), "Define which email provider the service is configured to rely on - Options: postmark, sendgrid, mailgun")
	rootCmd.Flags().StringVar(&env.Backend, "backend", viper.GetString("backend"), "Define which backend the service is configured to rely on - Options: asynq, kafka, nats")
	rootCmd.Flags().StringVar(&env.PostmarkServer, "postmark_server", viper.GetString("postmark.server"), "Set postmark server key")
	rootCmd.Flags().StringVar(&env.PostmarkAccount, "postmark_account", viper.GetString("postmark.account"), "Set postmark account key")
	rootCmd.Flags().StringVar(&env.SendgridAPIKey, "sendgrid_api_key", viper.GetString("sendgrid.api_key"), "Set sendgrid api key")
	rootCmd.Flags().StringVar(&env.MailgunDomain, "mailgun_domain", viper.GetString("mailgun.domain"), "Set mailgun domain")
	rootCmd.Flags().StringVar(&env.MailgunAPIKey, "mailgun_api_key", viper.GetString("mailgun.api_key"), "Set mailgun api key")
	rootCmd.Flags().StringVar(&env.RedisHost, "redis_host", viper.GetString("redis.host"), "Set host for redis backend")
	rootCmd.Flags().IntVar(&env.RedisPort, "redis_port", viper.GetInt("redis.port"), "Set port for redis backend")
	rootCmd.Flags().StringVar(&env.RedisPassword, "redis_password", viper.GetString("redis.password"), "Set password for redis backend")
	rootCmd.Flags().IntVar(&env.RedisDB, "redis_db", viper.GetInt("redis.db"), "Set redis database number")
	rootCmd.Flags().StringSliceVar(&env.KafkaAddresses, "kafka_addresses", viper.GetStringSlice("kafka.addresses"), "Set kafka brokers' address")
	rootCmd.Flags().StringVar(&env.KafkaTopic, "kafka_topic", viper.GetString("kafka.topic"), "Set kafka topic")
	rootCmd.Flags().StringVar(&env.KafkaGroup, "kafka_group", viper.GetString("kafka.group"), "Set kafka group")
	rootCmd.Flags().StringVar(&env.NatsHost, "nats_host", viper.GetString("nats.host"), "Set host for nats backend")
	rootCmd.Flags().IntVar(&env.NatsPort, "nats_port", viper.GetInt("nats.port"), "Set port for nats backend")
	rootCmd.Flags().StringVar(&env.NatsSubject, "nats_subject", viper.GetString("nats.subject"), "Set subject for nats subscriber")
	rootCmd.Flags().BoolVar(&env.OtelExporterJaegerEnable, "otel_exporter_jaeger_enable", viper.GetBool("otel.jaeger.enable"), "Enable OpenTelemetry based jager tracing")
	rootCmd.Flags().StringVar(&env.OtelExporterJaegerAgentHost, "otel_exporter_jaeger_agent_host", viper.GetString("otel.jaeger.host"), "Override Jaeger agent hostname")
	rootCmd.Flags().IntVar(&env.OtelExporterJaegerAgentPort, "otel_exporter_jaeger_agent_port", viper.GetInt("otel.jaeger.port"), "Override Jaeger agent port")
	rootCmd.Flags().BoolVar(&env.OtelExporterPrometheusEnable, "otel_exporter_prometheus_enable", viper.GetBool("otel.prometheus.enable"), "Enable OpenTelemetry based prometheus metrics")
	rootCmd.Flags().IntVar(&env.OtelExporterPrometheusPort, "otel_exporter_prometheus_port", viper.GetInt("otel.prometheus.port"), "Override Prometheus exposed port")

	if err = viper.BindPFlag("env", rootCmd.Flags().Lookup("env")); err != nil {
		return
	}
	if err = viper.BindPFlag("rest", rootCmd.Flags().Lookup("rest")); err != nil {
		return
	}
	if err = viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug")); err != nil {
		return
	}
	if err = viper.BindPFlag("name", rootCmd.Flags().Lookup("name")); err != nil {
		return
	}
	if err = viper.BindPFlag("template_dir", rootCmd.Flags().Lookup("template_dir")); err != nil {
		return
	}
	if err = viper.BindPFlag("sender", rootCmd.Flags().Lookup("sender")); err != nil {
		return
	}
	if err = viper.BindPFlag("frontend_host", rootCmd.Flags().Lookup("frontend_host")); err != nil {
		return
	}
	if err = viper.BindPFlag("concurrency", rootCmd.Flags().Lookup("concurrency")); err != nil {
		return
	}
	if err = viper.BindPFlag("queue", rootCmd.Flags().Lookup("queue")); err != nil {
		return
	}

	if err = viper.BindPFlag("http.host", rootCmd.Flags().Lookup("host")); err != nil {
		return
	}
	if err = viper.BindPFlag("http.port", rootCmd.Flags().Lookup("port")); err != nil {
		return
	}
	if err = viper.BindPFlag("provider", rootCmd.Flags().Lookup("provider")); err != nil {
		return
	}
	if err = viper.BindPFlag("backend", rootCmd.Flags().Lookup("backend")); err != nil {
		return
	}
	if err = viper.BindPFlag("postmark.server", rootCmd.Flags().Lookup("postmark_server")); err != nil {
		return
	}
	if err = viper.BindPFlag("postmark.account", rootCmd.Flags().Lookup("postmark_account")); err != nil {
		return
	}
	if err = viper.BindPFlag("sendgrid.api_key", rootCmd.Flags().Lookup("sendgrid_api_key")); err != nil {
		return
	}
	if err = viper.BindPFlag("mailgun.domain", rootCmd.Flags().Lookup("mailgun_domain")); err != nil {
		return
	}
	if err = viper.BindPFlag("mailgun.api_key", rootCmd.Flags().Lookup("mailgun_api_key")); err != nil {
		return
	}
	if err = viper.BindPFlag("redis.host", rootCmd.Flags().Lookup("redis_host")); err != nil {
		return
	}
	if err = viper.BindPFlag("redis.port", rootCmd.Flags().Lookup("redis_port")); err != nil {
		return
	}
	if err = viper.BindPFlag("redis.password", rootCmd.Flags().Lookup("redis_password")); err != nil {
		return
	}
	if err = viper.BindPFlag("redis.db", rootCmd.Flags().Lookup("redis_db")); err != nil {
		return
	}
	if err = viper.BindPFlag("kafka.addresses", rootCmd.Flags().Lookup("kafka_addresses")); err != nil {
		return
	}
	if err = viper.BindPFlag("kafka.topic", rootCmd.Flags().Lookup("kafka_topic")); err != nil {
		return
	}
	if err = viper.BindPFlag("kafka.group", rootCmd.Flags().Lookup("kafka_group")); err != nil {
		return
	}
	if err = viper.BindPFlag("nats.host", rootCmd.Flags().Lookup("nats_host")); err != nil {
		return
	}
	if err = viper.BindPFlag("nats.port", rootCmd.Flags().Lookup("nats_port")); err != nil {
		return
	}
	if err = viper.BindPFlag("nats.subject", rootCmd.Flags().Lookup("nats_subject")); err != nil {
		return
	}
	if err = viper.BindPFlag("otel.jaeger.enable", rootCmd.Flags().Lookup("otel_exporter_jaeger_enable")); err != nil {
		return
	}
	if err = viper.BindPFlag("otel.jaeger.host", rootCmd.Flags().Lookup("otel_exporter_jaeger_agent_host")); err != nil {
		return
	}
	if err = viper.BindPFlag("otel.jaeger.port", rootCmd.Flags().Lookup("otel_exporter_jaeger_agent_port")); err != nil {
		return
	}
	if err = viper.BindPFlag("otel.prometheus.enable", rootCmd.Flags().Lookup("otel_exporter_prometheus_enable")); err != nil {
		return
	}
	if err = viper.BindPFlag("otel.prometheus.port", rootCmd.Flags().Lookup("otel_exporter_prometheus_port")); err != nil {
		return
	}

	environment.Set(env)

	return
}
