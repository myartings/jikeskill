package jike

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetTopicFeed gets posts from a specific topic (圈子).
func (c *Client) GetTopicFeed(ctx context.Context, topicID string, loadMoreKey any) (*FeedResponse, error) {
	reqBody := map[string]any{
		"topicId": topicID,
	}
	if loadMoreKey != nil {
		reqBody["loadMoreKey"] = loadMoreKey
	}

	body, _, err := c.Do("POST", "/1.0/topicFeed/list", reqBody)
	if err != nil {
		return nil, fmt.Errorf("get topic feed: %w", err)
	}

	var resp FeedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse topic feed: %w", err)
	}
	return &resp, nil
}
