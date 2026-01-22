package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"solarman_exporter/internal/config"
	"solarman_exporter/internal/exporter"
	"solarman_exporter/internal/solarman"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Run(argv []string) {
	app := &cli.App{
		Name:  "Solarman Exporter",
		Usage: "Prometheus exporter for Solarman Smart (Solarman OpenAPI)",
		Flags: config.Flags(),
		Action: func(c *cli.Context) error {
			cfg, err := config.FromCLI(c)
			if err != nil {
				return err
			}

			zerolog.TimeFieldFormat = time.RFC3339Nano
			log.Logger = log.Level(cfg.LogLevel)

			client := solarman.NewClient(solarman.Config{
				BaseURL:     cfg.BaseURL,
				APIVersion:  cfg.APIVersion,
				Language:    cfg.Language,
				AppID:       cfg.AppID,
				AppSecret:   cfg.AppSecret,
				Email:       cfg.Email,
				Password:    cfg.Password,
				PasswordSHA: cfg.PasswordSHA256,
				HTTPTimeout: cfg.HTTPTimeout,
			})

			exp := exporter.New(exporter.Config{
				Client:        client,
				PollInterval:  cfg.PollInterval,
				DeviceSNs:     cfg.DeviceSN,
				StationID:     cfg.StationID,
				EnableGeneric: cfg.EnableGeneric,
			})

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			go exp.Run(ctx)

			mux := http.NewServeMux()
			mux.Handle(cfg.MetricsPath, promhttp.Handler())
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_, _ = w.Write([]byte("solarman_exporter\n"))
				_, _ = w.Write([]byte("GET " + cfg.MetricsPath + " for Prometheus metrics\n"))
			})
			mux.HandleFunc("/-/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok\n"))
			})
			mux.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
				if exp.Ready() {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("ready\n"))
					return
				}
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("not ready\n"))
			})

			srv := &http.Server{
				Addr:              cfg.Listen,
				Handler:           mux,
				ReadHeaderTimeout: 5 * time.Second,
			}

			go func() {
				<-ctx.Done()
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_ = srv.Shutdown(shutdownCtx)
			}()

			log.Info().
				Str("listen", cfg.Listen).
				Str("metrics_path", cfg.MetricsPath).
				Dur("poll_interval", cfg.PollInterval).
				Msg("Starting Solarman Exporter")

			err = srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Error().Err(err).Msg("server error")
				return err
			}
			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		log.Fatal().Err(err).Msg("fatal")
	}
}
