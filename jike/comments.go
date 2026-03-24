package jike

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetComments gets comments for a post.
func (c *Client) GetComments(ctx context.Context, targetID, targetType string, loadMoreKey any) (*CommentResponse, error) {
	reqBody := map[string]any{
		"targetId":   targetID,
		"targetType": targetType,
	}
	if loadMoreKey != nil {
		reqBody["loadMoreKey"] = loadMoreKey
	}

	body, _, err := c.Do("POST", "/1.0/comments/listPrimary", reqBody)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}

	var resp CommentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse comments: %w", err)
	}
	return &resp, nil
}

// AddComment adds a comment to a post.
func (c *Client) AddComment(ctx context.Context, targetID, targetType, content string) (*Comment, error) {
	reqBody := map[string]any{
		"targetId":   targetID,
		"targetType": targetType,
		"content":    content,
	}

	body, _, err := c.Do("POST", "/1.0/comments/add", reqBody)
	if err != nil {
		return nil, fmt.Errorf("add comment: %w", err)
	}

	var resp struct {
		Data Comment `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse comment: %w", err)
	}
	return &resp.Data, nil
}
