package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	parserv1 "github.com/Negat1v9/sum-tel/services/parser/internal/api/proto"
	"github.com/Negat1v9/sum-tel/services/parser/internal/domain"

	tgparser "github.com/Negat1v9/sum-tel/services/parser/internal/parser/tgParser"
	"github.com/Negat1v9/sum-tel/services/parser/internal/store"
	"github.com/Negat1v9/sum-tel/shared/kafka/producer"
)

const (
	DefaultTimePerMessage = 30 // base time in minutes
)

var (
	// error on parsing new channel and it not have any messages
	ErrChannelNoMessages = errors.New("channel has no messages")
)

type ParserService struct {
	parser  *tgparser.TgParser
	storage *store.Store

	rawMsgKafkaProducer *producer.Producer
}

func NewParserService(parser *tgparser.TgParser, storage *store.Store, rawMsgKafkaProducer *producer.Producer) *ParserService {
	return &ParserService{parser: parser, storage: storage, rawMsgKafkaProducer: rawMsgKafkaProducer}
}

func (s *ParserService) ParseNewChannel(ctx context.Context, channelID string, username string) (*parserv1.NewChannelResponse, error) {
	const mn = "ParserService.ParseNewChannel"
	// parse telegram meesage and receive information about it and latest ~12 messasges
	r, err := s.parser.ParseChannel(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%s.ParseChannel: %w", mn, err)
	}

	if len(r.Messages) == 0 {
		return nil, fmt.Errorf("%s: %w", mn, ErrChannelNoMessages)
	}
	// save in db messages if exitst
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

	err = s.storage.RawMsgRepo().CreateMessages(ctx, tx, convertToModel(channelID, r.Messages))
	if err != nil {
		return nil, err
	}

	err = s.rawMsgKafkaProducer.SendMessage(ctx, domain.ConvetParsedMessagesToAny(channelID, r.Messages)...)
	if err != nil {
		return nil, fmt.Errorf("%s.SendMessage: %w:", mn, err)
	}

	return &parserv1.NewChannelResponse{
		Success:     true,
		Username:    r.Username,
		Name:        r.Name,
		Description: r.Description,
		MsgInterval: DefaultTimePerMessage,
	}, nil
}

func (s *ParserService) ParseMessages(ctx context.Context, channelID string, username string) (*parserv1.ParseMessagesResponse, error) {
	const mn = "ParserService.ParseMessages"
	// get latest message from DB to recognize from which message to parse new ones
	latestMsg, err := s.storage.RawMsgRepo().GetLatestChannelMessage(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	// max count of errors in loop after return break the loop
	countErr := 3
	// stack error contains all error in loop if exists
	errorsStack := make([]error, 0, 3)
	// the number of iterations needed to calculate the average time between messages as a divisor
	numberIterations := 0
	// chennel for result parse from worker
	var msgsInteval int32 = DefaultTimePerMessage

	sleepOnErrTime := time.Second
	// call on error for not copy it every time then error
	onErrFn := func(internalErr error) {
		errorsStack = append(errorsStack, internalErr)
		countErr--
		time.Sleep(sleepOnErrTime)
	}
	lastMsgID := latestMsg.TelegramMessageID
	for {
		if countErr <= 0 {
			return nil, fmt.Errorf("%s: %v", mn, errors.Join(errorsStack...))
		}
		numberIterations++
		// parse telegram
		msgs, err := s.parser.ParseMessages(ctx, username, lastMsgID)
		if err != nil {
			onErrFn(err)
			continue
		}
		// no more messages
		if len(msgs) == 0 {
			break
		}
		// add medium duraion between messages
		msgsInteval += int32(calculateParseTime(msgs))
		// save new messages in db
		tx, err := s.storage.Transaction(ctx)
		if err != nil {
			onErrFn(err)
			continue
		}

		// save messages in db
		err = s.storage.RawMsgRepo().CreateMessages(ctx, tx, convertToModel(channelID, msgs))
		if err != nil {
			tx.Rollback()
			onErrFn(err)
		} else {
			err = tx.Commit()
			if err != nil {
				onErrFn(err)
			}
		}
		if len(msgs) > 0 {
			lastMsgID = msgs[len(msgs)-1].MsgId
		}

	}
	return &parserv1.ParseMessagesResponse{
		Success:     true,
		MsgInterval: DefaultTimePerMessage,
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

func convertToModel(channelID string, msgs []tgparser.ParsedMessage) []domain.RawMessage {
	modelMsgs := make([]domain.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		modelMsgs = append(modelMsgs, domain.NewRawMsg(
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
