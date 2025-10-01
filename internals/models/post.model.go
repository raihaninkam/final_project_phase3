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
