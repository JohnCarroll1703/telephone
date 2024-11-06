package http

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"telephone/internal/config"
	"telephone/internal/service"
)

type Handler struct {
	jaegerTrace trace.Tracer
	services    *service.Services
	config      *config.Config
}

func NewHandler(
	jaegerTrace trace.Tracer,
	services *service.Services,
	config *config.Config) *Handler {
	return &Handler{
		jaegerTrace: jaegerTrace,
		services:    services,
		config:      config,
	}
}

func (h Handler) Init() *gin.Engine {
	router := gin.New()

	_ = router.Group("/users")
	{

	}

	return router
}
