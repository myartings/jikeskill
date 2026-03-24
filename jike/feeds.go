package jike

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetFollowingFeeds returns the user's following timeline.
func (c *Client) GetFollowingFeeds(ctx context.Context, loadMoreKey any) (*FeedResponse, error) {
	reqBody := map[string]any{}
	if loadMoreKey != nil {
		reqBody["loadMoreKey"] = loadMoreKey
	}

	body, _, err := c.Do("POST", "/1.0/personalUpdate/followingUpdates", reqBody)
	if err != nil {
		return nil, fmt.Errorf("get following feeds: %w", err)
	}

	var resp FeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse feeds: %w", err)
	}
	return &resp, nil
}

// GetRecommendFeeds returns recommended feeds.
func (c *Client) GetRecommendFeeds(ctx context.Context, loadMoreKey any) (*FeedResponse, error) {
	reqBody := map[string]any{}
	if loadMoreKey != nil {
		reqBody["loadMoreKey"] = loadMoreKey
	}

	body, _, err := c.Do("POST", "/1.0/recommendFeed/list", reqBody)
	if err != nil {
		return nil, fmt.Errorf("get recommend feeds: %w", err)
	}

	var resp FeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse feeds: %w", err)
	}
	return &resp, nil
}
