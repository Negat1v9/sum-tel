package domain

type RawMessage struct {
	ChannelID string `json:"channel_id"` // source channel id from parser service
	MessageID int64  `json:"message_id"` // message id from telegram channel
	Text      string `json:"text"`       // text content
}
