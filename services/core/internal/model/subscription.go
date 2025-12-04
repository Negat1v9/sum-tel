package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSubscription struct {
	ID           int       `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	ChannelID    uuid.UUID `db:"channel_id" json:"channel_id"`
	SubscribedAt time.Time `db:"subscribed_at" json:"subscribed_at"`
}

type UserSubscriptionWithChannel struct {
	UserSubscription
	Channel Channel `json:"channel"`
}

func NewSub(userID int64, channelID uuid.UUID) *UserSubscription {
	return &UserSubscription{
		UserID:       userID,
		ChannelID:    channelID,
		SubscribedAt: time.Now(),
	}
}
