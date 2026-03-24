package jike

// User represents a Jike user profile.
type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	ScreenName     string `json:"screenName"`
	Bio            string `json:"bio,omitempty"`
	BriefIntro     string `json:"briefIntro,omitempty"`
	AvatarImage    any    `json:"avatarImage,omitempty"`
	Gender         string `json:"gender,omitempty"`
	IsBanned       bool   `json:"isBanned,omitempty"`
	IsBetaUser     bool   `json:"isBetaUser,omitempty"`
	StatsCount     any    `json:"statsCount,omitempty"`
	ProfileImageURL string `json:"profileImageUrl,omitempty"`
}

// Post represents a Jike post (originalPost or repost).
type Post struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Content      string    `json:"content"`
	User         User      `json:"user"`
	CreatedAt    string `json:"createdAt"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
	ShareCount   int       `json:"shareCount"`
	RepostCount  int       `json:"repostCount"`
	Topic        any       `json:"topic,omitempty"`
	Pictures     any       `json:"pictures,omitempty"`
	Liked        bool      `json:"liked"`
	Collected    bool      `json:"collected"`
	LinkInfo     any       `json:"linkInfo,omitempty"`
	TargetType   string    `json:"targetType,omitempty"`
	Target       any       `json:"target,omitempty"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	User       User      `json:"user"`
	CreatedAt  string `json:"createdAt"`
	LikeCount  int       `json:"likeCount"`
	Liked      bool      `json:"liked"`
	Level      int       `json:"level"`
	ReplyTo    any       `json:"replyTo,omitempty"`
	TargetType string    `json:"targetType"`
	TargetID   string    `json:"targetId"`
	HotReplies any       `json:"hotReplies,omitempty"`
}

// FeedResponse is the common response wrapper for feed-like APIs.
type FeedResponse struct {
	Data        []Post `json:"data"`
	LoadMoreKey any    `json:"loadMoreKey"`
}

type CommentResponse struct {
	Data        []Comment `json:"data"`
	LoadMoreKey any       `json:"loadMoreKey"`
}

// SearchResult for integrated search.
type SearchResult struct {
	Data        []Post `json:"data"`
	LoadMoreKey any    `json:"loadMoreKey"`
}

// SessionCreateResponse from sessions.create.
type SessionCreateResponse struct {
	UUID string `json:"uuid"`
}

// LoginConfirmResponse from sessions.wait_for_confirmation.
type LoginConfirmResponse struct {
	User User `json:"user"`
}
