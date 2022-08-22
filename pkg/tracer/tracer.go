package tracer

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = trace.NewNoopTracerProvider().Tracer("")
var middleware echo.MiddlewareFunc

// errorHandler is our custom Open Telemetry error handler
type errorHandler struct{}

// Handle will log and report to Bugsnag errors produced by the Open Telemetry library
func (e *errorHandler) Handle(err error) {
	logger.Error("otel", err, logger.Params{})
}

func NewTracerProvider(ctx context.Context, service string) (trace.Tracer, error) {
	exp, err := otlptrace.New(ctx, otlptracegrpc.NewClient())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(service),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(res),
	)

	otel.SetErrorHandler(&errorHandler{})
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	tracer = tp.Tracer(service)
	middleware = otelecho.Middleware(service)

	return tracer, nil
}

// func HandlersFilter(filter []string) otelhttp.Option {
// 	pathFilter := otelhttp.Filter(func(req *http.Request) bool {
// 		return !utils.StringInSlice(req.URL.Path, filter)
// 	})

// 	return otelhttp.WithFilter(pathFilter)
// }

func NewTracerSpan(ctx context.Context, name string) trace.Span {
	_, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
	return span
}

func Middleware() echo.MiddlewareFunc {
	return middleware
}
