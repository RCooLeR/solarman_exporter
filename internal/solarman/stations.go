package solarman

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"solarman_exporter/internal/models"

	"github.com/rs/zerolog/log"
)

func (c *Client) Stations(ctx context.Context) ([]models.Station, error) {
	raw, status, err := c.doJSONAuthRetry(ctx, "POST", fmt.Sprintf("/station/%s/list", c.cfg.APIVersion), false, map[string]any{})
	if err != nil {
		log.Error().Err(err).Msg("error getting stations")
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("station list failed: status=%d body=%s", status, strings.TrimSpace(string(raw)))
	}
	log.Debug().Str("stations", string(raw)).Msg("stations list")
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	listAny, _ := m["stationList"].([]any)
	out := make([]models.Station, 0, len(listAny))
	for _, it := range listAny {
		obj, ok := it.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, models.Station{
			ID:   toInt64(obj["id"]),
			Name: firstString(obj, "name"),
		})
	}
	return out, nil
}
