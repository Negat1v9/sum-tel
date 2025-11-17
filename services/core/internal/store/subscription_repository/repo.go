package subscription_repository

import (
	"context"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserSubscriptionRepository struct {
	db *sqlx.DB
}

func NewUserSubscriptionRepository(db *sqlx.DB) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{db: db}
}

func (r *UserSubscriptionRepository) Create(ctx context.Context, sub *model.UserSubscription) (*model.UserSubscription, error) {
	row := r.db.QueryRowxContext(
		ctx,
		createSubscriptionQuery,
		sub.UserID,
		sub.ChannelID,
	)

	if err := row.StructScan(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *UserSubscriptionRepository) GetByID(ctx context.Context, id int64) (*model.UserSubscription, error) {
	sub := &model.UserSubscription{}
	err := r.db.GetContext(ctx, sub, getSubscriptionByIDQuery, id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *UserSubscriptionRepository) GetAll(ctx context.Context, limit, offset int) ([]model.UserSubscription, error) {
	subs := []model.UserSubscription{}
	err := r.db.SelectContext(ctx, &subs, getAllSubscriptionsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *UserSubscriptionRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserSubscription, error) {
	subs := []model.UserSubscription{}
	err := r.db.SelectContext(ctx, &subs, getSubscriptionsByUserIDQuery, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *UserSubscriptionRepository) Update(ctx context.Context, sub *model.UserSubscription) (*model.UserSubscription, error) {
	row := r.db.QueryRowxContext(
		ctx,
		updateSubscriptionQuery,
		sub.ID,
		sub.UserID,
		sub.ChannelID,
	)

	if err := row.StructScan(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *UserSubscriptionRepository) Delete(ctx context.Context, id int64) (*model.UserSubscription, error) {
	sub := &model.UserSubscription{}
	err := r.db.GetContext(ctx, sub, deleteSubscriptionQuery, id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
