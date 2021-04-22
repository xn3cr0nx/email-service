package meter

import (
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

var (
	errMissingName = errors.New("service name not provided")
)

var meter *metric.Meter

// Config tracer configuration
type Config struct {
	Name string
	Port int
}

// NewMeter singleton implementation returns default meter
func NewMeter(conf *Config) (*metric.Meter, error) {
	if meter == nil {
		if conf.Name == "" {
			return nil, errMissingName
		}

		if err := configure(conf); err != nil {
			return nil, err
		}

		m := global.Meter(conf.Name)
		meter = &m
	}
	return meter, nil
}

func configure(conf *Config) (err error) {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return
	}
	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	}()

	fmt.Println(fmt.Sprintf("Prometheus server running on :%d", conf.Port))
	return
}
