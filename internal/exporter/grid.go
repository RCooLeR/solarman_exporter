package exporter

import "solarman_exporter/internal/models"

func NewGridGroup() *Group {
	g := &Group{
		Name:  "grid",
		Gauge: newGroup("solarman_grid_metric", "Grid-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)
			return containsAny(k, "grid", "ac", "meter") || containsAny(n, "grid", "ac", "meter")
		},
	}
	return g
}
