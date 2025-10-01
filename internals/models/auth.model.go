package models

type User struct {
	ID        string  `db:"id"`
	Email     string  `db:"email"`
	Password  string  `db:"password"`
	Name      *string `db:"name"`
	AvatarUrl *string `db:"avatar_url"`
	Bio       *string `db:"bio"`
}

type AuthRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type UserUpdate struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
}
