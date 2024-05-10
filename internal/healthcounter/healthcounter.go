// Package healthcounter provides a health counter that can be used to determine if a service is healthy or not.
package healthcounter

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/opentargets/opensearch-nanny/internal/config"
)

// HealthCounter is a health counter that can be used to determine if a service is healthy or not.
type HealthCounter struct {
	Healthy                  bool
	HealthURL                string
	SecondsInGreen           int
	SecondsInGreenForHealthy int
	TickerInterval           int
	m                        sync.Mutex
}

// New creates a new HealthCounter.
func New(oc config.OpensearchConfig) *HealthCounter {
	hc := &HealthCounter{
		Healthy:                  false,
		HealthURL:                oc.HealthURL,
		SecondsInGreen:           0,
		SecondsInGreenForHealthy: oc.SecondsInGreenForHealthy,
		TickerInterval:           oc.TickerInterval,
	}

	return hc
}

type clusterHealthResponse struct {
	ClusterName string `json:"cluster_name"`
	Status      string `json:"status"`
	TimedOut    bool   `json:"timed_out"`
}

// Start starts the health counter.
func (hc *HealthCounter) Start() {
	ticker := time.NewTicker(time.Duration(hc.TickerInterval) * time.Second)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				slog.Debug("tick", slog.Time("ticker", t), slog.Int("seconds_in_green", hc.SecondsInGreen))

				resp, err := http.Get(hc.HealthURL)
				if err != nil {
					slog.Error("error getting health URL", slog.Any("error", err))
					hc.m.Lock()
					hc.Healthy = false
					hc.m.Unlock()
					break
				}

				r := clusterHealthResponse{}
				err = json.NewDecoder(resp.Body).Decode(&r)
				if err != nil {
					slog.Error("error decoding JSON response", slog.Any("error", err))
					hc.m.Lock()
					hc.Healthy = false
					hc.m.Unlock()
					break
				}

				if r.Status == "green" {
					hc.m.Lock()
					hc.SecondsInGreen++
					hc.m.Unlock()

					slog.Debug("status is green", slog.Int("seconds_in_green", hc.SecondsInGreen))

					if !hc.Healthy && hc.SecondsInGreen >= hc.SecondsInGreenForHealthy {
						slog.Info("status has been green for long enough, setting healthy")
						hc.Healthy = true
					}
				} else {
					slog.Info("status is not green, resetting counter", slog.String("status", r.Status))
					hc.m.Lock()
					hc.SecondsInGreen = 0
					hc.Healthy = false
					hc.m.Unlock()

				}
			}
			time.Sleep(time.Second)
		}
	}()
}
