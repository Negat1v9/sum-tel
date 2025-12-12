package domain

type AggregatedNews struct {
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Sources  []Source `json:"sources"`
	Language string   `json:"language"`
	Category string   `json:"category"`
}

type Source struct {
	ChannelID string `json:"channel_id"`
	MessageID int64  `json:"message_id"`
}

type UnmatchedMessage struct {
	ChannelID string `json:"channel_id"`
	MessageID int64  `json:"message_id"`
	Reason    string `json:"reason"` // spam, adv, illegal,
}

type AggregationResponse struct {
	AggregatedNews    []AggregatedNews   `json:"aggregated_news"`
	UnmatchedMessages []UnmatchedMessage `json:"unmatched_messages"`
}

func ConvertAggregateNewsdResponseToAny(news []AggregatedNews) []any {
	result := make([]any, 0, len(news))
	for _, n := range news {
		result = append(result, n)
	}
	return result
}
