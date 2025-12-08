package newsservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	grpcclient "github.com/Negat1v9/sum-tel/services/core/internal/grpc/client"
	parserv1 "github.com/Negat1v9/sum-tel/services/core/internal/grpc/proto"
	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/internal/store"
	"github.com/Negat1v9/sum-tel/shared/kafka/consumer"
	"github.com/Negat1v9/sum-tel/shared/logger"
	"github.com/Negat1v9/sum-tel/shared/sqltransaction"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type NewsService struct {
	log          *logger.Logger
	store        *store.Storage
	newsConsumer *consumer.Consumer
	grpcClient   *grpcclient.TgParserClient
}

func NewNewsService(log *logger.Logger, store *store.Storage, newsConsumer *consumer.Consumer, grpcClient *grpcclient.TgParserClient) *NewsService {
	return &NewsService{
		log:          log,
		store:        store,
		newsConsumer: newsConsumer,
		grpcClient:   grpcClient,
	}
}

// TODO: test
func (s *NewsService) News(ctx context.Context, userID int, limit, offset int) (*model.NewsList, error) {
	mn := "NewsService.News"
	userNews, err := s.store.NewsRepo().GetByUserSubscription(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	return userNews, nil
}

// full info about news sources by news id: text, link, published at, etc
func (s *NewsService) NewsSourcesInfo(ctx context.Context, sNewsID string) (*model.NewsSourcesListResponse, error) {
	mn := "NewsService.NewsSourcesInfo"
	newsID, err := uuid.Parse(sNewsID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	newsSources, err := s.store.NewsRepo().GetNewsSourcesByNewsID(ctx, newsID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	sourceRawMessages, err := s.grpcClient.GetNewsSources(ctx, &parserv1.NewsSourcesRequest{
		Filters: convertSourcesToGrpcFilters(newsSources),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mn, err)
	}

	return &model.NewsSourcesListResponse{
		NewsID:  sNewsID,
		Sources: model.ClearGrpcNewsSources(sourceRawMessages.Messages),
	}, nil
}

// proccess news from broker kafka save to db
func (s *NewsService) ProcessNewsHandler() consumer.ProcFunc {
	return func(shutDownCtx context.Context, msgs []kafka.Message) bool {
		mn := "NewsService.ProcessNewsHandler"

		kafkaNews := make([]model.NewsKafkaType, 0, len(msgs))
		s.log.Debugf("Received message: len %d", len(msgs))

		for _, msg := range msgs {

			if shutDownCtx.Err() != nil {
				return false
			}

			var news model.NewsKafkaType

			err := json.Unmarshal(msg.Value, &news)
			if err != nil {
				s.log.Errorf("%s: %v", mn, err)
				continue
			}
			kafkaNews = append(kafkaNews, news)
		}

		for _, knews := range kafkaNews {
			ctx, cancel := context.WithTimeout(shutDownCtx, 10*time.Second)
			if err := s.saveWithRetrys(ctx, model.ConvertNewsKafkaTypeToNews(knews), knews.Sources); err != nil {
				s.log.Errorf("%s: %v", mn, err)
			}
			cancel()
		}

		return true
	}
}

func (s *NewsService) saveWithRetrys(ctx context.Context, newNews *model.News, sources []model.SourceKafkaType) error {
	// call on error for not copy it every time then error
	// stack error contains all error in loop if exists
	errorsStack := make([]error, 0, 3)
	maxRetrys := 3

	// onErrFn appends error to errorsStack and rollbacks transaction if exists
	onErrFn := func(tx sqltransaction.Txx, err error) {
		errorsStack = append(errorsStack, err)
		if tx != nil {
			tx.Rollback()
		}
	}
	for attempt := 1; attempt <= maxRetrys; attempt++ {
		tx, err := s.store.Transaction(ctx)
		if err != nil {
			onErrFn(nil, fmt.Errorf("saveWithRetrys: cannot start transaction: %w", err))
			continue
		}
		// save news
		err = s.store.NewsRepo().Create(ctx, tx, newNews)
		if err != nil {
			onErrFn(tx, fmt.Errorf("saveWithRetrys: cannot save news: %w", err))
			continue
		}
		// save news sources
		err = s.store.NewsRepo().CreateNewsSources(ctx, tx, model.ConvertSourceKafkaTypesToNewsSources(newNews.ID, sources))
		if err != nil {
			onErrFn(tx, fmt.Errorf("saveWithRetrys: cannot save news sources: %w", err))
			continue
		}
		// successful save, commit transaction and return
		if err := tx.Commit(); err != nil {
			onErrFn(tx, fmt.Errorf("saveWithRetrys: cannot commit transaction: %w", err))
		}
	}

	return nil
}

func convertSourcesToGrpcFilters(sources []model.NewsSource) []*parserv1.FiltersRawMessages {
	grpcSources := make([]*parserv1.FiltersRawMessages, 0, len(sources))
	for _, s := range sources {
		grpcSources = append(grpcSources, &parserv1.FiltersRawMessages{
			ChannelID: s.ChannelID.String(),
			TgMsgId:   s.MessageID,
			Username:  s.ChannelName,
		})
	}
	return grpcSources
}
