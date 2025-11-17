package store

import (
	"context"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/internal/store/channel_repository"
	"github.com/Negat1v9/sum-tel/services/core/internal/store/subscription_repository"
	"github.com/Negat1v9/sum-tel/services/core/internal/store/user_repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) (*model.User, error)
}

type ChannelRepository interface {
	Create(ctx context.Context, channel *model.Channel) (*model.Channel, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Channel, error)
	GetByUsername(ctx context.Context, username string) (*model.Channel, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Channel, error)
	GetUsernamesForParse(ctx context.Context, limit, offset int) ([]model.Channel, error)
	Update(ctx context.Context, channel *model.Channel) (*model.Channel, error)
	Delete(ctx context.Context, id uuid.UUID) (*model.Channel, error)
}

type UserChannelSubscriptionRepository interface {
	Create(ctx context.Context, sub *model.UserSubscription) (*model.UserSubscription, error)
	GetByID(ctx context.Context, id int64) (*model.UserSubscription, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.UserSubscription, error)
	Update(ctx context.Context, sub *model.UserSubscription) (*model.UserSubscription, error)
	Delete(ctx context.Context, id int64) (*model.UserSubscription, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserSubscription, error)
}

type Storage struct {
	UserRepo    UserRepository
	ChannelRepo ChannelRepository
	SubRepo     UserChannelSubscriptionRepository
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		UserRepo:    user_repository.NewUserRepository(db),
		ChannelRepo: channel_repository.NewChannelRepository(db),
		SubRepo:     subscription_repository.NewUserSubscriptionRepository(db),
	}
}
