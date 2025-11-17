package model

import (
	"time"

	"github.com/lib/pq"
)

type RawMessage struct {
	ID                int64          `json:"id" db:"id"`
	ChannelID         string         `json:"channel_id" db:"channel_id"`                   // channel ID from which the message was received
	ContentType       string         `json:"content_type" db:"content_type"`               // "text", "image", "text_image"
	TelegramMessageID int64          `json:"telegram_message_id" db:"telegram_message_id"` // Telegram message ID
	HTMLText          string         `json:"html_text" db:"html_text"`                     // text content in HTML format
	Status            string         `json:"status" db:"status"`                           // "new", "processed"
	MediaURLs         pq.StringArray `json:"media_urls" db:"media_urls"`                   // URLs of media (images) in the message
	MessageDate       time.Time      `json:"message_date" db:"message_date"`               // date of the message in the channel
	ReceivedAt        time.Time      `json:"received_at" db:"received_at"`                 // date when the message was received by the system
}

func NewRawMsg(chID string, contentType string, telegramMessageID int64, htmlText string, mediaURLs pq.StringArray, messageDate time.Time) RawMessage {
	return RawMessage{
		ChannelID:         chID,
		ContentType:       contentType,
		TelegramMessageID: telegramMessageID,
		HTMLText:          htmlText,
		Status:            "new",
		MediaURLs:         mediaURLs,
		MessageDate:       messageDate,
		ReceivedAt:        time.Now().UTC(),
	}
}
