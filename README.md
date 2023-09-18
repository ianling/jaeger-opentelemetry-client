# jaeger-opentelemetry-client

Initializes all the necessary `otel` things to enable sending traces to Jaeger.

Simply call the `InitializeJaeger()` function at the beginning of your application and provide a Jaeger Collector host via
the `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` environment variable. See: https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/#otel_exporter_otlp_traces_endpoint

Example: `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://my-api-endpoint:4318/v1/traces`

This package is safe to import and use even if you do not provide a Jaeger Collector host.
This is useful when running code in a local development environment, for example,
where you may not have access to a Jaeger Agent. All calls to this package, as well as the `otel` package, effectively
become no-ops.

# Usage

```go
package main

import (
    "context"
    "github.com/ianling/jaeger-opentelemetry-client"
    "log"
)

func main() {
    cleanShutdownFunc, err := jaeger_client.InitializeJaeger("service-name")
    if err == jaeger_client.ErrInvalidHost {
        log.Println("No Jaeger Agent host provided, no traces will be sent out!")
    } else if err != nil {
        log.Fatalf("Failed to initialize Jaeger: %v", err)
    }

    defer cleanShutdownFunc() // flushes remaining traces before the process ends, so they don't get lost

    DoThing() // example code
}

func DoThing() {
    // this trace object is just a plain OpenTelemetry trace from the otel package,
    // so you use it as you normally would.
    trace := jaeger_client.Trace()
    _, span := trace.Start(context.Background(), "DoThing")
    // ...
    defer span.End()
}
```
