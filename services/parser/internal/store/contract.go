package store

import (
	"context"

	"github.com/Negat1v9/sum-tel/services/parser/internal/model"
	rawMessageRepository "github.com/Negat1v9/sum-tel/services/parser/internal/store/rawMessageRepository"
	sqltransaction "github.com/Negat1v9/sum-tel/services/parser/internal/store/sqlTransaction"
	"github.com/jmoiron/sqlx"
)

type RawMsgRepository interface {
	CreateMessages(ctx context.Context, tx sqltransaction.Txx, msgs []model.RawMessage) error
	GetLatestChannelMessage(ctx context.Context, chID string) (model.RawMessage, error)
	GetChannelMessages(ctx context.Context, chID string, limit, offset int64) ([]model.RawMessage, error)
	// return messages sorted ASC on RawMessage.MessageDate and update status this messages
	GetAndProcessedChannelMessages(ctx context.Context, tx sqltransaction.Txx, chID string, limit int64) ([]model.RawMessage, error)
}

type Store struct {
	db         *sqlx.DB
	rawMsgRepo RawMsgRepository
	sqlTxx     sqltransaction.SqlTx
}

func NewStorage(db *sqlx.DB) *Store {
	return &Store{db: db,
		rawMsgRepo: rawMessageRepository.NewRawMessageRepository(db),
		sqlTxx:     sqltransaction.NewSqlTransaction(db),
	}
}

func (s *Store) RawMsgRepo() RawMsgRepository {
	if s.rawMsgRepo == nil {
		s.rawMsgRepo = rawMessageRepository.NewRawMessageRepository(s.db)
	}
	return s.rawMsgRepo
}

func (s *Store) Transaction(ctx context.Context) (sqltransaction.Txx, error) {
	tx, err := s.sqlTxx.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
