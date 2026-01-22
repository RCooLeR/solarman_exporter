package exporter

import "solarman_exporter/internal/models"

func NewBMSGroup() *Group {
	return &Group{
		Name:  "bms",
		Gauge: newGroup("solarman_bms_metric", "BMS-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)
			return containsAny(k, "bms", "cell", "balanc", "protect") || containsAny(n, "bms", "cell", "balanc", "protect")
		},
	}
}
