package tracer

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdkexporter "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

var (
	errMissingName = errors.New("service name not provided")
)

var tracer *trace.Tracer
var middleware echo.MiddlewareFunc

// Exporter list of exporter options
type Exporter int

const (
	// Stdout exporter to log tracing to stdout
	Stdout Exporter = iota
	// Jaeger export to send tracing to jaeger
	Jaeger
)

// Config tracer configuration
type Config struct {
	Host        string
	Port        int
	Name        string
	Exporter    Exporter
	Environment string
}

// NewTracer singleton implementation returns default tracer
func NewTracer(conf *Config) (*trace.Tracer, error) {
	if tracer == nil {
		if conf.Name == "" {
			return nil, errMissingName
		}

		t := otel.Tracer(conf.Name)
		tracer = &t

		if err := configure(conf); err != nil {
			return nil, err
		}
	}

	middleware = otelecho.Middleware(conf.Name)
	return tracer, nil
}

func Middleware() echo.MiddlewareFunc {
	return middleware
}

func configure(conf *Config) (err error) {
	var exporter sdkexporter.SpanExporter

	switch conf.Exporter {
	case Jaeger:
		exporter, err = jaeger.NewRawExporter(
			jaeger.WithCollectorEndpoint(
				fmt.Sprintf("http://%s:%d/api/traces", conf.Host, conf.Port)))
		if err != nil {
			return
		}
	default:
		exporter, err = stdout.NewExporter(stdout.WithPrettyPrint())
		if err != nil {
			return
		}
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(conf.Name),
			attribute.String("environment", conf.Environment),
			// attribute.Int64("ID", id),
		)),
	)
	if err != nil {
		return
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return
}
