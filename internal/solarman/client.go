package solarman

import (
	"net/http"
	"sync"
	"time"
)

type Config struct {
	BaseURL            string
	APIVersion         string
	Language           string
	AppID              string
	AppSecret          string
	Email              string
	Password           string
	PasswordSHA        string
	HTTPTimeout        time.Duration
	YearlyRequestLimit int64
}

type Client struct {
	cfg Config
	hc  *http.Client

	mu             sync.Mutex
	token          tokenResponse
	requestSpacing time.Duration
	nextRequestAt  time.Time
}

func NewClient(cfg Config) *Client {
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = 15 * time.Second
	}
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v1.0"
	}
	var requestSpacing time.Duration
	if cfg.YearlyRequestLimit > 0 {
		requestSpacing = (365 * 24 * time.Hour) / time.Duration(cfg.YearlyRequestLimit)
	}
	return &Client{
		cfg:            cfg,
		hc:             &http.Client{Timeout: cfg.HTTPTimeout},
		requestSpacing: requestSpacing,
	}
}
