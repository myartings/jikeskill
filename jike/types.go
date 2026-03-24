package jike

import "time"

// User represents a Jike user profile.
type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	ScreenName     string `json:"screenName"`
	Bio            string `json:"bio"`
	AvatarImage    string `json:"avatarImage"`
	Gender         string `json:"gender"`
	IsBanned       bool   `json:"isBanned"`
	IsBetaUser     bool   `json:"isBetaUser"`
	StatsCount     Stats  `json:"statsCount"`
	ProfileImageID string `json:"profileImageUrl"`
}

type Stats struct {
	TopicSubscribed int `json:"topicSubscribed"`
	TopicCreated    int `json:"topicCreated"`
	Following       int `json:"followingCount"`
	Follower        int `json:"followerCount"`
	Liked           int `json:"liked"`
	Highlighted     int `json:"highlightedPersonalUpdates"`
}

// Post represents a Jike post (originalPost or repost).
type Post struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Content           string    `json:"content"`
	User              User      `json:"user"`
	CreatedAt         time.Time `json:"createdAt"`
	LikeCount         int       `json:"likeCount"`
	CommentCount      int       `json:"commentCount"`
	ShareCount        int       `json:"shareCount"`
	Topic             *Topic    `json:"topic"`
	Pictures          []Picture `json:"pictures"`
	Liked             bool      `json:"liked"`
	Collected         bool      `json:"collected"`
	ReadTrackInfo     any       `json:"readTrackInfo"`
	LinkInfo          any       `json:"linkInfo"`
	RepostCount       int       `json:"repostCount"`
	TargetType        string    `json:"targetType,omitempty"`
	Target            *Post     `json:"target,omitempty"`
	ScrollingSubtitle string    `json:"scrollingSubtitle,omitempty"`
}

type Picture struct {
	PicURL    string `json:"picUrl"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CropperPosX float64 `json:"cropperPosX"`
	CropperPosY float64 `json:"cropperPosY"`
}

type Topic struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	SquarePicURL string `json:"squarePicUrl"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Content       string    `json:"content"`
	User          User      `json:"user"`
	CreatedAt     time.Time `json:"createdAt"`
	LikeCount     int       `json:"likeCount"`
	Liked         bool      `json:"liked"`
	Level         int       `json:"level"`
	ReplyTo       *Comment  `json:"replyTo,omitempty"`
	TargetType    string    `json:"targetType"`
	TargetID      string    `json:"targetId"`
	ThreadID      string    `json:"threadId,omitempty"`
	HotReplies    []Comment `json:"hotReplies,omitempty"`
}

// FeedResponse is the common response wrapper for feed-like APIs.
type FeedResponse struct {
	Data       []Post `json:"data"`
	LoadMoreKey any    `json:"loadMoreKey"`
}

type CommentResponse struct {
	Data       []Comment `json:"data"`
	LoadMoreKey any      `json:"loadMoreKey"`
}

// SearchResult for integrated search.
type SearchResult struct {
	Data       []Post `json:"data"`
	LoadMoreKey any   `json:"loadMoreKey"`
}

// SessionCreateResponse from sessions.create.
type SessionCreateResponse struct {
	UUID string `json:"uuid"`
}

// LoginConfirmResponse from sessions.wait_for_confirmation.
type LoginConfirmResponse struct {
	User User `json:"user"`
}
