package channel_repository

import (
	"context"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/shared/sqltransaction"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChannelRepository struct {
	db *sqlx.DB
}

func NewChannelRepository(db *sqlx.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(ctx context.Context, tx sqltransaction.Txx, channel *model.Channel) (*model.Channel, error) {
	row := tx.QueryRowxContext(
		ctx,
		createChannelQuery,
		channel.ID,
		channel.Status,
		channel.CreatedBy,
		channel.Username,
		channel.Title,
		channel.Description,
		channel.ParseInterval,
	)

	if err := row.StructScan(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Channel, error) {
	channel := &model.Channel{}
	err := r.db.GetContext(ctx, channel, getChannelByIDQuery, id)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) GetByUsername(ctx context.Context, username string) (*model.Channel, error) {
	channel := &model.Channel{}
	err := r.db.GetContext(ctx, channel, getChannelByUsernameQuery, username)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Channel, error) {
	channels := []model.Channel{}
	err := r.db.SelectContext(ctx, &channels, getAllChannelsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

// avgMsgs - average number of messages parsed per interval
func (r *ChannelRepository) GetUsernamesForParse(ctx context.Context, avgMsgs int, limit, offset int) ([]model.Channel, error) {
	channels := []model.Channel{}
	err := r.db.SelectContext(ctx, &channels, getUsernamesForParseQuery, limit, offset, avgMsgs)
	if err != nil {
		return nil, err
	}
	// JUST EDIT FOR TESTING DOCKDER MUGAGAGAGA

	return channels, nil
}

func (r *ChannelRepository) Update(ctx context.Context, channel *model.Channel) (*model.Channel, error) {
	row := r.db.QueryRowxContext(
		ctx,
		updateChannelQuery,
		channel.ID,
		channel.Username,
		channel.Title,
		channel.Description,
		channel.ParseInterval,
		channel.LastParsedAt,
		channel.Status,
	)

	if err := row.StructScan(channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *ChannelRepository) Delete(ctx context.Context, id uuid.UUID) (*model.Channel, error) {
	channel := &model.Channel{}
	err := r.db.GetContext(ctx, channel, deleteChannelQuery, id)
	if err != nil {
		return nil, err
	}

	return channel, nil
}
