package model

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Username      string    `db:"username" json:"username"`
	Title         string    `db:"title" json:"title"`
	Description   string    `db:"description" json:"description,omitempty"`
	ParseInterval int       `db:"parse_interval" json:"-"`
	LastParsedAt  time.Time `db:"last_parsed_at" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
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
		Description:   description,
		ParseInterval: parseInterval,
		LastParsedAt:  time.Now(),
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}
}
