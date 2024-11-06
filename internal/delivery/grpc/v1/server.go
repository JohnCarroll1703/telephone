package v1

import (
	"go.opentelemetry.io/otel/trace"
	pb "telephone/internal/proto"
	"telephone/internal/service"
)

type Server struct {
	services *service.Services
	tr       trace.Tracer
	pb.UnimplementedUserContactServiceServer
}

func NewServer(services *service.Services,
	tr trace.Tracer) *Server {
	return &Server{
		services: services,
		tr:       tr,
	}
}
