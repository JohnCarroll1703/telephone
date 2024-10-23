package server

import (
	"context"
	"fmt"
	"git.tarlanpayments.kz/pkg/golog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_sentry "github.com/johnbellone/grpc-middleware-sentry"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"net"
	"telephone/internal/config"
	"telephone/internal/service"
	"telephone/pkg/tracing"
)

type GrpcServer struct {
	server   *grpc.Server
	port     string
	services *service.Services
	log      golog.ContextLogger
	trace    trace.Tracer
}

func NewGRPCServer(
	cfg *config.Config,
	services *service.Services,
	jaegerTrace trace.Tracer,
	prom *prometheus.Registry,
	logger *golog.ZapLogger,
) (*GrpcServer, error) {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	prom.MustRegister(srvMetrics)
	spanProm := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{tracing.TraceIdCTX: span.TraceID().String()}
		}
		return nil
	}
	grpcSrv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(spanProm)),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(spanProm)),
		),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_tags.UnaryServerInterceptor(),
			grpc_sentry.UnaryServerInterceptor(),
			tracing.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_tags.StreamServerInterceptor(),
			grpc_sentry.StreamServerInterceptor(),
			tracing.StreamServerInterceptor(),
		)),
	)
	return &GrpcServer{
		server:   grpcSrv,
		port:     cfg.Service.GrpcPort,
		services: services,
		trace:    jaegerTrace,
		log:      logger,
	}, nil
}

// Here we run the gRPC server
func (gs *GrpcServer) Run() error {
	defer gs.server.GracefulStop()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gs.port))
	if err != nil {
		return fmt.Errorf("failed to listen to gRPC port (%s): %v", gs.port, err)
	}

	return gs.server.Serve(listen)
}
