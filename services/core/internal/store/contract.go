package store

import (
	"context"
	"errors"

	"github.com/Negat1v9/sum-tel/shared/sqltransaction"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/internal/store/channel_repository"
	newsrepository "github.com/Negat1v9/sum-tel/services/core/internal/store/news_repository"
	"github.com/Negat1v9/sum-tel/services/core/internal/store/subscription_repository"
	"github.com/Negat1v9/sum-tel/services/core/internal/store/user_repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoUserSubscriptions = errors.New("no user subscriptions found")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) (*model.User, error)
}

type ChannelRepository interface {
	Create(ctx context.Context, tx sqltransaction.Txx, channel *model.Channel) (*model.Channel, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Channel, error)
	GetByUsername(ctx context.Context, username string) (*model.Channel, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Channel, error)
	GetUsernamesForParse(ctx context.Context, avgMsgs int, limit, offset int) ([]model.Channel, error)
	Update(ctx context.Context, channel *model.Channel) (*model.Channel, error)
	Delete(ctx context.Context, id uuid.UUID) (*model.Channel, error)
}

type UserChannelSubscriptionRepository interface {
	Create(ctx context.Context, tx sqltransaction.Txx, sub *model.UserSubscription) (*model.UserSubscription, error)
	GetByID(ctx context.Context, id int64) (*model.UserSubscription, error)
	GetByUserAndChannelID(ctx context.Context, userID int64, channelID uuid.UUID) (*model.UserSubscription, error)
	Delete(ctx context.Context, id int64) (*model.UserSubscription, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserSubscriptionWithChannel, error)
}

type NewsRepository interface {
	Create(ctx context.Context, tx sqltransaction.Txx, news *model.News) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.News, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.News, error)
	Delete(ctx context.Context, id uuid.UUID) (*model.News, error)
	CreateNewsSource(ctx context.Context, tx sqltransaction.Txx, source *model.NewsSource) error
	CreateNewsSources(ctx context.Context, tx sqltransaction.Txx, sources []model.NewsSource) error
	DeleteNewsSource(ctx context.Context, id int) (*model.NewsSource, error)
	DeleteNewsSourcesByNewsID(ctx context.Context, newsID uuid.UUID) error
}

type Storage struct {
	db          *sqlx.DB
	sqlTx       sqltransaction.SqlTx
	userRepo    UserRepository
	channelRepo ChannelRepository
	subRepo     UserChannelSubscriptionRepository
	newsRepo    NewsRepository
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db:          db,
		sqlTx:       sqltransaction.NewSqlTransaction(db),
		userRepo:    user_repository.NewUserRepository(db),
		channelRepo: channel_repository.NewChannelRepository(db),
		subRepo:     subscription_repository.NewUserSubscriptionRepository(db),
		newsRepo:    newsrepository.NewNewsRepository(db),
	}
}

func (s *Storage) UserRepo() UserRepository {
	if s.userRepo == nil {
		s.userRepo = user_repository.NewUserRepository(s.db)
	}
	return s.userRepo
}

func (s *Storage) ChannelRepo() ChannelRepository {
	if s.channelRepo == nil {
		s.channelRepo = channel_repository.NewChannelRepository(s.db)
	}
	return s.channelRepo
}

func (s *Storage) SubRepo() UserChannelSubscriptionRepository {
	if s.subRepo == nil {
		s.subRepo = subscription_repository.NewUserSubscriptionRepository(s.db)
	}
	return s.subRepo
}

func (s *Storage) NewsRepo() NewsRepository {
	if s.newsRepo == nil {
		s.newsRepo = newsrepository.NewNewsRepository(s.db)
	}
	return s.newsRepo
}

func (s *Storage) Transaction(ctx context.Context) (sqltransaction.Txx, error) {
	tx, err := s.sqlTx.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
