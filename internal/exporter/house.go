package exporter

import "solarman_exporter/internal/models"

func NewHouseGroup() *Group {
	return &Group{
		Name:  "house",
		Gauge: newGroup("solarman_house_metric", "House/consumption-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)

			// consumption/house-only signals
			return containsAny(k, "house", "home", "consumption", "consume", "load_cons", "pcons", "p_load_total") ||
				containsAny(n, "house", "home", "consumption", "consumed", "usage", "demand")
		},
	}
}
