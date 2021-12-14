package jaeger_client

import (
    "context"
    "errors"
    "fmt"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    tracesdk "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
    "go.opentelemetry.io/otel/trace"
    "net/http"
    "os"
    "time"
)

type contextKey int

const spanNameContextKey contextKey = iota

// ErrInvalidHost is used when the Jaeger Agent host is either not given or invalid.
var ErrInvalidHost = errors.New("jaeger: invalid agent host")
var name string
var cleanShutdownFunc = func() error { return nil }

// InitializeJaeger initializes a tracer provider that sends traces to Jaeger, and then sets it as
// the global tracer provider.
// additionalAttributes is an optional parameter that specifies additional attributes that should be added to every trace.
func InitializeJaeger(serviceName string, additionalAttributes ...attribute.KeyValue) error {
    // if there is no URL set for the trace collector, do not configure the tracer.
    // This is fine, it just means that any traces we generate will be discarded.
    // Prevents any weirdness involving traces when running the service locally.
    jaegerHost := os.Getenv("OTEL_EXPORTER_JAEGER_AGENT_HOST")
    if jaegerHost == "" {
        return ErrInvalidHost
    }

    if serviceName == "" {
        return errors.New("jaeger: invalid service name")
    }

    // store this at the package level for later use
    name = serviceName

    exp, err := jaeger.New(jaeger.WithAgentEndpoint())
    if err != nil {
        return err
    }

    attributes := []attribute.KeyValue{semconv.ServiceNameKey.String(name)}
    attributes = append(attributes, additionalAttributes...)

    tracerProvider := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exp),
        tracesdk.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            attributes...,
        )),
    )
    otel.SetTracerProvider(tracerProvider)

    // set up inter-service trace propagation
    otel.SetTextMapPropagator(
        propagation.NewCompositeTextMapPropagator(
            propagation.TraceContext{},
        ),
    )

    // set up a function that will ensure all traces get flushed from memory before the process is killed.
    // This has a timeout of 5 seconds, so it will not hang indefinitely if there is some problem with flushing.
    cleanShutdownFunc = func() error {
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
        defer cancel()
        if err := tracerProvider.Shutdown(ctx); err != nil {
            return fmt.Errorf("failed to cleanly shut down tracer provider: %w", err)
        }

        return nil
    }

    return nil
}

func Shutdown() error {
    return cleanShutdownFunc()
}

func Trace() trace.Tracer {
    tp := otel.GetTracerProvider()

    return tp.Tracer(name)
}

func SpanFromContext(ctx context.Context, name string) (context.Context, trace.Span) {
    tr := Trace()

    return tr.Start(ctx, name)
}

func SpanNameFromContext(ctx context.Context) string {
    spanName, ok := ctx.Value(spanNameContextKey).(string)
    if !ok {
        return ""
    }

    return spanName
}

func InjectSpanName(ctx context.Context, name string) context.Context {
    return context.WithValue(ctx, spanNameContextKey, name)
}

func UninjectSpanName(ctx context.Context) context.Context {
    // injecting a blank span name is effectively the same as "uninjecting" it
    return InjectSpanName(ctx, "")
}

func SpanNameFormatter(operation string, req *http.Request) string {
    spanName := SpanNameFromContext(req.Context())
    if spanName == "" {
        spanName = fmt.Sprintf("%s %s", req.Method, req.URL.Path)
    }

    return spanName
}
