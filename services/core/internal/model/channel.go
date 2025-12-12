package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	// channel is active - parse it messages and create news on channel messages
	ChannelStatusActive string = "active"
	// channel is inactive- not parsing it
	ChannelStatusInActive string = "inactive"
	// channel is banned - not show it, no parsing it and not subsciribe on it
	ChannelStatusBanned string = "banned"
)

type Channel struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Status        string    `db:"status" json:"status"`
	CreatedBy     int       `db:"created_by" json:"-"`
	Username      string    `db:"username" json:"username"`
	Title         string    `db:"title" json:"title"`
	Description   string    `db:"description" json:"description,omitempty"`
	ParseInterval int       `db:"parse_interval" json:"-"`
	LastParsedAt  time.Time `db:"last_parsed_at" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func NewChannel(id uuid.UUID, status, username, title, description string, createdBy, parseInterval int, createdAt time.Time) *Channel {
	return &Channel{
		ID:            id,
		Status:        status,
		CreatedBy:     createdBy,
		Username:      username,
		Title:         title,
		Description:   description,
		ParseInterval: parseInterval,
		LastParsedAt:  time.Now(),
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}
}
