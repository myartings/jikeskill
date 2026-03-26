package main

import (
	"context"

	"github.com/myartings/jikeskill/jike"
	"github.com/myartings/jikeskill/tokens"
)

type JikeService struct {
	client *jike.Client
}

func NewJikeService(tokenPath string) *JikeService {
	store := tokens.NewStore(tokenPath)
	client := jike.NewClient(store)
	return &JikeService{client: client}
}

// Login flow
func (s *JikeService) CreateLoginSession(ctx context.Context) (uuid string, qrcodeBase64 string, err error) {
	uuid, err = s.client.CreateSession(ctx)
	if err != nil {
		return "", "", err
	}
	qr, err := jike.GenerateQRCode(uuid)
	if err != nil {
		return "", "", err
	}
	return uuid, qr, nil
}

func (s *JikeService) WaitForLogin(ctx context.Context, uuid string) (*jike.User, error) {
	return s.client.WaitForLogin(ctx, uuid)
}

func (s *JikeService) CheckLoginStatus(ctx context.Context) (bool, *jike.User, error) {
	return s.client.CheckLoginStatus(ctx)
}

func (s *JikeService) Logout() error {
	return s.client.Store().Delete()
}

// Feeds
func (s *JikeService) GetFollowingFeeds(ctx context.Context, loadMoreKey any) (*jike.FeedResponse, error) {
	return s.client.GetFollowingFeeds(ctx, loadMoreKey)
}

func (s *JikeService) GetRecommendFeeds(ctx context.Context, loadMoreKey any) (*jike.FeedResponse, error) {
	return s.client.GetRecommendFeeds(ctx, loadMoreKey)
}

// Search
func (s *JikeService) Search(ctx context.Context, keyword string, loadMoreKey any) (*jike.SearchResult, error) {
	return s.client.Search(ctx, keyword, loadMoreKey)
}

// Posts
func (s *JikeService) GetPostDetail(ctx context.Context, postID, postType string) (*jike.Post, error) {
	return s.client.GetPostDetail(ctx, postID, postType)
}

func (s *JikeService) CreatePost(ctx context.Context, content, topicID string, pictureKeys []string) (*jike.Post, error) {
	return s.client.CreatePost(ctx, content, topicID, pictureKeys)
}

func (s *JikeService) RemovePost(ctx context.Context, postID string) error {
	return s.client.RemovePost(ctx, postID)
}

// Comments
func (s *JikeService) GetComments(ctx context.Context, targetID, targetType string, loadMoreKey any) (*jike.CommentResponse, error) {
	return s.client.GetComments(ctx, targetID, targetType, loadMoreKey)
}

func (s *JikeService) AddComment(ctx context.Context, targetID, targetType, content string) (*jike.Comment, error) {
	return s.client.AddComment(ctx, targetID, targetType, content)
}

// User
func (s *JikeService) GetUserProfile(ctx context.Context, username string) (*jike.User, error) {
	return s.client.GetUserProfile(ctx, username)
}

func (s *JikeService) GetUserPosts(ctx context.Context, username string, loadMoreKey any) (*jike.FeedResponse, error) {
	return s.client.GetUserPosts(ctx, username, loadMoreKey)
}

// Topics
func (s *JikeService) GetTopicFeed(ctx context.Context, topicID string, loadMoreKey any) (*jike.FeedResponse, error) {
	return s.client.GetTopicFeed(ctx, topicID, loadMoreKey)
}

// URL resolution
func (s *JikeService) ResolveURL(rawURL string) (string, error) {
	return jike.ResolveShortURL(rawURL)
}

// Interactions
func (s *JikeService) LikePost(ctx context.Context, postID, targetType string) error {
	return s.client.LikePost(ctx, postID, targetType)
}

func (s *JikeService) UnlikePost(ctx context.Context, postID, targetType string) error {
	return s.client.UnlikePost(ctx, postID, targetType)
}

func (s *JikeService) FollowUser(ctx context.Context, username string) error {
	return s.client.FollowUser(ctx, username)
}

func (s *JikeService) UnfollowUser(ctx context.Context, username string) error {
	return s.client.UnfollowUser(ctx, username)
}
