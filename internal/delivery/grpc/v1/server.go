package v1

import (
	"git.tarlanpayments.kz/pkg/golog"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/service"
)

type Server struct {
	services *service.Services
	tr       trace.Tracer
	logger   golog.ContextLogger
}

func NewServer(services *service.Services,
	tr trace.Tracer,
	logger golog.ContextLogger) *Server {
	return &Server{
		services: services,
		tr:       tr,
		logger:   logger,
	}
}
