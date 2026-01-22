package solarman

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) buildURL(path string, withAppLang bool) (string, error) {
	base := strings.TrimRight(c.cfg.BaseURL, "/")
	p := strings.TrimLeft(path, "/")
	u, err := url.Parse(base + "/" + p)
	if err != nil {
		return "", err
	}
	if withAppLang {
		q := u.Query()
		q.Set("appId", c.cfg.AppID)
		q.Set("language", c.cfg.Language)
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

// doJSON returns (rawBody, statusCode, error).
// withAppLang -> add appId & language query params.
// withAuth -> add Authorization header.
func (c *Client) doJSON(ctx context.Context, method, path string, withAppLang bool, withAuth bool, body any) ([]byte, int, error) {
	u, err := c.buildURL(path, withAppLang)
	if err != nil {
		return nil, 0, err
	}

	var rdr io.Reader
	if body == nil {
		rdr = bytes.NewReader([]byte(`{}`))
	} else {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		rdr = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, rdr)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if withAuth {
		req.Header.Set("Authorization", c.authHeader())
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	raw, _ := io.ReadAll(resp.Body)
	return raw, resp.StatusCode, nil
}

// doJSONAuthRetry retries once on 401 by refreshing token.
func (c *Client) doJSONAuthRetry(ctx context.Context, method, path string, withAppLang bool, body any) ([]byte, int, error) {
	if err := c.EnsureToken(ctx); err != nil {
		return nil, 0, err
	}

	raw, status, err := c.doJSON(ctx, method, path, withAppLang, true, body)
	if err != nil {
		return nil, 0, err
	}
	if status != http.StatusUnauthorized {
		return raw, status, nil
	}

	// refresh token and retry once
	if err := c.obtainToken(ctx); err != nil {
		return raw, status, fmt.Errorf("401 then token refresh failed: %w", err)
	}
	return c.doJSON(ctx, method, path, withAppLang, true, body)
}
