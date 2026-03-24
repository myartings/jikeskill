package jike

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/myartings/jikeskill/tokens"
)

const (
	BaseURL = "https://api.ruguoapp.com"
	Origin  = "https://web.okjike.com"
)

type Client struct {
	httpClient *http.Client
	store      *tokens.Store
	mu         sync.Mutex
}

func NewClient(store *tokens.Store) *Client {
	return &Client{
		httpClient: &http.Client{},
		store:      store,
	}
}

func (c *Client) Store() *tokens.Store {
	return c.store
}

// Do sends an authenticated request. Automatically refreshes token on 401.
func (c *Client) Do(method, path string, body any) ([]byte, http.Header, error) {
	respBody, respHeader, statusCode, err := c.doOnce(method, path, body)
	if err != nil {
		return nil, nil, err
	}
	if statusCode == http.StatusUnauthorized {
		if refreshErr := c.refreshToken(); refreshErr != nil {
			return nil, nil, fmt.Errorf("token refresh failed: %w", refreshErr)
		}
		respBody, respHeader, statusCode, err = c.doOnce(method, path, body)
		if err != nil {
			return nil, nil, err
		}
	}
	if statusCode >= 400 {
		return nil, nil, fmt.Errorf("API error %d: %s", statusCode, string(respBody))
	}
	return respBody, respHeader, nil
}

// DoRaw sends a request without auth (for login endpoints).
func (c *Client) DoRaw(method, path string, body any) ([]byte, http.Header, int, error) {
	return c.doRequest(method, path, body, "")
}

func (c *Client) doOnce(method, path string, body any) ([]byte, http.Header, int, error) {
	token := ""
	if td := c.store.Get(); td != nil {
		token = td.AccessToken
	}
	return c.doRequest(method, path, body, token)
}

func (c *Client) doRequest(method, path string, body any, token string) ([]byte, http.Header, int, error) {
	url := BaseURL + path

	var reqBody io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", Origin)
	req.Header.Set("Referer", Origin+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	if token != "" {
		req.Header.Set("x-jike-access-token", token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.Header, resp.StatusCode, nil
}

func (c *Client) refreshToken() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	td := c.store.Get()
	if td == nil || td.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	url := BaseURL + "/app_auth_tokens.refresh"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", Origin)
	req.Header.Set("x-jike-refresh-token", td.RefreshToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refresh failed %d: %s", resp.StatusCode, string(body))
	}

	newAccess := resp.Header.Get("x-jike-access-token")
	newRefresh := resp.Header.Get("x-jike-refresh-token")
	if newAccess == "" {
		return fmt.Errorf("no access token in refresh response")
	}

	newTD := &tokens.TokenData{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}
	if newTD.RefreshToken == "" {
		newTD.RefreshToken = td.RefreshToken
	}
	return c.store.Save(newTD)
}
