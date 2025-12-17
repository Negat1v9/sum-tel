package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	parserv1 "github.com/Negat1v9/sum-tel/services/parser/internal/api/proto"
	"github.com/Negat1v9/sum-tel/services/parser/internal/domain"
	"github.com/Negat1v9/sum-tel/services/parser/pkg/metrics"

	tgparser "github.com/Negat1v9/sum-tel/services/parser/internal/parser/tgParser"
	"github.com/Negat1v9/sum-tel/services/parser/internal/store"
	"github.com/Negat1v9/sum-tel/shared/kafka/producer"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

const (
	DefaultTimePerMessage = 30 // base time in minutes
)

var (
	// error on parsing new channel and it not have any messages
	ErrChannelNoMessages = errors.New("channel has no messages")
	ErrNoRawMessages     = errors.New("no raw messages to process")
)

type ParserService struct {
	log     *logger.Logger
	parser  *tgparser.TgParser
	storage *store.Store

	metrics             *metrics.PrometheusMetrics
	rawMsgKafkaProducer *producer.Producer
}

func NewParserService(log *logger.Logger, parser *tgparser.TgParser, storage *store.Store, metrics *metrics.PrometheusMetrics, rawMsgKafkaProducer *producer.Producer) *ParserService {
	return &ParserService{log: log, parser: parser, storage: storage, metrics: metrics, rawMsgKafkaProducer: rawMsgKafkaProducer}
}

func (s *ParserService) ParseNewChannel(ctx context.Context, channelID string, username string) (*parserv1.NewChannelResponse, error) {
	const mn = "ParserService.ParseNewChannel"
	// parse telegram meesage and receive information about it and latest ~12 messasges
	r, err := s.parser.ParseChannel(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%s.ParseChannel: %w", mn, err)
	}

	// add metrics for parsed channel and messages
	s.metrics.IncParsedChannels()

	if len(r.Messages) == 0 {
		return nil, fmt.Errorf("%s: %w", mn, ErrChannelNoMessages)
	}

	s.metrics.AddParsedMessages(username, len(r.Messages))

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

	err = s.storage.RawMsgRepo().CreateMessages(ctx, tx, convertToModel(channelID, domain.RawMessageStatusProcessed, r.Messages))
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
		s.log.Debugf("%s parsing channel %s, lastMessageID: %d", mn, channelID, lastMsgID)

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
		// add metrics for parsed messages
		s.metrics.AddParsedMessages(username, len(msgs))
		// add medium duraion between messages
		msgsInteval += int32(calculateParseTime(msgs))
		// save new messages in db
		tx, err := s.storage.Transaction(ctx)
		if err != nil {
			onErrFn(err)
			continue
		}

		// save messages in db
		err = s.storage.RawMsgRepo().CreateMessages(ctx, tx, convertToModel(channelID, domain.RawMessageStatusNew, msgs))
		if err != nil {
			tx.Rollback()
			onErrFn(err)
		} else {
			err = tx.Commit()
			if err != nil {
				onErrFn(err)
			}
		}

		lastMsgID = msgs[len(msgs)-1].MsgId
	}
	return &parserv1.ParseMessagesResponse{
		Success:     true,
		MsgInterval: DefaultTimePerMessage,
	}, nil
}

func (s *ParserService) GetNewsSources(ctx context.Context, filters []*parserv1.FiltersRawMessages) (*parserv1.NewsSourcesResponse, error) {
	const mn = "ParserService.GetNewsSources"

	// Get messages from database by filters
	rawMessages, err := s.storage.RawMsgRepo().GetMessagesByFilters(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("%s.GetMessagesByFilters: %w", mn, err)
	}

	// Create a map of filters for quick lookup of username by channelID
	filterMap := make(map[string]string)
	for _, filter := range filters {
		filterMap[filter.ChannelID+strconv.FormatInt(filter.TgMsgId, 10)] = filter.Username
	}

	// Convert domain messages to proto messages and populate username field
	messages := make([]*parserv1.Message, 0, len(rawMessages))
	for _, rawMsg := range rawMessages {
		messages = append(messages, &parserv1.Message{
			Type:          rawMsg.ContentType,
			HtmlText:      rawMsg.HTMLText,
			ChannelId:     rawMsg.ChannelID,
			Username:      filterMap[rawMsg.ChannelID+strconv.FormatInt(rawMsg.TelegramMessageID, 10)], // Populate username from filter map
			TelegramMsgId: rawMsg.TelegramMessageID,
			Date:          rawMsg.MessageDate.Unix(),
			PhotoUrls:     rawMsg.MediaURLs,
		})
	}

	return &parserv1.NewsSourcesResponse{
		Success:  true,
		Messages: messages,
	}, nil
}

func (s *ParserService) ProccessRawMessages(ctx context.Context, limit int) error {
	mn := "ParserService.ProccessRawMessages"
	tx, err := s.storage.Transaction(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", mn, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	msgs, err := s.storage.RawMsgRepo().GetAndProcessedChannelMessages(ctx, tx, limit)
	if err != nil {
		return fmt.Errorf("%s: %w", mn, err)
	}

	if len(msgs) == 0 {
		return fmt.Errorf("%s: %w", mn, ErrNoRawMessages)
	}

	return s.rawMsgKafkaProducer.SendMessage(ctx, domain.ConvertRawMessagesToAny(msgs)...)
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

func convertToModel(channelID, status string, msgs []tgparser.ParsedMessage) []domain.RawMessage {
	modelMsgs := make([]domain.RawMessage, 0, len(msgs))
	for _, msg := range msgs {
		modelMsgs = append(modelMsgs, domain.NewRawMsg(
			channelID,
			status,
			msg.Type,
			msg.MsgId,
			msg.HtmlText,
			msg.PhotoUrls,
			msg.Date,
		))
	}
	return modelMsgs
}
