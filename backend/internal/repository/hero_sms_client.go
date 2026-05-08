package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

// heroSmsRetryBackoffs is the wait duration before each subsequent attempt of
// GetNumberWithRetry. With these 5 entries, GetNumberWithRetry performs up to
// 6 attempts (1 initial + 5 retries), waiting 1/3/5/8/10 seconds between them
// — matching the configured retry policy.
var heroSmsRetryBackoffs = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	5 * time.Second,
	8 * time.Second,
	10 * time.Second,
}

type heroSmsClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	service    string
	country    string
}

func NewHeroSmsClient(cfg *config.Config) service.HeroSmsClient {
	return &heroSmsClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		baseURL:    cfg.HeroSms.BaseURL,
		apiKey:     cfg.HeroSms.APIKey,
		service:    cfg.HeroSms.Service,
		country:    cfg.HeroSms.Country,
	}
}

func (c *heroSmsClient) GetNumber(ctx context.Context) (*service.HeroSmsNumber, error) {
	params := url.Values{}
	params.Set("action", "getNumber")
	params.Set("service", c.service)
	params.Set("country", c.country)
	params.Set("api_key", c.apiKey)

	reqURL := c.baseURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hero-sms getNumber: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	text := strings.TrimSpace(string(body))
	// Expected format: ACCESS_NUMBER:id:phone
	parts := strings.SplitN(text, ":", 3)
	if len(parts) != 3 || parts[0] != "ACCESS_NUMBER" {
		return nil, fmt.Errorf("unexpected getNumber response: %s", text)
	}

	return &service.HeroSmsNumber{
		ID:    parts[1],
		Phone: parts[2],
	}, nil
}

// GetNumberWithRetry retries GetNumber on failure with the configured backoff.
// All hero-sms errors (network and business such as NO_NUMBERS / BAD_KEY) are
// treated uniformly because hero-sms occasionally returns business errors
// while it warms up.
func (c *heroSmsClient) GetNumberWithRetry(ctx context.Context) (*service.HeroSmsNumber, error) {
	var lastErr error
	for attempt := 0; attempt <= len(heroSmsRetryBackoffs); attempt++ {
		if attempt > 0 {
			wait := heroSmsRetryBackoffs[attempt-1]
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}

		num, err := c.GetNumber(ctx)
		if err == nil {
			return num, nil
		}
		lastErr = err
		slog.Warn("hero-sms getNumber attempt failed",
			"attempt", attempt+1,
			"max_attempts", len(heroSmsRetryBackoffs)+1,
			"error", err,
		)
	}
	return nil, fmt.Errorf("hero-sms getNumber failed after %d attempts: %w", len(heroSmsRetryBackoffs)+1, lastErr)
}

func (c *heroSmsClient) GetStatus(ctx context.Context, id string) (*service.HeroSmsStatus, error) {
	params := url.Values{}
	params.Set("action", "getStatusV2")
	params.Set("id", id)
	params.Set("api_key", c.apiKey)

	reqURL := c.baseURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hero-sms getStatus: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result struct {
		VerificationType int `json:"verificationType"`
		Sms              *struct {
			DateTime string `json:"dateTime"`
			Code     string `json:"code"`
			Text     string `json:"text"`
		} `json:"sms"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse getStatus response: %w (body: %s)", err, string(body))
	}

	status := &service.HeroSmsStatus{}
	if result.Sms != nil && result.Sms.Code != "" {
		status.Received = true
		status.Code = result.Sms.Code
		status.Text = result.Sms.Text
	}
	return status, nil
}
