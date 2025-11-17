package channelservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	grpcclient "github.com/Negat1v9/sum-tel/services/core/internal/grpc/client"
	parserv1 "github.com/Negat1v9/sum-tel/services/core/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/internal/store"
	"github.com/google/uuid"
)

type ChannelService struct {
	store      *store.Storage
	grpcClient *grpcclient.TgParserClient
}

func NewChannelService(stor *store.Storage, grpcClient *grpcclient.TgParserClient) *ChannelService {
	return &ChannelService{
		store:      stor,
		grpcClient: grpcClient,
	}
}

func (s *ChannelService) CreateChannel(ctx context.Context, userID int64, username string) (res *model.Channel, err error) {
	// check channel exists
	existsChannel, err := s.store.ChannelRepo.GetByUsername(ctx, username)
	if err == nil {
		return existsChannel, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// create new channel ID
	channelID := uuid.New()

	// parse channel via grpc
	parsedChannel, err := s.grpcClient.ParseNewChannel(ctx, &parserv1.NewChannelRequest{
		ChannelID: channelID.String(),
		Username:  username,
	})
	if err != nil {
		return nil, err
	}

	// TODO: add calculationg parse interval based on channel activity
	createdChannel, err := s.store.ChannelRepo.Create(
		ctx, model.NewChannel(channelID, parsedChannel.Username, parsedChannel.Name,
			parsedChannel.Description, 5, time.Now()))
	if err != nil {
		return nil, err
	}
	return createdChannel, nil
}

// return full info about channel by username
// username without @
func (s *ChannelService) GetChannelByUsername(ctx context.Context, username string) (*model.Channel, error) {
	return s.store.ChannelRepo.GetByUsername(ctx, username)
}

func (s *ChannelService) SubscribeChannel(ctx context.Context, userID int64, channelID string) (*model.UserSubscription, error) {

	chID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}
	// check channel exists
	_, err = s.store.ChannelRepo.GetByID(ctx, chID)
	if err != nil {
		return nil, err
	}
	sub := &model.UserSubscription{
		UserID:       userID,
		ChannelID:    chID,
		SubscribedAt: time.Now(),
	}

	createdSub, err := s.store.SubRepo.Create(ctx, sub)
	return createdSub, err
}

// get all user subscriptions
func (s *ChannelService) UsersSubscriptions(ctx context.Context, userID int64, limit, offset int) ([]model.UserSubscription, error) {
	return s.store.SubRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *ChannelService) ChannelsToParse(ctx context.Context, limit, offset int) ([]model.Channel, error) {
	channels, err := s.store.ChannelRepo.GetUsernamesForParse(ctx, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Channel{}, nil
		}
		return nil, err
	}
	return channels, nil
}

func (s *ChannelService) ParseChannel(ctx context.Context, channelID uuid.UUID, username string) error {
	result, err := s.grpcClient.ParseMessages(ctx, &parserv1.ParseMessagesRequest{
		ChannelID: channelID.String(),
		Username:  username,
	})
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("newsservice.ParseChannel: failed to parse messages for channel %s", username)
	}

	updChannel := &model.Channel{
		ID:            channelID,
		ParseInterval: int(result.MsgInterval), // calucate with
		UpdatedAt:     time.Now(),
		LastParsedAt:  sql.NullTime{Time: time.Now(), Valid: true},
	}
	_, err = s.store.ChannelRepo.Update(ctx, updChannel)
	if err != nil {
		return err
	}
	return nil

}
