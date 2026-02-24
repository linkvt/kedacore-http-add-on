package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// ExternalScaler holds static configuration info for the external scaler
type ExternalScaler struct {
	ServiceName string `env:"KEDAHTTP_OPERATOR_EXTERNAL_SCALER_SERVICE,required"`
	Port        int32  `env:"KEDAHTTP_OPERATOR_EXTERNAL_SCALER_PORT" envDefault:"8091"`
}

type Base struct {
	// The current namespace in which the operator is running.
	CurrentNamespace string `env:"KEDA_HTTP_OPERATOR_NAMESPACE" envDefault:""`
	// The namespace the operator should watch. Leave blank to
	// tell the operator to watch all namespaces.
	WatchNamespace string `env:"KEDA_HTTP_OPERATOR_WATCH_NAMESPACE" envDefault:""`
	// Leader election durations. Nil means use controller-runtime defaults.
	LeaseDuration *time.Duration `env:"KEDA_HTTP_OPERATOR_LEADER_ELECTION_LEASE_DURATION"`
	RenewDeadline *time.Duration `env:"KEDA_HTTP_OPERATOR_LEADER_ELECTION_RENEW_DEADLINE"`
	RetryPeriod   *time.Duration `env:"KEDA_HTTP_OPERATOR_LEADER_ELECTION_RETRY_PERIOD"`
}

func NewBaseFromEnv() (Base, error) {
	return env.ParseAs[Base]()
}

func (e ExternalScaler) HostName(namespace string) string {
	return fmt.Sprintf(
		"%s.%s:%d",
		e.ServiceName,
		namespace,
		e.Port,
	)
}

func NewExternalScalerFromEnv() (ExternalScaler, error) {
	return env.ParseAs[ExternalScaler]()
}
