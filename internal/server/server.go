package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/etherlabsio/healthcheck"
	v "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xn3cr0nx/email-service/pkg/pprof"
	"github.com/xn3cr0nx/email-service/pkg/validator"
)

// Server struct initialized with port
type (
	Server struct {
		port   string
		router *echo.Echo
	}
)

const (
	defaultStatusTimeout = 5 * time.Second
)

var server *Server

// NewServer singleton pattern that returns pointer to server
func NewServer(port int) *Server {
	if server != nil {
		return server
	}
	server = &Server{
		port:   fmt.Sprintf(":%d", port),
		router: echo.New(),
	}
	return server
}

func (s *Server) Listen() {
	pprof.Wrap(s.router)

	s.router.HideBanner = true
	s.router.Debug = viper.GetBool("debug")
	s.router.Use(middleware.Logger())
	s.router.Logger.SetLevel(log.INFO)
	s.router.Validator = validator.NewValidator()

	s.router.HTTPErrorHandler = customHTTPErrorHandler

	s.router.Use(middleware.Recover())
	s.router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	s.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: viper.GetStringSlice(("auth.cors")),
		AllowMethods: viper.GetStringSlice(("auth.methods")),
	}))

	s.router.Use(middleware.RequestID())

	s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	s.router.GET("/status", handleStatus())

	log.Printf(
		"mailer (PID: %d) is starting on %s\n=> Ctrl-C to shutdown server\n",
		os.Getpid(),
		s.port)
	go func() {
		if err := s.router.Start(s.port); err != nil {
			s.router.Logger.Error(err)
			s.router.Logger.Fatal("shutting down the server")
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.router.Logger.Info("signal caught. gracefully shutting down...")
	if err := s.router.Shutdown(ctx); err != nil {
		s.router.Logger.Fatal(err)
	}
}

func timeout() time.Duration {
	timeoutMillis, err := strconv.Atoi(os.Getenv("TIMEOUT_MILLIS"))
	if err != nil {
		log.Panic(err)
	}
	return time.Duration(timeoutMillis) * time.Millisecond
}

func handleStatus() echo.HandlerFunc {
	timeout := healthcheck.WithTimeout(defaultStatusTimeout)
	handler := healthcheck.HandlerFunc(
		timeout,
	)
	return echo.WrapHandler(handler)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	code := http.StatusInternalServerError
	m := ""

	if e, ok := err.(*echo.HTTPError); ok {
		code = e.Code
		if httpError, ok := e.Message.(*echo.HTTPError); ok {
			m = httpError.Message.(string)
		} else if _, ok := e.Message.(v.ValidationErrors); ok {
		} else {
			if stringError, ok := e.Message.(string); ok {
				m = stringError
			} else {
				m = err.Error()
			}
		}
	}

	message := map[string]interface{}{"code": code, "error": http.StatusText(code)}
	if m != "" && m != message["error"] {
		message["type"] = m
	}
	c.JSON(code, message)
}
