package main

import (
	"encoding/json"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newMCPServer(app *AppServer) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "jike-mcp",
		Version: "1.0.0",
	}, nil)

	inputSchema := func(props map[string]any, required []string) json.RawMessage {
		schema := map[string]any{
			"type":       "object",
			"properties": props,
		}
		if len(required) > 0 {
			schema["required"] = required
		}
		raw, _ := json.Marshal(schema)
		return raw
	}

	// Auth tools
	s.AddTool(&mcp.Tool{
		Name:        "check_login_status",
		Description: "Check if currently logged in to Jike",
		InputSchema: inputSchema(map[string]any{}, nil),
	}, app.handleCheckLoginStatus)

	s.AddTool(&mcp.Tool{
		Name:        "get_login_qrcode",
		Description: "Get QR code for Jike login. Returns a base64 PNG image that the user should scan with the Jike app.",
		InputSchema: inputSchema(map[string]any{}, nil),
	}, app.handleGetLoginQRCode)

	s.AddTool(&mcp.Tool{
		Name:        "wait_for_login",
		Description: "Wait for the user to scan the QR code and confirm login. Call this after get_login_qrcode.",
		InputSchema: inputSchema(map[string]any{
			"uuid": map[string]any{"type": "string", "description": "The session UUID returned by get_login_qrcode"},
		}, []string{"uuid"}),
	}, app.handleWaitForLogin)

	s.AddTool(&mcp.Tool{
		Name:        "logout",
		Description: "Logout from Jike by deleting stored tokens",
		InputSchema: inputSchema(map[string]any{}, nil),
	}, app.handleLogout)

	// Feed tools
	s.AddTool(&mcp.Tool{
		Name:        "get_following_feeds",
		Description: "Get posts from users you follow on Jike",
		InputSchema: inputSchema(map[string]any{
			"load_more_key": map[string]any{"type": "string", "description": "Pagination key for loading more results. Omit for first page."},
		}, nil),
	}, app.handleGetFollowingFeeds)

	s.AddTool(&mcp.Tool{
		Name:        "get_recommend_feeds",
		Description: "Get recommended posts on Jike",
		InputSchema: inputSchema(map[string]any{
			"load_more_key": map[string]any{"type": "string", "description": "Pagination key for loading more results. Omit for first page."},
		}, nil),
	}, app.handleGetRecommendFeeds)

	// Search tool
	s.AddTool(&mcp.Tool{
		Name:        "search",
		Description: "Search for posts, users, or topics on Jike",
		InputSchema: inputSchema(map[string]any{
			"keyword":       map[string]any{"type": "string", "description": "Search keyword"},
			"load_more_key": map[string]any{"type": "string", "description": "Pagination key for loading more results"},
		}, []string{"keyword"}),
	}, app.handleSearch)

	// Post tools
	s.AddTool(&mcp.Tool{
		Name:        "get_post_detail",
		Description: "Get detailed information about a specific Jike post",
		InputSchema: inputSchema(map[string]any{
			"post_id":   map[string]any{"type": "string", "description": "The post ID"},
			"post_type": map[string]any{"type": "string", "description": "Post type: ORIGINAL_POST or REPOST. Defaults to ORIGINAL_POST.", "enum": []string{"ORIGINAL_POST", "REPOST"}},
		}, []string{"post_id"}),
	}, app.handleGetPostDetail)

	s.AddTool(&mcp.Tool{
		Name:        "create_post",
		Description: "Create a new post on Jike",
		InputSchema: inputSchema(map[string]any{
			"content":  map[string]any{"type": "string", "description": "Post content text"},
			"topic_id": map[string]any{"type": "string", "description": "Optional topic ID to post under"},
		}, []string{"content"}),
	}, app.handleCreatePost)

	// Topic tools
	s.AddTool(&mcp.Tool{
		Name:        "get_topic_feed",
		Description: "Get posts from a Jike topic (圈子). Use search to find topic IDs first.",
		InputSchema: inputSchema(map[string]any{
			"topic_id":      map[string]any{"type": "string", "description": "Topic ID (from search results)"},
			"load_more_key": map[string]any{"type": "string", "description": "Pagination key for loading more results"},
		}, []string{"topic_id"}),
	}, app.handleGetTopicFeed)

	// Comment tools
	s.AddTool(&mcp.Tool{
		Name:        "get_comments",
		Description: "Get comments for a Jike post",
		InputSchema: inputSchema(map[string]any{
			"target_id":     map[string]any{"type": "string", "description": "The post ID to get comments for"},
			"target_type":   map[string]any{"type": "string", "description": "Target type: ORIGINAL_POST or REPOST. Defaults to ORIGINAL_POST."},
			"load_more_key": map[string]any{"type": "string", "description": "Pagination key for loading more comments"},
		}, []string{"target_id"}),
	}, app.handleGetComments)

	s.AddTool(&mcp.Tool{
		Name:        "add_comment",
		Description: "Add a comment to a Jike post",
		InputSchema: inputSchema(map[string]any{
			"target_id":   map[string]any{"type": "string", "description": "The post ID to comment on"},
			"target_type": map[string]any{"type": "string", "description": "Target type: ORIGINAL_POST or REPOST. Defaults to ORIGINAL_POST."},
			"content":     map[string]any{"type": "string", "description": "Comment content"},
		}, []string{"target_id", "content"}),
	}, app.handleAddComment)

	// User tools
	s.AddTool(&mcp.Tool{
		Name:        "get_user_profile",
		Description: "Get a Jike user's profile by username or Jike URL (e.g., https://okjk.co/xxx)",
		InputSchema: inputSchema(map[string]any{
			"username": map[string]any{"type": "string", "description": "Jike username or profile URL (e.g., https://okjk.co/xxx)"},
		}, []string{"username"}),
	}, app.handleGetUserProfile)

	s.AddTool(&mcp.Tool{
		Name:        "get_user_posts",
		Description: "Get posts by a specific Jike user. Accepts username or Jike URL (e.g., https://okjk.co/xxx)",
		InputSchema: inputSchema(map[string]any{
			"username":      map[string]any{"type": "string", "description": "Jike username or profile URL (e.g., https://okjk.co/xxx)"},
			"load_more_key": map[string]any{"type": "string", "description": "Pagination key for loading more results"},
		}, []string{"username"}),
	}, app.handleGetUserPosts)

	// Interaction tools
	s.AddTool(&mcp.Tool{
		Name:        "like_post",
		Description: "Like a Jike post",
		InputSchema: inputSchema(map[string]any{
			"post_id":     map[string]any{"type": "string", "description": "The post ID to like"},
			"target_type": map[string]any{"type": "string", "description": "Target type: ORIGINAL_POST or REPOST. Defaults to ORIGINAL_POST."},
		}, []string{"post_id"}),
	}, app.handleLikePost)

	s.AddTool(&mcp.Tool{
		Name:        "unlike_post",
		Description: "Unlike a Jike post",
		InputSchema: inputSchema(map[string]any{
			"post_id":     map[string]any{"type": "string", "description": "The post ID to unlike"},
			"target_type": map[string]any{"type": "string", "description": "Target type: ORIGINAL_POST or REPOST. Defaults to ORIGINAL_POST."},
		}, []string{"post_id"}),
	}, app.handleUnlikePost)

	s.AddTool(&mcp.Tool{
		Name:        "follow_user",
		Description: "Follow a Jike user",
		InputSchema: inputSchema(map[string]any{
			"username": map[string]any{"type": "string", "description": "Username of the user to follow"},
		}, []string{"username"}),
	}, app.handleFollowUser)

	s.AddTool(&mcp.Tool{
		Name:        "unfollow_user",
		Description: "Unfollow a Jike user",
		InputSchema: inputSchema(map[string]any{
			"username": map[string]any{"type": "string", "description": "Username of the user to unfollow"},
		}, []string{"username"}),
	}, app.handleUnfollowUser)

	return s
}

func newStreamableHandler(s *mcp.Server) http.Handler {
	return mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return s
	}, nil)
}
