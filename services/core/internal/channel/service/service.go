package channelservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	grpcclient "github.com/Negat1v9/sum-tel/services/core/internal/grpc/client"
	parserv1 "github.com/Negat1v9/sum-tel/services/core/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/internal/store"
	"github.com/google/uuid"
)

const (
	gRPCCallContextTimeout = time.Second * 30
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

func (s *ChannelService) CreateChannel(ctx context.Context, userID int, username string) (res *model.Channel, err error) {
	mn := "ChannelService.CreateChannel"
	// delete "@" if exists
	username, _ = strings.CutPrefix(username, "@")
	// check channel exists
	existsChannel, err := s.store.ChannelRepo().GetByUsername(ctx, username)
	if err == nil {
		return existsChannel, nil
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("%s.GetByUsername: %w", mn, err)
	}

	// create new channel ID
	channelID := uuid.New()

	ctxWithoutCancel := context.WithoutCancel(ctx)

	// parse channel via grpc
	grpcCtx, grpcCancel := context.WithTimeout(ctxWithoutCancel, gRPCCallContextTimeout)
	defer grpcCancel()

	parsedChannel, err := s.grpcClient.ParseNewChannel(grpcCtx, &parserv1.NewChannelRequest{
		ChannelID: channelID.String(),
		Username:  username,
	})
	if err != nil {
		return nil, fmt.Errorf("%s.ParseNewChannel: %w", mn, err)
	}
	// create another context Need save channel
	dbCtx, dbCancel := context.WithTimeout(ctxWithoutCancel, time.Second*60)
	tx, err := s.store.Transaction(dbCtx)
	if err != nil {
		dbCancel()
		return nil, fmt.Errorf("%s.Transaction: %w", mn, err)
	}

	defer func() {
		var txErr error
		if err != nil {
			txErr = tx.Rollback()

		} else {
			txErr = tx.Commit()
		}
		if txErr != nil {
			fmt.Println(txErr)
		}
		dbCancel()
	}()

	// add channel information in db
	createdChannel, err := s.store.ChannelRepo().Create(
		dbCtx, tx, model.NewChannel(channelID, model.ChannelStatusActive, username, parsedChannel.Name,
			parsedChannel.Description, userID, int(parsedChannel.GetMsgInterval()), time.Now()))
	if err != nil {
		return nil, fmt.Errorf("%s.Create: %w", mn, err)
	}

	return createdChannel, nil
}

// return full info about channel by username
// username without @
func (s *ChannelService) GetChannelByUsername(ctx context.Context, username string) (*model.Channel, error) {
	mn := "ChannelService.GetChannelByUsername"
	channel, err := s.store.ChannelRepo().GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%s.GetByUsername: %w", mn, err)
	}
	return channel, nil
}

func (s *ChannelService) SubscribeChannel(ctx context.Context, userID int, channelID string) (*model.UserSubscription, error) {
	mn := "ChannelService.SubscribeChannel"
	chID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("%s invalid channel ID: %w", mn, err)
	}

	_, err = s.store.ChannelRepo().GetByID(ctx, chID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	existsSub, err := s.store.SubRepo().GetByUserAndChannelID(ctx, userID, chID)
	switch {
	case existsSub != nil && err == nil:
		return existsSub, nil
	case err != nil && err != sql.ErrNoRows:
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	tx, err := s.store.Transaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s.Transaction: %w", mn, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	createdSub, err := s.store.SubRepo().Create(ctx, tx, model.NewSub(userID, chID))
	if err != nil {
		return nil, fmt.Errorf("%s.Create: %w", mn, err)
	}

	return createdSub, nil
}

// get all user subscriptions
func (s *ChannelService) UsersSubscriptions(ctx context.Context, userID int, limit, offset int) (*model.UserSubscriptionWithChannelList, error) {
	mn := "ChannelService.UsersSubscriptions"
	subsINfo, err := s.store.SubRepo().GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s.GetByUserID: %w", mn, err)
	}
	return subsINfo, nil
}

func (s *ChannelService) ChannelsToParse(ctx context.Context, limit, offset int) ([]model.Channel, error) {
	mn := "ChannelService.ChannelsToParse"
	channels, err := s.store.ChannelRepo().GetUsernamesForParse(ctx, 10, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Channel{}, nil
		}
		return nil, fmt.Errorf("%s.GetUsernamesForParse: %w", mn, err)
	}
	return channels, nil
}

func (s *ChannelService) ParseChannel(ctx context.Context, channelID uuid.UUID, username string) error {
	mn := "ChannelService.ParseChannel"
	result, err := s.grpcClient.ParseMessages(ctx, &parserv1.ParseMessagesRequest{
		ChannelID: channelID.String(),
		Username:  username,
	})
	if err != nil {
		return fmt.Errorf("%s.ParseMessages: %w", mn, err)
	}
	if !result.Success {
		return fmt.Errorf("%s: failed to parse messages for channel %s", mn, username)
	}

	updChannel := &model.Channel{
		ID:            channelID,
		ParseInterval: int(result.MsgInterval),
		UpdatedAt:     time.Now(),
		LastParsedAt:  time.Now(),
	}
	_, err = s.store.ChannelRepo().Update(ctx, updChannel)
	if err != nil {
		return fmt.Errorf("%s.Update: %w", mn, err)
	}
	return nil

}
