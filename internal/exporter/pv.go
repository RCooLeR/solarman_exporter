package exporter

import "solarman_exporter/internal/models"

func NewPVGroup() *Group {
	return &Group{
		Name:  "pv",
		Gauge: newGroup("solarman_pv_metric", "PV-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)
			return containsAny(k, "pv", "solar", "string", "dc") || containsAny(n, "pv", "solar", "string", "dc")
		},
	}
}
