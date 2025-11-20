package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID            uuid.UUID      `db:"id"`
	Username      string         `db:"username"`
	Title         string         `db:"title"`
	Description   sql.NullString `db:"description"`
	ParseInterval int            `db:"parse_interval"`
	LastParsedAt  sql.NullTime   `db:"last_parsed_at"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

type ChannelMessages struct {
	Type      string   `json:"type"`
	Text      string   `json:"text"`
	HTMLText  string   `json:"html_text"`
	Link      string   `json:"link"`
	MsgID     int64    `json:"msg_id"`
	Date      int64    `json:"date"`
	PhotoURLs []string `json:"photo_urls,omitempty"`
}

func NewChannel(id uuid.UUID, username, title, description string, parseInterval int, createdAt time.Time) *Channel {
	return &Channel{
		ID:            id,
		Username:      username,
		Title:         title,
		Description:   sql.NullString{String: description, Valid: description != ""},
		ParseInterval: parseInterval,
		LastParsedAt:  sql.NullTime{Valid: false},
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}
}
