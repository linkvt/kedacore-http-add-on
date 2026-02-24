package config

import (
	"github.com/caarlos0/env/v11"
)

// Metrics is the configuration for configuring metrics in the interceptor.
type Metrics struct {
	// Sets whether or not to enable the Prometheus metrics exporter
	OtelPrometheusExporterEnabled bool `env:"OTEL_PROM_EXPORTER_ENABLED" envDefault:"true"`
	// Sets the port which the Prometheus compatible metrics endpoint should be served on
	OtelPrometheusExporterPort int `env:"OTEL_PROM_EXPORTER_PORT" envDefault:"2223"`
	// Sets whether or not to enable the OTEL metrics exporter
	OtelHTTPExporterEnabled bool `env:"OTEL_EXPORTER_OTLP_METRICS_ENABLED" envDefault:"false"`
}

// MustParseMetrics parses standard configs and returns the
// newly created config. Panics if parsing fails.
func MustParseMetrics() Metrics {
	return env.Must(env.ParseAs[Metrics]())
}
