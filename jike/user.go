package jike

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetUserProfile gets a user's profile by username.
func (c *Client) GetUserProfile(ctx context.Context, username string) (*User, error) {
	path := fmt.Sprintf("/1.0/users/profile?username=%s", username)
	body, _, err := c.Do("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}

	var resp struct {
		User User `json:"user"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse user profile: %w", err)
	}
	return &resp.User, nil
}

// GetUserPosts gets posts by a specific user.
func (c *Client) GetUserPosts(ctx context.Context, username string, loadMoreKey any) (*FeedResponse, error) {
	reqBody := map[string]any{
		"username": username,
	}
	if loadMoreKey != nil {
		reqBody["loadMoreKey"] = loadMoreKey
	}

	body, _, err := c.Do("POST", "/1.0/personalUpdate/single", reqBody)
	if err != nil {
		return nil, fmt.Errorf("get user posts: %w", err)
	}

	var resp FeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse user posts: %w", err)
	}
	return &resp, nil
}
