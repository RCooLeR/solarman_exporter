package solarman

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"solarman_exporter/internal/models"

	"github.com/rs/zerolog/log"
)

func (c *Client) StationDevices(ctx context.Context, stationID int64) ([]models.Device, error) {
	body := map[string]any{"stationId": stationID}

	raw, status, err := c.doJSONAuthRetry(ctx, "POST", fmt.Sprintf("/station/%s/device", c.cfg.APIVersion), false, body)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("station device failed: status=%d body=%s", status, strings.TrimSpace(string(raw)))
	}
	log.Debug().Str("devices", string(raw)).Msg("station device list")
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	listAny, _ := m["deviceListItems"].([]any)
	out := make([]models.Device, 0, len(listAny))
	for _, it := range listAny {
		obj, ok := it.(map[string]any)
		if !ok {
			continue
		}
		sn := firstString(obj, "deviceSn")
		if sn == "" {
			continue
		}
		name := firstString(obj, "deviceName", "name")
		out = append(out, models.Device{
			DeviceID: toInt64(obj["deviceId"]),
			DeviceSN: sn,
			Name:     name,
		})
	}
	return out, nil
}
