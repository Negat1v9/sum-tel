package service

import (
	"context"
	"fmt"

	parserv1 "github.com/Negat1v9/sum-tel/services/parser/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/services/parser/internal/model"
	tgparser "github.com/Negat1v9/sum-tel/services/parser/internal/parser/tgParser"
	"github.com/Negat1v9/sum-tel/services/parser/internal/store"
)

const (
	DefaultTimePerMessage = 30 // base time in minutes
)

type ParserService struct {
	parser  *tgparser.TgParser
	storage *store.Store
}

func NewParserService(parser *tgparser.TgParser, storage *store.Store) *ParserService {
	return &ParserService{parser: parser, storage: storage}
}

func (s *ParserService) ParseNewChannel(ctx context.Context, channelID string, username string) (*parserv1.NewChannelResponse, error) {
	r, err := s.parser.ParseChannel(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("messages.MessageService.ParseNewChannel: %w", err)
	}

	if len(r.Messages) == 0 {
		return nil, fmt.Errorf("messages.MessageService.ParseNewChannel: no messages found for channel %s", username)
	}

	tx, err := s.storage.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	// commit or rollback transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// TODO: save channels messages to DB
	err = s.storage.RawMsgRepo().CreateMessages(ctx, tx, convertToModel(channelID, r.Messages))
	if err != nil {
		return nil, err
	}

	return &parserv1.NewChannelResponse{
		Success:     true,
		Username:    r.Username,
		Name:        r.Name,
		Description: r.Description,
		MsgInterval: int32(calculateParseTime(r.Messages)),
	}, nil
}

func (s *ParserService) ParseMessages(ctx context.Context, channelID string, username string) (*parserv1.ParseMessagesResponse, error) {
	// get latest message from DB to recognize from which message to parse new ones
	latestMsg, err := s.storage.RawMsgRepo().GetLatestChannelMessage(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("messages.MessageService.ParseMessages: %w", err)
	}

	// parese new messages from telegram

	msgs, err := s.parser.ParseMessages(ctx, username, latestMsg.TelegramMessageID)
	if err != nil {
		return nil, fmt.Errorf("messages.MessageService.ParseMessages: %w", err)
	}

	tx, err := s.storage.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	// commit or rollback transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = s.storage.RawMsgRepo().CreateMessages(ctx, tx, convertToModel(channelID, msgs))
	if err != nil {
		return nil, err
	}

	return &parserv1.ParseMessagesResponse{
		Success:     true,
		MsgInterval: int32(calculateParseTime(msgs)),
	}, nil
}

// calculateParseTime calculates the time required to parse messages based on their content in minutes
func calculateParseTime(msgs []tgparser.ParsedMessage) int {
	// const baseTimePerMessage = 1 // base time in minutes
	l := len(msgs)

	// deltaTime := 0
	if l <= 1 {
		return DefaultTimePerMessage
	}

	totalDeltaTime := 0

	for i := 0; i < l-1; i++ {
		deltaTime := msgs[i+1].Date.Sub(msgs[i].Date)
		if deltaTime < 0 {
			// module deltaTime
			deltaTime *= -1
		}
		// accumulate total time between publications (in seconds)
		totalDeltaTime += int(deltaTime.Minutes()) // convert to minutes

	}

	return totalDeltaTime / (l - 1) // convert seconds to minutes
}

func convertToModel(channelID string, msgs []tgparser.ParsedMessage) []model.RawMessage {
	modelMsgs := make([]model.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		modelMsgs = append(modelMsgs, model.NewRawMsg(
			channelID,
			msg.Type,
			msg.MsgId,
			msg.HtmlText,
			msg.PhotoUrls,
			msg.Date,
		))
	}
	return modelMsgs
}
