package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"time"

	"github.com/majidbl/wallet/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *server) runHttpServer() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	s.mapRoutes()

	s.echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		s.echo.Server.ReadTimeout = time.Second * s.cfg.HTTP.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.HTTP.WriteTimeout
		s.echo.Server.MaxHeaderBytes = maxHeaderBytes
		if err := s.echo.Start(s.cfg.HTTP.Port); err != nil {
			s.log.Fatalf("Error starting http Server: ", err)
		}
	}()
}

func (s *server) mapRoutes() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Wallet microservice"
	docs.SwaggerInfo.Description = "Arvan Wallet manager  microservice."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXRequestID,
			csrfTokenHeader,
		},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit(bodyLimit))

}
