package solarman

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type tokenResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`

	AccessToken          string `json:"access_token"`
	TokenType            string `json:"token_type"`
	RefreshToken         string `json:"refresh_token"`
	ExpiresIn            int64  `json:"expires_in,string"`
	AccessTokenExpiresAt time.Time
}

func (c *Client) authHeader() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	tt := strings.TrimSpace(c.token.TokenType)
	if tt == "" {
		tt = "Bearer"
	}
	return tt + " " + strings.TrimSpace(c.token.AccessToken)
}

func (c *Client) hasToken() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token.AccessToken != ""
}

func (c *Client) passwordSHA256Hex() (string, error) {
	if strings.TrimSpace(c.cfg.PasswordSHA) != "" {
		return strings.ToLower(strings.TrimSpace(c.cfg.PasswordSHA)), nil
	}
	if c.cfg.Password == "" {
		return "", errors.New("missing password")
	}
	sum := sha256.Sum256([]byte(c.cfg.Password))
	return hex.EncodeToString(sum[:]), nil
}

func (c *Client) EnsureToken(ctx context.Context) error {
	if time.Now().After(c.token.AccessTokenExpiresAt) && (c.token.RefreshToken != "") {
		return c.obtainToken(ctx)
	}
	if c.hasToken() {
		return nil
	}
	return c.obtainToken(ctx)
}

func (c *Client) obtainToken(ctx context.Context) error {
	passHex, err := c.passwordSHA256Hex()
	if err != nil {
		return err
	}

	body := map[string]any{
		"appSecret": c.cfg.AppSecret,
		"email":     c.cfg.Email,
		"password":  passHex,
	}

	raw, status, err := c.doJSON(ctx, "POST", fmt.Sprintf("/account/%s/token", c.cfg.APIVersion), true, false, body)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("token request failed: status=%d body=%s", status, strings.TrimSpace(string(raw)))
	}
	var tr tokenResponse
	if err := json.Unmarshal(raw, &tr); err != nil {
		return fmt.Errorf("token Unmarshal failed: %w", err)
	}
	if !tr.Success || tr.AccessToken == "" {
		return fmt.Errorf("token error: %s", tr.Msg)
	}
	tr.AccessTokenExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn-5) * time.Second)
	c.mu.Lock()
	c.token = tr
	c.mu.Unlock()

	log.Info().Msg("obtained solarman access token")
	return nil
}
