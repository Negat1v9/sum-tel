package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSubscription struct {
	ID           int       `db:"id"`
	UserID       int64     `db:"user_id"`
	ChannelID    uuid.UUID `db:"channel_id"`
	SubscribedAt time.Time `db:"subscribed_at"`
}

func NewSub(userID int64, channelID uuid.UUID) *UserSubscription {
	return &UserSubscription{
		UserID:       userID,
		ChannelID:    channelID,
		SubscribedAt: time.Now(),
	}
}
