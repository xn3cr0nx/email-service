package meter

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	errMissingName = errors.New("service name not provided")
)

var meter = metric.NewNoopMeterProvider().Meter("")

// Config tracer configuration
type Config struct {
	Name string
	Port int
}

// NewMeter singleton implementation returns default meter
func NewMeter(conf *Config) (metric.Meter, error) {
	if meter == nil {
		if conf.Name == "" {
			return nil, errMissingName
		}

		if err := configure(conf); err != nil {
			return nil, err
		}

		m := global.Meter(conf.Name)
		meter = m
	}
	return meter, nil
}

func configure(conf *Config) (err error) {
	config := prometheus.Config{
		// DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		return fmt.Errorf("failed to initialize prometheus exporter: %w", err)
	}

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	}()

	logger.Info("otel", fmt.Sprintf("Prometheus server running on :%d", conf.Port), logger.Params{})
	return
}
