package domain

import tgparser "github.com/Negat1v9/sum-tel/services/parser/internal/parser/tgParser"

type KafkaRawMessage struct {
	ChannelID string `json:"channel_id"` // source channel id from parser service
	MessageID int64  `json:"message_id"` // message id from parser service
	Text      string `json:"text"`       // text content
}

func ConvetParsedMessagesToAny(chID string, msgs []tgparser.ParsedMessage) []any {
	r := make([]any, 0, len(msgs))
	for _, m := range msgs {
		r = append(r, KafkaRawMessage{
			ChannelID: chID,
			MessageID: m.MsgId,
			Text:      m.Text,
		})
	}
	return r
}

// convert RawMessage slice to any slice with KafkaRawMessage
func ConvertRawMessagesToAny(msgs []RawMessage) []any {
	r := make([]any, 0, len(msgs))
	for _, m := range msgs {
		r = append(r, KafkaRawMessage{
			ChannelID: m.ChannelID,
			MessageID: m.TelegramMessageID,
			Text:      tgparser.CleanMessageText(m.HTMLText),
		})
	}
	return r
}
