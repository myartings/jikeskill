package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myartings/jikeskill/jike"
)

func resolveUsernameParam(input string) (string, error) {
	if strings.Contains(input, "://") || strings.Contains(input, "okjk.co") || strings.Contains(input, "okjike.com") {
		return jike.ResolveShortURL(input)
	}
	return input, nil
}

func (a *AppServer) setupRoutes(r *gin.Engine) {
	// MCP endpoint (Streamable HTTP)
	mcpHandler := newStreamableHandler(a.mcpServer)
	r.Any("/mcp", gin.WrapH(mcpHandler))

	// REST API endpoints
	api := r.Group("/api/v1")
	{
		api.GET("/status", a.apiCheckStatus)
		api.POST("/login/qrcode", a.apiGetQRCode)
		api.POST("/login/wait", a.apiWaitForLogin)
		api.POST("/logout", a.apiLogout)
		api.POST("/feeds/following", a.apiGetFollowingFeeds)
		api.POST("/feeds/recommend", a.apiGetRecommendFeeds)
		api.POST("/search", a.apiSearch)
		api.POST("/post/detail", a.apiGetPostDetail)
		api.POST("/post/create", a.apiCreatePost)
		api.POST("/comments/list", a.apiGetComments)
		api.POST("/comments/add", a.apiAddComment)
		api.GET("/user/:username", a.apiGetUserProfile)
		api.POST("/user/:username/posts", a.apiGetUserPosts)
		api.POST("/like", a.apiLikePost)
		api.POST("/unlike", a.apiUnlikePost)
		api.POST("/follow", a.apiFollowUser)
		api.POST("/unfollow", a.apiUnfollowUser)
	}
}

// REST API handlers

func (a *AppServer) apiCheckStatus(c *gin.Context) {
	loggedIn, user, _ := a.service.CheckLoginStatus(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"logged_in": loggedIn, "user": user})
}

func (a *AppServer) apiGetQRCode(c *gin.Context) {
	uuid, qr, err := a.service.CreateLoginSession(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uuid": uuid, "qrcode_base64": qr})
}

func (a *AppServer) apiWaitForLogin(c *gin.Context) {
	var req struct {
		UUID string `json:"uuid"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}
	user, err := a.service.WaitForLogin(c.Request.Context(), req.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (a *AppServer) apiLogout(c *gin.Context) {
	if err := a.service.Logout(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (a *AppServer) apiGetFollowingFeeds(c *gin.Context) {
	var req struct {
		LoadMoreKey any `json:"loadMoreKey"`
	}
	c.ShouldBindJSON(&req)
	resp, err := a.service.GetFollowingFeeds(c.Request.Context(), req.LoadMoreKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *AppServer) apiGetRecommendFeeds(c *gin.Context) {
	var req struct {
		LoadMoreKey any `json:"loadMoreKey"`
	}
	c.ShouldBindJSON(&req)
	resp, err := a.service.GetRecommendFeeds(c.Request.Context(), req.LoadMoreKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *AppServer) apiSearch(c *gin.Context) {
	var req struct {
		Keyword     string `json:"keyword"`
		LoadMoreKey any    `json:"loadMoreKey"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}
	resp, err := a.service.Search(c.Request.Context(), req.Keyword, req.LoadMoreKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *AppServer) apiGetPostDetail(c *gin.Context) {
	var req struct {
		PostID   string `json:"post_id"`
		PostType string `json:"post_type"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	post, err := a.service.GetPostDetail(c.Request.Context(), req.PostID, req.PostType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (a *AppServer) apiCreatePost(c *gin.Context) {
	var req struct {
		Content     string   `json:"content"`
		TopicID     string   `json:"topic_id"`
		PictureKeys []string `json:"picture_keys"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}
	post, err := a.service.CreatePost(c.Request.Context(), req.Content, req.TopicID, req.PictureKeys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (a *AppServer) apiGetComments(c *gin.Context) {
	// Accept multiple parameter name formats (target_id, targetId, post_id, id)
	var raw map[string]any
	c.ShouldBindJSON(&raw)
	targetID := ""
	for _, key := range []string{"target_id", "targetId", "post_id", "id"} {
		if v, ok := raw[key].(string); ok && v != "" {
			targetID = v
			break
		}
	}
	if targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_id is required"})
		return
	}
	targetType := "ORIGINAL_POST"
	for _, key := range []string{"target_type", "targetType"} {
		if v, ok := raw[key].(string); ok && v != "" {
			targetType = v
			break
		}
	}
	var loadMoreKey any
	if v, ok := raw["loadMoreKey"]; ok {
		loadMoreKey = v
	}
	resp, err := a.service.GetComments(c.Request.Context(), targetID, targetType, loadMoreKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *AppServer) apiAddComment(c *gin.Context) {
	var req struct {
		TargetID   string `json:"target_id"`
		TargetType string `json:"target_type"`
		Content    string `json:"content"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_id and content are required"})
		return
	}
	if req.TargetType == "" {
		req.TargetType = "ORIGINAL_POST"
	}
	comment, err := a.service.AddComment(c.Request.Context(), req.TargetID, req.TargetType, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (a *AppServer) apiGetUserProfile(c *gin.Context) {
	username := c.Param("username")
	resolved, err := resolveUsernameParam(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := a.service.GetUserProfile(c.Request.Context(), resolved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (a *AppServer) apiGetUserPosts(c *gin.Context) {
	username := c.Param("username")
	resolved, err := resolveUsernameParam(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		LoadMoreKey any `json:"loadMoreKey"`
	}
	c.ShouldBindJSON(&req)
	resp, err := a.service.GetUserPosts(c.Request.Context(), resolved, req.LoadMoreKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (a *AppServer) apiLikePost(c *gin.Context) {
	var req struct {
		PostID     string `json:"post_id"`
		TargetType string `json:"target_type"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	if err := a.service.LikePost(c.Request.Context(), req.PostID, req.TargetType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}

func (a *AppServer) apiUnlikePost(c *gin.Context) {
	var req struct {
		PostID     string `json:"post_id"`
		TargetType string `json:"target_type"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}
	if err := a.service.UnlikePost(c.Request.Context(), req.PostID, req.TargetType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

func (a *AppServer) apiFollowUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	if err := a.service.FollowUser(c.Request.Context(), req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

func (a *AppServer) apiUnfollowUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	if err := a.service.UnfollowUser(c.Request.Context(), req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}
