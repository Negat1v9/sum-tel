package model

import "time"

type User struct {
	ID         int       `db:"id" json:"id"`
	TelegramID int64     `db:"telegram_id" json:"-"`
	Username   string    `db:"username" json:"username"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	Role       string    `db:"role" json:"-"`
}
