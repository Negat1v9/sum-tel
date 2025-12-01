package model

import (
	"time"

	"github.com/google/uuid"
)

// model message from kafka topic news-aggregated
// TODO: move to shared module ?
type NewsKafkaType struct {
	Title    string            `json:"title"`
	Summary  string            `json:"summary"`
	Sources  []SourceKafkaType `json:"sources"`
	Language string            `json:"language"`
}

type SourceKafkaType struct {
	ChannelID string `json:"channel_id"`
	MessageID int64  `json:"message_id"`
}

type News struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Summary   string    `db:"summary" json:"summary"`
	Language  string    `db:"language" json:"language"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type NewsSource struct {
	ID        int       `db:"id" json:"id"`
	MessageID int64     `db:"message_id" json:"message_id"`
	NewsID    uuid.UUID `db:"news_id" json:"news_id"`
	ChannelID uuid.UUID `db:"channel_id" json:"channel_id"`
}

func ConvertNewsKafkaTypeToNews(kafkaNews NewsKafkaType) *News {
	return &News{
		ID:       uuid.New(),
		Title:    kafkaNews.Title,
		Summary:  kafkaNews.Summary,
		Language: kafkaNews.Language,
	}
}

func ConvertSourceKafkaTypesToNewsSources(newsID uuid.UUID, kafkaSource []SourceKafkaType) []NewsSource {
	r := make([]NewsSource, 0, len(kafkaSource))
	for _, m := range kafkaSource {
		channelID, _ := uuid.Parse(m.ChannelID)
		r = append(r, NewsSource{
			MessageID: m.MessageID,
			NewsID:    newsID,
			ChannelID: channelID,
		})
	}
	return r
}

func NewNews(id uuid.UUID, title, summary, language string) *News {
	return &News{
		ID:        id,
		Title:     title,
		Summary:   summary,
		Language:  language,
		CreatedAt: time.Now(),
	}
}

func NewNewsSource(id int, newsID, channelID uuid.UUID) *NewsSource {
	return &NewsSource{
		ID:        id,
		NewsID:    newsID,
		ChannelID: channelID,
	}
}
