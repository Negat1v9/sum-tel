package subscription_repository

import (
	"context"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/shared/sqltransaction"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserSubscriptionRepository struct {
	db *sqlx.DB
}

func NewUserSubscriptionRepository(db *sqlx.DB) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{db: db}
}

func (r *UserSubscriptionRepository) Create(ctx context.Context, tx sqltransaction.Txx, sub *model.UserSubscription) (*model.UserSubscription, error) {
	row := tx.QueryRowxContext(
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

func (r *UserSubscriptionRepository) GetByUserAndChannelID(ctx context.Context, userID int64, channelID uuid.UUID) (*model.UserSubscription, error) {
	sub := &model.UserSubscription{}
	err := r.db.GetContext(ctx, sub, getSubscriptionByUserAndChannelQuary, userID, channelID)
	if err != nil {
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

func (r *UserSubscriptionRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.UserSubscriptionWithChannel, error) {

	rows, err := r.db.QueryxContext(ctx, getSubscriptionsByUserIDQuery, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	subsInfo := make([]model.UserSubscriptionWithChannel, 0)
	for rows.Next() {
		var subInfo model.UserSubscriptionWithChannel
		var ch model.Channel

		if err = rows.Scan(&subInfo.ID, &subInfo.UserID, &subInfo.ChannelID, &subInfo.SubscribedAt, &ch.ID,
			&ch.Username, &ch.Title, &ch.Description, &ch.ParseInterval, &ch.LastParsedAt, &ch.CreatedAt,
			&ch.UpdatedAt,
		); err != nil {
			return nil, err
		}

		subInfo.Channel = ch

		subsInfo = append(subsInfo, subInfo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return subsInfo, nil
}

func (r *UserSubscriptionRepository) Delete(ctx context.Context, id int64) (*model.UserSubscription, error) {
	sub := &model.UserSubscription{}
	err := r.db.GetContext(ctx, sub, deleteSubscriptionQuery, id)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
