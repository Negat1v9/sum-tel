package rawMessageRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Negat1v9/sum-tel/services/parser/internal/model"
	sqltransaction "github.com/Negat1v9/sum-tel/services/parser/internal/store/sqlTransaction"
	"github.com/jmoiron/sqlx"
)

type rawMessageRepository struct {
	db *sqlx.DB
}

func NewRawMessageRepository(db *sqlx.DB) *rawMessageRepository {
	return &rawMessageRepository{db: db}
}

// create a new message
func (r *rawMessageRepository) CreateMessages(ctx context.Context, tx sqltransaction.Txx, msgs []model.RawMessage) error {
	if len(msgs) == 0 {
		return nil
	}

	var placeholder []string
	args := make([]any, 0, len(msgs)*8)

	for i, msg := range msgs {
		offset := i * 8
		placeholder = append(placeholder, fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %d)", offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8))
		args = append(args, msg.ChannelID, msg.ContentType, msg.TelegramMessageID, msg.HTMLText, msg.Status, msg.MediaURLs, msg.MessageDate, msg.ReceivedAt)
	}

	query := fmt.Sprintf(createMessagesQuery, strings.Join(placeholder, ","))

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

// return a channel messages
func (r *rawMessageRepository) GetChannelMessages(ctx context.Context, chID string, limit, offset int64) ([]model.RawMessage, error) {
	var msgs []model.RawMessage
	err := r.db.SelectContext(ctx, &msgs, getChannelMessagesQuery, chID, limit, offset)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *rawMessageRepository) GetLatestChannelMessage(ctx context.Context, chID string) (model.RawMessage, error) {
	var msg model.RawMessage
	err := r.db.GetContext(ctx, &msg, getLatestChannelMessageQuery, chID)
	if err != nil {
		return model.RawMessage{}, err
	}
	return msg, nil
}

func (r *rawMessageRepository) GetAndProcessedChannelMessages(ctx context.Context, tx sqltransaction.Txx, chID string, limit int64) ([]model.RawMessage, error) {
	var msgs []model.RawMessage
	err := tx.SelectContext(ctx, &msgs, getAndProcessMessagesQuery, chID, limit)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
