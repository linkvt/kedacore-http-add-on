//go:build e2e

package helpers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// GenerateLoad sends sustained HTTP traffic through the interceptor proxy
// at the given rate for the specified duration. It returns a stop function
// that halts the attack, waits for all results to be collected, and returns
// the aggregated metrics.
func (f *Framework) GenerateLoad(host string, rate vegeta.Rate, duration time.Duration) func() vegeta.Metrics {
	f.t.Helper()

	target := vegeta.Target{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("http://%s/", f.proxyAddr),
		Header: http.Header{"Host": []string{host}},
	}
	targeter := vegeta.NewStaticTargeter(target)

	attacker := vegeta.NewAttacker()
	resultCh := attacker.Attack(targeter, rate, duration, host)

	var wg sync.WaitGroup
	var metrics vegeta.Metrics
	wg.Go(func() {
		for res := range resultCh {
			metrics.Add(res)
		}
		metrics.Close()
	})

	return func() vegeta.Metrics {
		attacker.Stop()
		wg.Wait()
		return metrics
	}
}
