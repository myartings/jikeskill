package jike

import (
	"context"
	"fmt"
)

// LikePost likes a post.
func (c *Client) LikePost(ctx context.Context, postID, targetType string) error {
	if targetType == "" {
		targetType = "ORIGINAL_POST"
	}
	reqBody := map[string]any{
		"targetId":   postID,
		"targetType": targetType,
	}
	_, _, err := c.Do("POST", "/1.0/likes/save", reqBody)
	if err != nil {
		return fmt.Errorf("like post: %w", err)
	}
	return nil
}

// UnlikePost unlikes a post.
func (c *Client) UnlikePost(ctx context.Context, postID, targetType string) error {
	if targetType == "" {
		targetType = "ORIGINAL_POST"
	}
	reqBody := map[string]any{
		"targetId":   postID,
		"targetType": targetType,
	}
	_, _, err := c.Do("POST", "/1.0/likes/remove", reqBody)
	if err != nil {
		return fmt.Errorf("unlike post: %w", err)
	}
	return nil
}

// FollowUser follows a user.
func (c *Client) FollowUser(ctx context.Context, username string) error {
	reqBody := map[string]any{
		"username": username,
	}
	_, _, err := c.Do("POST", "/1.0/userRelation/follow", reqBody)
	if err != nil {
		return fmt.Errorf("follow user: %w", err)
	}
	return nil
}

// UnfollowUser unfollows a user.
func (c *Client) UnfollowUser(ctx context.Context, username string) error {
	reqBody := map[string]any{
		"username": username,
	}
	_, _, err := c.Do("POST", "/1.0/userRelation/unfollow", reqBody)
	if err != nil {
		return fmt.Errorf("unfollow user: %w", err)
	}
	return nil
}
