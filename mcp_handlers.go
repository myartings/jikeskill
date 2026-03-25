package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/myartings/jikeskill/jike"
)

func parseArgs(req *mcp.CallToolRequest) map[string]any {
	var args map[string]any
	if req.Params.Arguments != nil {
		json.Unmarshal(req.Params.Arguments, &args)
	}
	if args == nil {
		args = map[string]any{}
	}
	return args
}

func getStringArg(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func toJSON(v any) string {
	raw, _ := json.MarshalIndent(v, "", "  ")
	return string(raw)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %s", err.Error())}},
		IsError: true,
	}
}

// resolveUsername resolves a username from either a plain username or a Jike URL (e.g., https://okjk.co/xxx).
func resolveUsername(input string) (string, error) {
	if strings.Contains(input, "://") || strings.Contains(input, "okjk.co") || strings.Contains(input, "okjike.com") {
		return jike.ResolveShortURL(input)
	}
	return input, nil
}

// Auth handlers

func (a *AppServer) handleCheckLoginStatus(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	loggedIn, user, err := a.service.CheckLoginStatus(ctx)
	if err != nil {
		return errorResult(err), nil
	}
	if !loggedIn {
		return textResult("Not logged in. Use get_login_qrcode to start login."), nil
	}
	return textResult(fmt.Sprintf("Logged in as: %s (%s)", user.ScreenName, user.Username)), nil
}

func (a *AppServer) handleGetLoginQRCode(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uuid, qrBase64, err := a.service.CreateLoginSession(ctx)
	if err != nil {
		return errorResult(err), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("QR code generated. Session UUID: %s\nPlease scan the QR code with the Jike app to login.\nThen call wait_for_login with this UUID.", uuid),
			},
			&mcp.ImageContent{
				Data:     []byte(qrBase64),
				MIMEType: "image/png",
			},
		},
	}, nil
}

func (a *AppServer) handleWaitForLogin(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	uuid := getStringArg(args, "uuid")
	if uuid == "" {
		return errorResult(fmt.Errorf("uuid is required")), nil
	}

	user, err := a.service.WaitForLogin(ctx, uuid)
	if err != nil {
		return errorResult(err), nil
	}

	if user != nil {
		return textResult(fmt.Sprintf("Login successful! Welcome, %s (%s)", user.ScreenName, user.Username)), nil
	}
	return textResult("Login successful!"), nil
}

func (a *AppServer) handleLogout(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if err := a.service.Logout(); err != nil {
		return errorResult(err), nil
	}
	return textResult("Logged out successfully."), nil
}

// Feed handlers

func (a *AppServer) handleGetFollowingFeeds(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	loadMoreKey := getStringArg(args, "load_more_key")
	var lmk any
	if loadMoreKey != "" {
		json.Unmarshal([]byte(loadMoreKey), &lmk)
	}

	resp, err := a.service.GetFollowingFeeds(ctx, lmk)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(resp)), nil
}

func (a *AppServer) handleGetRecommendFeeds(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	loadMoreKey := getStringArg(args, "load_more_key")
	var lmk any
	if loadMoreKey != "" {
		json.Unmarshal([]byte(loadMoreKey), &lmk)
	}

	resp, err := a.service.GetRecommendFeeds(ctx, lmk)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(resp)), nil
}

// Search handler

func (a *AppServer) handleSearch(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	keyword := getStringArg(args, "keyword")
	if keyword == "" {
		return errorResult(fmt.Errorf("keyword is required")), nil
	}

	loadMoreKey := getStringArg(args, "load_more_key")
	var lmk any
	if loadMoreKey != "" {
		json.Unmarshal([]byte(loadMoreKey), &lmk)
	}

	resp, err := a.service.Search(ctx, keyword, lmk)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(resp)), nil
}

// Post handlers

func (a *AppServer) handleGetPostDetail(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	postID := getStringArg(args, "post_id")
	if postID == "" {
		return errorResult(fmt.Errorf("post_id is required")), nil
	}
	postType := getStringArg(args, "post_type")

	post, err := a.service.GetPostDetail(ctx, postID, postType)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(post)), nil
}

