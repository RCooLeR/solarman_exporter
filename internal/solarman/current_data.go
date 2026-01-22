package solarman

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"solarman_exporter/internal/models"

	"github.com/rs/zerolog/log"
)

func (c *Client) CurrentData(ctx context.Context, deviceSN string) ([]models.MetricPoint, error) {
	body := map[string]any{"deviceSn": deviceSN}

	raw, status, err := c.doJSONAuthRetry(ctx, "POST", fmt.Sprintf("/device/%s/currentData", c.cfg.APIVersion), true, body)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("currentData failed: status=%d body=%s", status, strings.TrimSpace(string(raw)))
	}
	log.Info().Msg("Obtained current data")
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	if s, ok := m["success"].(bool); ok && !s {
		msg, _ := m["msg"].(string)
		if msg == "" {
			msg = "unknown error"
		}
		return nil, fmt.Errorf("solarman error: %s", msg)
	}

	dataListAny, _ := m["dataList"].([]any)
	points := make([]models.MetricPoint, 0, len(dataListAny))

	for _, it := range dataListAny {
		obj, ok := it.(map[string]any)
		if !ok {
			continue
		}
		key := firstString(obj, "key", "dataKey", "id", "sn")
		name := firstString(obj, "name", "dataName", "title", "paramName")
		unit := firstString(obj, "unit", "dataUnit")

		val, ok := toFloat64(obj["value"])
		if !ok {
			val, ok = toFloat64(obj["val"])
			if !ok {
				continue
			}
		}
		if key == "" {
			key = sanitizeKey(name)
		}
		points = append(points, models.MetricPoint{Key: key, Name: name, Unit: unit, Value: val})
	}

	return points, nil
}
