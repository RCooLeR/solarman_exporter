package exporter

import (
	"strings"

	"solarman_exporter/internal/models"

	"github.com/prometheus/client_golang/prometheus"
)

type Group struct {
	Name  string
	Gauge *prometheus.GaugeVec
	Match func(p models.MetricPoint) bool
}

func newGroup(metricName, help string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: metricName, Help: help},
		[]string{"device_sn", "key", "name", "unit"},
	)
}

func norm(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
