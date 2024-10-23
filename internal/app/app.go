package app

import (
	"context"
	"errors"
	"fmt"
	"git.tarlanpayments.kz/pkg/golog"
	"git.tarlanpayments.kz/pkg/gosentry"
	"git.tarlanpayments.kz/processing/jusan/docs"
	"git.tarlanpayments.kz/processing/jusan/pkg/openapi"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel/trace"
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
	logger, err := createLogger(config)
	if err != nil {
		panic(err)
	}

	openApiInit(config, logger)

	postgres, err := pgx.Connect(context.Background(), config.Database.PostgreDSN)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	err = gosentry.SentryInit(config.Sentry.DSN, config.Service.Environment.String())
	if err != nil {
		logger.Fatalw(err.Error())
	}

	tr := tracing.JaegerTraceProvider(
		config.Jeager.DSN,
		config.Service.Environment.String(),
		config.Service.Namespace+"-"+config.Service.AppName)

	if err != nil {
		logger.Fatalw(err.Error())
	}

	services := service.NewServices(service.Deps{
		Repos:        repository.NewRepositories(config, tr, logger, postgres),
		Cgf:          config,
		Logger:       logger,
		JeagerTracer: tr,
	})

	promRegistry := newPrometheusRegistry()

	runServers(config, logger, promRegistry, services, tr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	gracefulShutdown(logger)
}

func openApiInit(cfg *config.Config, log *golog.ZapLogger) {
	if cfg.Service.Domain == "127.0.0.1" || cfg.Service.Domain == "localhost" {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Service.Port)
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	} else {
		docs.SwaggerInfo.Host = cfg.Service.Domain
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}

	if err := openapi.NewOpenApiClient(
		cfg.Service.AppName,
		cfg.Service.Openapi,
		docs.SwaggerInfo.ReadDoc()).Send(
		context.Background()).Error(); err != nil {
		log.Fatalw(err.Error())
	}
}

func runGRPCServer(srv *server.GrpcServer, log *golog.ZapLogger) {
	if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("🔥 Server stopped due error", "err", err.Error())
	} else {
		log.Infow("✅ Server shutdown successfully")
	}
}

func createLogger(cfg *config.Config) (*golog.ZapLogger, error) {
	loggerConfig := golog.Config{
		Mode:              golog.ProductionMode,
		Level:             golog.InfoLevel,
		AppName:           cfg.Service.AppName,
		DisableStacktrace: true,
	}

	if cfg.Service.Environment.IsLocal() {
		loggerConfig.Mode = golog.DevelopmentMode
		loggerConfig.Level = golog.DebugLevel
	}

	if cfg.Service.Environment.IsProduction() {
		loggerConfig.Mode = golog.ProductionMode
		loggerConfig.Level = golog.ErrorLevel
	}

	return golog.NewZapLogger(loggerConfig)
}

func newPrometheusRegistry() *prometheus.Registry {
	promReg := prometheus.NewRegistry()
	promReg.MustRegister(collectors.NewGoCollector())
	promReg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	promReg.MustRegister(collectors.NewBuildInfoCollector())
	return promReg
}

func gracefulShutdown(logger *golog.ZapLogger) {
	logger.Errorw("Shutting down...")

	for {
		time.Sleep(time.Second * 1)

		logger.Errorw("goroutines", "count", runtime.NumGoroutine())

		if runtime.NumGoroutine() <= minGoroutines {
			break
		}
	}
}

func runServers(
	cfg *config.Config,
	logger *golog.ZapLogger,
	promReg *prometheus.Registry,
	services *service.Services,
	jaegerTrace trace.Tracer) {

	grpcSrv, err := server.NewGRPCServer(cfg, services, jaegerTrace, promReg, logger)
	if err != nil {
		logger.Fatalw(err.Error())
	}

	go runGRPCServer(grpcSrv, logger)

	logger.Infow("🚀 Starting gRPC server")
}
