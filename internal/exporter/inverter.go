package exporter

import "solarman_exporter/internal/models"

func NewInverterGroup() *Group {
	return &Group{
		Name:  "inverter",
		Gauge: newGroup("solarman_inverter_metric", "Inverter-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)
			return containsAny(k, "inv", "inverter", "power", "freq", "temperature", "temp") ||
				containsAny(n, "inv", "inverter", "power", "frequency", "temp", "temperature")
		},
	}
}
