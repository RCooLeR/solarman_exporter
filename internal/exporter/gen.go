package exporter

import "solarman_exporter/internal/models"

func NewGenGroup() *Group {
	return &Group{
		Name:  "gen",
		Gauge: newGroup("solarman_gen_metric", "Generator-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)
			return containsAny(k, "gen", "generator", "dg", "diesel") ||
				containsAny(n, "gen", "generator", "dg", "diesel")
		},
	}
}
