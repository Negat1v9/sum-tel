package model

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID         int       `db:"id" json:"id"`
	TelegramID int64     `db:"telegram_id" json:"-"`
	Username   string    `db:"username" json:"username"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	Role       string    `db:"role" json:"-"`
}

func NewUser(tgID int64, username string, role string) *User {
	return &User{
		TelegramID: tgID,
		Username:   username,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		IsActive:   true,
		Role:       role,
	}
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
