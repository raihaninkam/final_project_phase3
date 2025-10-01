package models

type UserProfile struct {
	Id     string  `json:"id"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar_url"`
	Bio    *string `json:"bio"`
}
