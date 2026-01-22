package exporter

import "solarman_exporter/internal/models"

func NewLoadGroup() *Group {
	return &Group{
		Name:  "load",
		Gauge: newGroup("solarman_load_metric", "Load/UPS/backup output metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)

			// keep this group for UPS/backup/output load (NOT general household consumption)
			return containsAny(k, "load", "ups", "backup", "eps", "critical", "output") ||
				containsAny(n, "load", "ups", "backup", "eps", "critical", "output")
		},
	}
}
