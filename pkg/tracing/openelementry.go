package tracing

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger" // nolint
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const TraceIdCTX = "traceId"

func JaegerTraceProvider(url, env, service string) trace.Tracer {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.DeploymentEnvironmentKey.String(env),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tr := tp.Tracer("App Initialization")

	return tr
}

func TraceIdFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)

	return span.SpanContext().TraceID().String()
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		traceId := TraceIdFromContext(ctx)

		resp, err := handler(ctx, req)
		if err != nil {
			tags := grpc_tags.Extract(ctx)
			tags.Set(TraceIdCTX, traceId)
		}

		return resp, err
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := ss.Context()
		stream := grpc_middleware.WrapServerStream(ss)
		stream.WrappedContext = ctx

		traceId := TraceIdFromContext(ctx)

		err := handler(srv, stream)
		if err != nil {
			tags := grpc_tags.Extract(ctx)
			tags.Set(TraceIdCTX, traceId)
		}
		return err
	}
}

func CreateSpan(ctxt context.Context, tr trace.Tracer, funcName string) (ctx context.Context, span trace.Span) {
	ctx, span = tr.Start(ctxt, funcName)
	traceId := span.SpanContext().TraceID()
	span.SetAttributes(attribute.String(TraceIdCTX, traceId.String()))

	return ctx, span
}
