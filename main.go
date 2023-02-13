package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"rosetta_exporter/pkg/config"
	prometheusexporter "rosetta_exporter/pkg/prometheus-handlers"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load the environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env vars")
	}

	// Load config
	cfg := config.Get()

	// exporter register
	exporter := prometheusexporter.NewExporter(cfg)
	prometheus.MustRegister(exporter)

	// set handler for metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
