package jike

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/myartings/jikeskill/tokens"
	qrcode "github.com/skip2/go-qrcode"
)

// CreateSession creates a new login session and returns the UUID.
func (c *Client) CreateSession(ctx context.Context) (string, error) {
	body, _, _, err := c.DoRaw("POST", "/sessions.create", nil)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	var resp SessionCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parse session response: %w", err)
	}
	if resp.UUID == "" {
		return "", fmt.Errorf("empty uuid in response: %s", string(body))
	}
	return resp.UUID, nil
}

// GenerateQRCode generates a QR code PNG image (base64 encoded) for the given session UUID.
func GenerateQRCode(uuid string) (string, error) {
	scanURL := fmt.Sprintf("https://web.okjike.com/scan-login?uuid=%s", uuid)

	png, err := qrcode.Encode(scanURL, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("generate qrcode: %w", err)
	}
	return base64.StdEncoding.EncodeToString(png), nil
}

// WaitForLogin polls for login confirmation. Returns the logged-in user on success.
// Timeout is 180 seconds.
func (c *Client) WaitForLogin(ctx context.Context, uuid string) (*User, error) {
	timeout := time.After(180 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("login timeout after 180 seconds")
		case <-ticker.C:
			user, done, err := c.checkConfirmation(uuid)
			if err != nil {
				continue // keep polling on errors
			}
			if done {
				return user, nil
			}
		}
	}
}

func (c *Client) checkConfirmation(uuid string) (*User, bool, error) {
	path := fmt.Sprintf("/sessions.wait_for_confirmation?uuid=%s", uuid)
	body, header, statusCode, err := c.DoRaw("GET", path, nil)
	if err != nil {
		return nil, false, err
	}

	if statusCode != 200 {
		return nil, false, fmt.Errorf("status %d", statusCode)
	}

	// Save tokens from response headers
	accessToken := header.Get("x-jike-access-token")
	refreshToken := header.Get("x-jike-refresh-token")
	if accessToken == "" {
		return nil, false, fmt.Errorf("no access token in response")
	}

	if err := c.store.Save(&tokens.TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}); err != nil {
		return nil, false, fmt.Errorf("save tokens: %w", err)
	}

	var resp LoginConfirmResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, true, nil // logged in but can't parse user
	}

	return &resp.User, true, nil
}

// CheckLoginStatus checks if the current tokens are valid.
func (c *Client) CheckLoginStatus(ctx context.Context) (bool, *User, error) {
	td := c.store.Get()
	if td == nil || td.AccessToken == "" {
		return false, nil, nil
	}

	body, _, err := c.Do("GET", "/1.0/users/profile", nil)
	if err != nil {
		return false, nil, nil
	}

	var resp struct {
		User User `json:"user"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, nil, nil
	}

	return true, &resp.User, nil
}
