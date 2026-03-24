package jike

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetPostDetail gets a post by its ID.
func (c *Client) GetPostDetail(ctx context.Context, postID, postType string) (*Post, error) {
	if postType == "" {
		postType = "ORIGINAL_POST"
	}

	var path string
	switch postType {
	case "REPOST":
		path = "/1.0/reposts/get"
	default:
		path = "/1.0/originalPosts/get"
	}

	reqBody := map[string]any{
		"id": postID,
	}

	body, _, err := c.Do("POST", path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("get post detail: %w", err)
	}

	var resp struct {
		Data Post `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse post: %w", err)
	}
	return &resp.Data, nil
}

// CreatePost creates a new original post.
func (c *Client) CreatePost(ctx context.Context, content string, topicID string, pictureKeys []string) (*Post, error) {
	reqBody := map[string]any{
		"content": content,
	}
	if topicID != "" {
		reqBody["topicId"] = topicID
	}
	if len(pictureKeys) > 0 {
		reqBody["pictureKeys"] = pictureKeys
	}

	body, _, err := c.Do("POST", "/1.0/originalPosts/create", reqBody)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	var resp struct {
		Data Post `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse post: %w", err)
	}
	return &resp.Data, nil
}

// RemovePost removes a post by its ID.
func (c *Client) RemovePost(ctx context.Context, postID string) error {
	reqBody := map[string]any{
		"id": postID,
	}
	_, _, err := c.Do("POST", "/1.0/originalPosts/remove", reqBody)
	if err != nil {
		return fmt.Errorf("remove post: %w", err)
	}
	return nil
}