func (a *AppServer) handleCreatePost(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	content := getStringArg(args, "content")
	if content == "" {
		return errorResult(fmt.Errorf("content is required")), nil
	}
	topicID := getStringArg(args, "topic_id")

	post, err := a.service.CreatePost(ctx, content, topicID, nil)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(fmt.Sprintf("Post created successfully!\n%s", toJSON(post))), nil
}

// Comment handlers

func (a *AppServer) handleGetComments(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	targetID := getStringArg(args, "target_id")
	if targetID == "" {
		return errorResult(fmt.Errorf("target_id is required")), nil
	}
	targetType := getStringArg(args, "target_type")
	if targetType == "" {
		targetType = "ORIGINAL_POST"
	}

	loadMoreKey := getStringArg(args, "load_more_key")
	var lmk any
	if loadMoreKey != "" {
		json.Unmarshal([]byte(loadMoreKey), &lmk)
	}

	resp, err := a.service.GetComments(ctx, targetID, targetType, lmk)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(resp)), nil
}

func (a *AppServer) handleAddComment(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	targetID := getStringArg(args, "target_id")
	if targetID == "" {
		return errorResult(fmt.Errorf("target_id is required")), nil
	}
	content := getStringArg(args, "content")
	if content == "" {
		return errorResult(fmt.Errorf("content is required")), nil
	}
	targetType := getStringArg(args, "target_type")
	if targetType == "" {
		targetType = "ORIGINAL_POST"
	}

	comment, err := a.service.AddComment(ctx, targetID, targetType, content)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(fmt.Sprintf("Comment added successfully!\n%s", toJSON(comment))), nil
}

// User handlers

func (a *AppServer) handleGetUserProfile(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	username := getStringArg(args, "username")
	if username == "" {
		return errorResult(fmt.Errorf("username is required")), nil
	}

	resolved, err := resolveUsername(username)
	if err != nil {
		return errorResult(fmt.Errorf("resolve URL: %w", err)), nil
	}

	user, err := a.service.GetUserProfile(ctx, resolved)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(user)), nil
}

func (a *AppServer) handleGetUserPosts(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	username := getStringArg(args, "username")
	if username == "" {
		return errorResult(fmt.Errorf("username is required")), nil
	}

	resolved, err := resolveUsername(username)
	if err != nil {
		return errorResult(fmt.Errorf("resolve URL: %w", err)), nil
	}

	loadMoreKey := getStringArg(args, "load_more_key")
	var lmk any
	if loadMoreKey != "" {
		json.Unmarshal([]byte(loadMoreKey), &lmk)
	}

	resp, err := a.service.GetUserPosts(ctx, resolved, lmk)
	if err != nil {
		return errorResult(err), nil
	}
	return textResult(toJSON(resp)), nil
}

// Interaction handlers

func (a *AppServer) handleLikePost(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	postID := getStringArg(args, "post_id")
	if postID == "" {
		return errorResult(fmt.Errorf("post_id is required")), nil
	}
	targetType := getStringArg(args, "target_type")

	if err := a.service.LikePost(ctx, postID, targetType); err != nil {
		return errorResult(err), nil
	}
	return textResult("Post liked successfully!"), nil
}

func (a *AppServer) handleUnlikePost(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	postID := getStringArg(args, "post_id")
	if postID == "" {
		return errorResult(fmt.Errorf("post_id is required")), nil
	}
	targetType := getStringArg(args, "target_type")

	if err := a.service.UnlikePost(ctx, postID, targetType); err != nil {
		return errorResult(err), nil
	}
	return textResult("Post unliked successfully!"), nil
}

func (a *AppServer) handleFollowUser(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	username := getStringArg(args, "username")
	if username == "" {
		return errorResult(fmt.Errorf("username is required")), nil
	}

	if err := a.service.FollowUser(ctx, username); err != nil {
		return errorResult(err), nil
	}
	return textResult(fmt.Sprintf("Successfully followed user: %s", username)), nil
}

func (a *AppServer) handleUnfollowUser(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := parseArgs(req)
	username := getStringArg(args, "username")
	if username == "" {
		return errorResult(fmt.Errorf("username is required")), nil
	}

	if err := a.service.UnfollowUser(ctx, username); err != nil {
		return errorResult(err), nil
	}
	return textResult(fmt.Sprintf("Successfully unfollowed user: %s", username)), nil
}
