module github.com/ianling/jaeger-opentelemetry-client

go 1.17

require (
	go.opentelemetry.io/contrib/propagators/jaeger v1.3.0
	go.opentelemetry.io/contrib/propagators/ot v1.3.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
)

require (
	github.com/go-logr/logr v1.2.1 // indirect
	github.com/go-logr/stdr v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
)
