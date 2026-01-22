package exporter

import (
	"context"
	"strings"
	"sync"
	"time"

	"solarman_exporter/internal/models"
	"solarman_exporter/internal/solarman"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Client       *solarman.Client
	PollInterval time.Duration
	DeviceSNs    []string
	StationID    int64

	EnableGeneric bool
}

type Exporter struct {
	cfg Config

	mu          sync.RWMutex
	ready       bool
	lastSuccess time.Time

	deviceUp   *prometheus.GaugeVec
	lastUpdate *prometheus.GaugeVec

	groups []*Group

	generic *prometheus.GaugeVec
}

func New(cfg Config) *Exporter {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 60 * time.Second
	}

	e := &Exporter{
		cfg: cfg,
		deviceUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "solarman_device_up", Help: "1 if last poll succeeded for device, else 0."},
			[]string{"device_sn"},
		),
		lastUpdate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "solarman_last_update_timestamp_seconds", Help: "Unix timestamp of last successful refresh per device."},
			[]string{"device_sn"},
		),
		groups: []*Group{
			NewGridGroup(),
			NewPVGroup(),
			NewInverterGroup(),
			NewBatteryGroup(),
			NewBMSGroup(),
			NewGenGroup(),
			NewLoadGroup(),
			NewHouseGroup(),
		},
	}
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.MustRegister(e.deviceUp, e.lastUpdate)
	for _, g := range e.groups {
		prometheus.MustRegister(g.Gauge)
	}

	if cfg.EnableGeneric {
		e.generic = newGroup("solarman_metric", "All currentData numeric metrics (labels: device_sn,key,name,unit).")
		prometheus.MustRegister(e.generic)
	}

	return e
}

func (e *Exporter) Ready() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if !e.ready {
		return false
	}
	return time.Since(e.lastSuccess) <= 3*e.cfg.PollInterval
}

func (e *Exporter) Run(ctx context.Context) {
	t := time.NewTicker(e.cfg.PollInterval)
	defer t.Stop()

	e.poll(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("exporter stopped")
			return
		case <-t.C:
			e.poll(ctx)
		}
	}
}

func (e *Exporter) poll(ctx context.Context) {
	deviceSNs := e.cfg.DeviceSNs
	if len(deviceSNs) == 0 {
		sns, err := e.discoverDeviceSNs(ctx)
		if err != nil {
			log.Error().Err(err).Msg("auto-discovery failed")
			e.setNotReady()
			return
		}
		deviceSNs = sns
		if len(deviceSNs) == 0 {
			log.Error().Msg("no devices found")
			e.setNotReady()
			return
		}
	}

	for _, sn := range deviceSNs {
		sn = strings.TrimSpace(sn)
		if sn == "" {
			continue
		}

		points, err := e.cfg.Client.CurrentData(ctx, sn)
		if err != nil {
			log.Warn().Err(err).Str("device_sn", sn).Msg("poll failed")
			e.deviceUp.WithLabelValues(sn).Set(0)
			continue
		}

		e.deviceUp.WithLabelValues(sn).Set(1)
		e.lastUpdate.WithLabelValues(sn).Set(float64(time.Now().Unix()))

		e.updateGroups(sn, points)
		e.setReady()
	}
}

func (e *Exporter) updateGroups(deviceSN string, points []models.MetricPoint) {
	for _, p := range points {
		if e.generic != nil {
			e.generic.WithLabelValues(deviceSN, p.Key, p.Name, p.Unit).Set(p.Value)
		}
		for _, g := range e.groups {
			if g.Match != nil && g.Match(p) {
				g.Gauge.WithLabelValues(deviceSN, p.Key, p.Name, p.Unit).Set(p.Value)
			}
		}
	}
}

func (e *Exporter) discoverDeviceSNs(ctx context.Context) ([]string, error) {
	stations, err := e.cfg.Client.Stations(ctx)
	if err != nil {
		return nil, err
	}
	if len(stations) == 0 {
		return nil, nil
	}

	stationID := e.cfg.StationID
	if stationID == 0 {
		stationID = stations[0].ID
		log.Info().Int64("station_id", stationID).Str("station_name", stations[0].Name).Msg("using first station")
	} else {
		log.Info().Int64("station_id", stationID).Msg("using configured station")
	}

	devs, err := e.cfg.Client.StationDevices(ctx, stationID)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(devs))
	for _, d := range devs {
		if d.DeviceSN != "" {
			out = append(out, d.DeviceSN)
		}
	}
	log.Info().Int("devices", len(out)).Msg("auto-discovery complete")
	return out, nil
}

func (e *Exporter) setReady() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ready = true
	e.lastSuccess = time.Now()
}

func (e *Exporter) setNotReady() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ready = false
}
