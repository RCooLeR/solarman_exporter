package config

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Config struct {
	Listen      string
	MetricsPath string
	LogLevel    zerolog.Level

	BaseURL    string
	APIVersion string
	Language   string

	AppID          string
	AppSecret      string
	Email          string
	Password       string
	PasswordSHA256 string

	DeviceSN     []string
	StationID    int64
	PollInterval time.Duration
	HTTPTimeout  time.Duration

	EnableGeneric bool
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "listen", Value: ":9206", EnvVars: []string{"SOLARMAN_EXPORTER_LISTEN"}},
		&cli.StringFlag{Name: "metrics-path", Value: "/metrics", EnvVars: []string{"SOLARMAN_EXPORTER_METRICS_PATH"}},
		&cli.StringFlag{Name: "log-level", Value: "info", EnvVars: []string{"SOLARMAN_EXPORTER_LOG_LEVEL"}},

		&cli.StringFlag{Name: "base-url", Value: "https://globalapi.solarmanpv.com", EnvVars: []string{"SOLARMAN_BASE_URL"}},
		&cli.StringFlag{Name: "api-version", Value: "v1.0", EnvVars: []string{"SOLARMAN_API_VERSION"}},
		&cli.StringFlag{Name: "language", Value: "en", EnvVars: []string{"SOLARMAN_LANGUAGE"}},

		&cli.StringFlag{Name: "app-id", EnvVars: []string{"SOLARMAN_APP_ID"}, Required: true},
		&cli.StringFlag{Name: "app-secret", EnvVars: []string{"SOLARMAN_APP_SECRET"}, Required: true},
		&cli.StringFlag{Name: "email", EnvVars: []string{"SOLARMAN_EMAIL"}, Required: true},
		&cli.StringFlag{Name: "password", EnvVars: []string{"SOLARMAN_PASSWORD"}},
		&cli.StringFlag{Name: "password-sha256", EnvVars: []string{"SOLARMAN_PASSWORD_SHA256"}},

		&cli.StringSliceFlag{Name: "device-sn", EnvVars: []string{"SOLARMAN_DEVICE_SN"}},
		&cli.Int64Flag{Name: "station-id", Value: 0, EnvVars: []string{"SOLARMAN_STATION_ID"}},

		&cli.DurationFlag{Name: "poll-interval", Value: 30 * time.Second, EnvVars: []string{"SOLARMAN_POLL_INTERVAL"}},
		&cli.DurationFlag{Name: "http-timeout", Value: 20 * time.Second, EnvVars: []string{"SOLARMAN_HTTP_TIMEOUT"}},

		&cli.BoolFlag{Name: "enable-generic", Value: true, EnvVars: []string{"SOLARMAN_ENABLE_GENERIC"}},
	}
}

func FromCLI(c *cli.Context) (Config, error) {
	lvl, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(c.String("log-level"))))
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	cfg := Config{
		Listen:      c.String("listen"),
		MetricsPath: c.String("metrics-path"),
		LogLevel:    lvl,

		BaseURL:    c.String("base-url"),
		APIVersion: c.String("api-version"),
		Language:   c.String("language"),

		AppID:          c.String("app-id"),
		AppSecret:      c.String("app-secret"),
		Email:          c.String("email"),
		Password:       c.String("password"),
		PasswordSHA256: c.String("password-sha256"),

		DeviceSN:     c.StringSlice("device-sn"),
		StationID:    c.Int64("station-id"),
		PollInterval: c.Duration("poll-interval"),
		HTTPTimeout:  c.Duration("http-timeout"),

		EnableGeneric: c.Bool("enable-generic"),
	}

	if cfg.Password == "" && cfg.PasswordSHA256 == "" {
		log.Error().Err(errors.New("either --password or --password-sha256 must be provided")).Msg("password is required")
		return Config{}, errors.New("either --password or --password-sha256 must be provided")
	}
	return cfg, nil
}
