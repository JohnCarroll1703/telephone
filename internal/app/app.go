package app

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"telephone/internal/config"
	"telephone/internal/repository"
	"telephone/internal/server"
	"telephone/internal/service"
	"telephone/pkg/tracing"
	"time"
)

const minGoroutines = 10

func Run(config *config.Config) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(config.Database.PostgreDSN), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}

	tr := tracing.JaegerTraceProvider(
		config.Jeager.DSN,
		config.Service.Environment.String(),
		config.Service.Namespace+"-"+config.Service.AppName)

	if err != nil {
		logger.Fatal(err.Error())
	}

	services := service.NewServices(service.Deps{
		Repos:        repository.NewRepositories(config, tr, db),
		Cgf:          config,
		JeagerTracer: tr,
	})

	promRegistry := newPrometheusRegistry()

	runServers(config, logger, promRegistry, services, tr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	gracefulShutdown(logger)
}

func runGRPCServer(srv *server.GrpcServer, log *zap.Logger) {
	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("🔥 Server stopped due error")
	} else {
		log.Info("✅ Server shutdown successfully")
	}
}

func newPrometheusRegistry() *prometheus.Registry {
	promReg := prometheus.NewRegistry()
	promReg.MustRegister(collectors.NewGoCollector())
	promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	promReg.MustRegister(collectors.NewBuildInfoCollector())
	return promReg
}

func gracefulShutdown(logger *zap.Logger) {
	logger.Error("Shutting down...")

	for {
		time.Sleep(time.Second * 1)

		logger.Error("goroutines")

		if runtime.NumGoroutine() <= minGoroutines {
			break
		}
	}
}

func runServers(
	cfg *config.Config,
	logger *zap.Logger,
	promReg *prometheus.Registry,
	services *service.Services,
	jaegerTrace trace.Tracer) {

	grpcSrv, err := server.NewGRPCServer(cfg, services, jaegerTrace, promReg)
	if err != nil {
		logger.Fatal(err.Error())
	}

	go runGRPCServer(grpcSrv, logger)

	logger.Info("🚀 Starting gRPC server")
}
