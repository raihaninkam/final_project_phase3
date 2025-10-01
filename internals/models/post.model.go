package models

import (
	"time"
)

type Posts struct {
	Id        string     `db:"id"`
	UserId    string     `db:"user_id"`
	Content   string     `db:"content_text"`
	ImageUrl  string     `db:"image_url"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type PostsRequest struct {
	Content  *string `json:"content_text" form:"content_text"`
	ImageUrl *string `json:"image_url" form:"image_url"`
}

type PostWithUser struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	Content    string    `json:"content_text"`
	ImageUrl   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UserName   *string   `json:"user_name"`
	UserAvatar *string   `json:"user_avatar"`
}

type Posting struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	Content       string    `json:"content" db:"content"`
	ImageUrl      string    `json:"image_url,omitempty" db:"image_url"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	UserName      string    `json:"user_name,omitempty" db:"user_name"`
	UserAvatarUrl *string   `json:"user_avatar_url,omitempty" db:"user_avatar_url"`
	LikeCount     int       `json:"like_count" db:"like_count"`
	CommentCount  int       `json:"comment_count" db:"comment_count"`
	FollowerCount int       `json:"follower_count" db:"follower_count"`
	IsLiked       bool      `json:"is_liked" db:"is_liked"`
}
