package jike

import (
	"context"
	"encoding/json"
	"fmt"
)

// Search performs an integrated search.
func (c *Client) Search(ctx context.Context, keyword string, loadMoreKey any) (*SearchResult, error) {
	reqBody := map[string]any{
		"keywords": keyword,
		"type":     "ALL",
	}
	if loadMoreKey != nil {
		reqBody["loadMoreKey"] = loadMoreKey
	}

	body, _, err := c.Do("POST", "/1.0/search/integrate", reqBody)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	var resp SearchResult
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse search result: %w", err)
	}
	return &resp, nil
}
