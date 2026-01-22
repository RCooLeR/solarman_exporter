package exporter

import "solarman_exporter/internal/models"

func NewBatteryGroup() *Group {
	return &Group{
		Name:  "battery",
		Gauge: newGroup("solarman_battery_metric", "Battery-related metrics from currentData."),
		Match: func(p models.MetricPoint) bool {
			k := norm(p.Key)
			n := norm(p.Name)
			return containsAny(k, "bat", "battery", "soc", "soh", "charge", "discharge") ||
				containsAny(n, "bat", "battery", "soc", "soh", "charge", "discharge")
		},
	}
}
