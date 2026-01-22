package solarman

import (
	"net/http"
	"sync"
	"time"
)

type Config struct {
	BaseURL     string
	APIVersion  string
	Language    string
	AppID       string
	AppSecret   string
	Email       string
	Password    string
	PasswordSHA string
	HTTPTimeout time.Duration
}

type Client struct {
	cfg Config
	hc  *http.Client

	mu    sync.Mutex
	token tokenResponse
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
	return &Client{
		cfg: cfg,
		hc:  &http.Client{Timeout: cfg.HTTPTimeout},
	}
}
