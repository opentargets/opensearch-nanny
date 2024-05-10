// Package main contains the main function that starts the application.
package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"flag"

	"github.com/opentargets/opensearch-nanny/internal/config"
	"github.com/opentargets/opensearch-nanny/internal/healthcounter"
	"github.com/opentargets/opensearch-nanny/internal/log"
)

func main() {
	var p string

	flag.StringVar(&p, "c", "./config.toml", "Configuration file path")
	flag.StringVar(&p, "config", "./config.toml", "Configuration file path")
	flag.Parse()

	c := config.InitConfig(p)
	log.InitLogger(c.Server)

	hc := healthcounter.New(c.Opensearch)
	hc.Start()

	h := func(w http.ResponseWriter, r *http.Request) {
		if hc.Healthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Service unavailable"))
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.Server.Address, c.Server.Port),
		Handler: http.HandlerFunc(h),
	}

	slog.Info("Starting server", slog.String("address", srv.Addr), slog.String("health_url", c.Opensearch.HealthURL))
	srv.ListenAndServe()
}
