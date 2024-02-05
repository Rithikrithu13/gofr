package exporters

import (
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Prometheus(appName, appVersion string) metric.Meter {
	exporter, err := prometheus.New(prometheus.WithoutTargetInfo())
	if err != nil {
		return nil
	}

	meter := metricSdk.NewMeterProvider(
		metricSdk.WithReader(exporter),
		metricSdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		))).Meter(appName, metric.WithInstrumentationVersion(appVersion))

	return meter
}

func OTLPStdOut(appName, appVersion string) metric.Meter {
	exporter, err := stdoutmetric.New()
	if err != nil {
		return nil
	}

	meter := metricSdk.NewMeterProvider(
		metricSdk.WithResource(resource.NewSchemaless(semconv.ServiceName(appName))),
		metricSdk.WithReader(metricSdk.NewPeriodicReader(exporter,
			metricSdk.WithInterval(3*time.Second)))).Meter(appName, metric.WithInstrumentationVersion(appVersion))

	return meter
}

func OTLPMetricHTTP(appName, appVersion string) metric.Meter {
	exporter, err := otlpmetrichttp.New(nil,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithURLPath("/metrics"),
		otlpmetrichttp.WithEndpoint("localhost:8000"))
	if err != nil {
		return nil
	}

	meter := metricSdk.NewMeterProvider(metricSdk.WithReader(metricSdk.NewPeriodicReader(exporter,
		metricSdk.WithInterval(3*time.Second)))).Meter(appName, metric.WithInstrumentationVersion(appVersion))

	return meter
}
