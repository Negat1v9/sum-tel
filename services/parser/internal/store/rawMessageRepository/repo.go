package rawMessageRepository

import (
	"context"
	"fmt"
	"strings"

	parserv1 "github.com/Negat1v9/sum-tel/services/parser/internal/api/proto"
	"github.com/Negat1v9/sum-tel/services/parser/internal/domain"
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
func (r *rawMessageRepository) CreateMessages(ctx context.Context, tx sqltransaction.Txx, msgs []domain.RawMessage) error {
	if len(msgs) == 0 {
		return nil
	}

	var placeholder []string
	args := make([]any, 0, len(msgs)*8)

	for i, msg := range msgs {
		offset := i * 8
		placeholder = append(placeholder, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8))
		args = append(args, msg.ChannelID, msg.ContentType, msg.TelegramMessageID, msg.HTMLText, msg.Status, msg.MediaURLs, msg.MessageDate, msg.ReceivedAt)
	}

	query := fmt.Sprintf(createMessagesQuery, strings.Join(placeholder, ","))

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

// return a channel messages
func (r *rawMessageRepository) GetChannelMessages(ctx context.Context, chID string, limit, offset int64) ([]domain.RawMessage, error) {
	var msgs []domain.RawMessage
	err := r.db.SelectContext(ctx, &msgs, getChannelMessagesQuery, chID, limit, offset)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *rawMessageRepository) GetLatestChannelMessage(ctx context.Context, chID string) (domain.RawMessage, error) {
	var msg domain.RawMessage
	err := r.db.GetContext(ctx, &msg, getLatestChannelMessageQuery, chID)
	if err != nil {
		return domain.RawMessage{}, err
	}
	return msg, nil
}

func (r *rawMessageRepository) GetAndProcessedChannelMessages(ctx context.Context, tx sqltransaction.Txx, limit int) ([]domain.RawMessage, error) {
	var msgs []domain.RawMessage
	err := tx.SelectContext(ctx, &msgs, getAndProcessMessagesQuery, limit)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// GetMessagesByFilters returns messages that match the given filters
// Each filter contains channel_id and telegram_message_id pair, and both conditions must be satisfied
func (r *rawMessageRepository) GetMessagesByFilters(ctx context.Context, filters []*parserv1.FiltersRawMessages) ([]domain.RawMessage, error) {

	if len(filters) == 0 {
		return []domain.RawMessage{}, nil
	}

	args := make([]any, 0, len(filters)*2)
	placeholders := make([]string, len(filters))

	for i, filter := range filters {
		offset := i * 2
		placeholders[i] = fmt.Sprintf("($%d, $%d)", offset+1, offset+2)
		args = append(args, filter.ChannelID, filter.TgMsgId)
	}

	query := fmt.Sprintf(getMessagesByFiltersQuery, strings.Join(placeholders, ","))

	var msgs []domain.RawMessage
	err := r.db.SelectContext(ctx, &msgs, query, args...)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
