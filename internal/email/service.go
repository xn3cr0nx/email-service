package email

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xn3cr0nx/email-service/internal/provider"
	"github.com/xn3cr0nx/email-service/pkg/validator"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Service interface exports available methods for user service
type Service interface {
	Send(context.Context, *WelcomeEmailBody) error
	// SendBatch() error
}

type service struct {
	Mailer provider.Mailer
	tracer trace.Tracer
	meter  metric.Meter
}

// NewService instantiates a new Service layer for customer
func NewService(m provider.Mailer, tracer trace.Tracer, meter metric.Meter) *service {
	return &service{
		Mailer: m,
		tracer: tracer,
		meter:  meter,
	}
}

// email godoc
// @ID email
//
// @Router /email [post]
// @Summary Email
// @Description Process email request
// @Tags email
//
// @Accept  json
// @Produce  json
//
// @Param email body WelcomeEmailBody true "welcome email parameters"
//
// @Success 200 {string} Ok
// @Failure 400 {string} string
// @Failure 500 {string} string
func Handler(s Service) func(echo.Context) error {
	return func(c echo.Context) error {
		b := new(WelcomeEmailBody)
		if err := validator.Struct(&c, b); err != nil {
			return err
		}

		err := s.Send(c.Request().Context(), b)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, "Ok")
	}
}

// Send processes email request and send using injected email client
func (s *service) Send(ctx context.Context, body *WelcomeEmailBody) (err error) {
	err = body.Process(ctx, s.Mailer)
	return
}
