package jike

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
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

// parseCreatedAt tries common Jike timestamp formats.
func parseCreatedAt(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05.000Z"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// GetTopicFeedPages fetches multiple pages of topic posts.
// maxPosts=0 means no post-count limit.
// since=0 means no time filter.
// Posts are assumed newest-first; stops early when a post is older than cutoff.
func (c *Client) GetTopicFeedPages(ctx context.Context, topicID string, maxPosts int, since time.Duration) ([]Post, error) {
	var allPosts []Post
	var loadMoreKey any
	var cutoff time.Time
	if since > 0 {
		cutoff = time.Now().Add(-since)
	}

	for {
		resp, err := c.GetTopicFeed(ctx, topicID, loadMoreKey)
		if err != nil {
			return nil, err
		}
		for _, p := range resp.Data {
			if since > 0 {
				if t, ok := parseCreatedAt(p.CreatedAt); ok && t.Before(cutoff) {
					return allPosts, nil
				}
			}
			allPosts = append(allPosts, p)
			if maxPosts > 0 && len(allPosts) >= maxPosts {
				return allPosts, nil
			}
		}
		if resp.LoadMoreKey == nil {
			break
		}
		loadMoreKey = resp.LoadMoreKey
	}
	return allPosts, nil
}
