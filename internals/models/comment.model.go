package models

import "time"

type Comment struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"user_id" db:"user_id"`
	PostId    string    `json:"post_id" db:"post_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CommentWithUser struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	PostId     string    `json:"post_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UserName   *string   `json:"user_name"`
	UserAvatar *string   `json:"user_avatar"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}
