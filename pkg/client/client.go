// Package client provides an HTTP client for the DeepCoin REST API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/deepcoinapi/agent-cli/pkg/auth"
)

// Client is a synchronous HTTP client for DeepCoin's REST API.
type Client struct {
	APIKey     string
	SecretKey  string
	Passphrase string
	BaseURL    string
	http       *http.Client
}

// New creates a Client from environment variables or explicit values.
func New() *Client {
	base := envOrAny("https://api.deepcoin.com", "DEEPCOIN_BASE_URL", "DC_BASE_URL")
	return &Client{
		APIKey:     envOrAny("", "DEEPCOIN_API_KEY", "DC_API_KEY"),
		SecretKey:  envOrAny("", "DEEPCOIN_SECRET_KEY", "DC_SECRET_KEY"),
		Passphrase: envOrAny("", "DEEPCOIN_PASSPHRASE", "DC_PASSPHRASE"),
		BaseURL:    strings.TrimRight(base, "/"),
		http:       &http.Client{Timeout: 30 * time.Second},
	}
}

func envOrAny(fallback string, keys ...string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}
	return fallback
}

// ── Public (no auth) ────────────────────────────────────────────────

// GetPublic sends an unauthenticated GET request.
func (c *Client) GetPublic(path string, params map[string]string) (map[string]any, error) {
	u := c.buildURL(path, params)
	resp, err := c.http.Get(u)
	if err != nil {
		return nil, err
	}
	return parseResponse(resp)
}

// ── Private (authenticated) ─────────────────────────────────────────

// Get sends an authenticated GET request.
func (c *Client) Get(path string, params map[string]string) (map[string]any, error) {
	qs := buildQuery(params)
	signPath := path
	if qs != "" {
		signPath = path + "?" + qs
	}

	ts := auth.Timestamp()
	sig := auth.Sign(ts, "GET", signPath, "", c.SecretKey)

	u := c.BaseURL + signPath
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	c.setAuthHeaders(req, ts, sig)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return parseResponse(resp)
}

// Post sends an authenticated POST request with a JSON body.
func (c *Client) Post(path string, body map[string]any) (map[string]any, error) {
	var bodyStr string
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyStr = string(b)
		bodyReader = bytes.NewReader(b)
	} else {
		bodyReader = strings.NewReader("")
	}

	ts := auth.Timestamp()
	sig := auth.Sign(ts, "POST", path, bodyStr, c.SecretKey)

	u := c.BaseURL + path
	req, err := http.NewRequest("POST", u, bodyReader)
	if err != nil {
		return nil, err
	}
	c.setAuthHeaders(req, ts, sig)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return parseResponse(resp)
}

// ── Helpers ─────────────────────────────────────────────────────────

func (c *Client) setAuthHeaders(req *http.Request, ts, sig string) {
	req.Header.Set("DC-ACCESS-KEY", c.APIKey)
	req.Header.Set("DC-ACCESS-SIGN", sig)
	req.Header.Set("DC-ACCESS-TIMESTAMP", ts)
	req.Header.Set("DC-ACCESS-PASSPHRASE", c.Passphrase)
	req.Header.Set("Content-Type", "application/json")
}

func (c *Client) buildURL(path string, params map[string]string) string {
	qs := buildQuery(params)
	if qs != "" {
		return c.BaseURL + path + "?" + qs
	}
	return c.BaseURL + path
}

func buildQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	v := url.Values{}
	for k, val := range params {
		if val != "" {
			v.Set(k, val)
		}
	}
	return v.Encode()
}

func parseResponse(resp *http.Response) (map[string]any, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode >= 400 {
		code := result["code"]
		msg := result["msg"]
		return nil, fmt.Errorf("API error %v: %v", code, msg)
	}

	if code, ok := result["code"]; ok && fmt.Sprintf("%v", code) != "0" {
		return nil, fmt.Errorf("API error %v: %v", code, result["msg"])
	}

	return result, nil
}

// GetData is a convenience that extracts the "data" field from the response.
func GetData(resp map[string]any) any {
	if d, ok := resp["data"]; ok {
		return d
	}
	return resp
}

// GetDataSlice extracts "data" as []map[string]any.
func GetDataSlice(resp map[string]any) []map[string]any {
	d, ok := resp["data"]
	if !ok {
		return nil
	}
	switch v := d.(type) {
	case []any:
		out := make([]map[string]any, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	case []map[string]any:
		return v
	}
	return nil
}

// GetDataMap extracts "data" as map[string]any.
func GetDataMap(resp map[string]any) map[string]any {
	d, ok := resp["data"]
	if !ok {
		return nil
	}
	if m, ok := d.(map[string]any); ok {
		return m
	}
	return nil
}
