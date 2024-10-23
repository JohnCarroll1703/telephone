package http

import (
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/service"
)

type Handler struct {
	jaegerTrace trace.Tracer
	logger      golog.ContextLogger
	services    *service.Services
	config      *config.Config
}

func NewHandler(
	jaegerTrace trace.Tracer,
	logger golog.ContextLogger,
	services *service.Services,
	config *config.Config) *Handler {
	return &Handler{
		jaegerTrace: jaegerTrace,
		logger:      logger,
		services:    services,
		config:      config,
	}
}
