package model

import (
	"time"

	parserv1 "github.com/Negat1v9/sum-tel/services/core/internal/grpc/proto"
	"github.com/google/uuid"
)

// model message from kafka topic news-aggregated
// TODO: move to shared module ?
type NewsKafkaType struct {
	Title    string            `json:"title"`
	Summary  string            `json:"summary"`
	Sources  []SourceKafkaType `json:"sources"`
	Language string            `json:"language"`
	Category string            `json:"category"`
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
	Category  string    `db:"category" json:"category"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	NumberOfSources int `db:"number_of_sources" json:"number_of_sources"` // number of sources for this news geneated by DB
}

type NewsSource struct {
	ID        int       `db:"id" json:"id"`
	MessageID int64     `db:"message_id" json:"message_id"` // telegram message id
	NewsID    uuid.UUID `db:"news_id" json:"news_id"`
	ChannelID uuid.UUID `db:"channel_id" json:"channel_id"`

	ChannelName string `db:"channel_name" json:"channel_name"`
}

type NewsList struct {
	TotalRecords int    `json:"total_records"`
	News         []News `json:"news"`
}

func ConvertNewsKafkaTypeToNews(kafkaNews NewsKafkaType) *News {
	return &News{
		ID:       uuid.New(),
		Title:    kafkaNews.Title,
		Summary:  kafkaNews.Summary,
		Language: kafkaNews.Language,
		Category: kafkaNews.Category,
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

func NewNews(id uuid.UUID, title, summary, language, category string) *News {
	return &News{
		ID:        id,
		Title:     title,
		Summary:   summary,
		Language:  language,
		Category:  category,
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

// JSON data about news source for response to user

type NewsSourcesListResponse struct {
	NewsID  string              `json:"news_id"`
	Sources []*parserv1.Message `json:"sources"`
}

func ClearGrpcNewsSources(msgs []*parserv1.Message) []*parserv1.Message {
	for _, msg := range msgs {
		if msg.Date <= 0 {
			msg.Date = 0
		}
		msg.ChannelId = ""
	}
	return msgs
}
