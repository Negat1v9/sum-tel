package narrator

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Negat1v9/sum-tel/services/narrator/internal/domain"
	aihttpclient "github.com/Negat1v9/sum-tel/services/narrator/internal/infrastructure/aiHttpClient"
	"github.com/Negat1v9/sum-tel/shared/kafka/consumer"
	"github.com/Negat1v9/sum-tel/shared/kafka/producer"
	"github.com/Negat1v9/sum-tel/shared/logger"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	log                    *logger.Logger
	aiClient               *aihttpclient.Client
	rawMsgCosnumer         *consumer.Consumer
	newsAggregatorProducer *producer.Producer

	wg *sync.WaitGroup
}

func NewService(log *logger.Logger, aiClient *aihttpclient.Client, consumer *consumer.Consumer, producer *producer.Producer, wg *sync.WaitGroup) *Service {
	s := &Service{
		log:                    log,
		aiClient:               aiClient,
		rawMsgCosnumer:         consumer,
		newsAggregatorProducer: producer,
		wg:                     wg,
	}

	return s
}

func (s *Service) Stop(ctx context.Context) error {
	done := make(chan bool)
	go func() {
		s.wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		s.log.Infof("Service stopped gracefully")
		return s.newsAggregatorProducer.Close() // close producer and flush messages
	case <-ctx.Done():
		s.log.Errorf("Service stop timeout: %v", ctx.Err())
		return ctx.Err()
	}
}

func (s *Service) RawMessagesHandler() consumer.ProcFunc {
	return func(shutDownCtx context.Context, msgs []kafka.Message) bool {
		mn := "Service.RawMessagesHandler"

		rawMsgs := make([]domain.RawMessage, 0, len(msgs))
		s.log.Debugf("Received message: len %d", len(msgs))
		for _, msg := range msgs {

			if shutDownCtx.Err() != nil {
				return false
			}

			var rawMsg domain.RawMessage

			err := json.Unmarshal(msg.Value, &rawMsg)
			if err != nil {
				s.log.Errorf("%s: %v", mn, err)
				continue
			}
			rawMsgs = append(rawMsgs, rawMsg)
		}

		aiClientCtx, cancel := context.WithTimeout(shutDownCtx, time.Second*120)
		defer cancel()
		// proccess raw RawMessage
		aggregatedMsgs, err := s.aiClient.DoAggregation(aiClientCtx, rawMsgs)
		if err != nil {
			s.log.Errorf("%s: %v", mn, err)
			return false
		}

		// pruducerCtx not block main shutdown process
		pruducerCtx, producerCancel := context.WithTimeout(context.Background(), time.Second*30)
		defer producerCancel()
		// send aggregated messages to news-aggregator service after LLM proccesing
		s.wg.Add(1) // try send aggregated message in a separate goroutine if receive shutdown signal
		s.sendAggregatedMessage(pruducerCtx, aggregatedMsgs)
		s.wg.Done()

		return true
	}
}

func (s *Service) sendAggregatedMessage(ctx context.Context, msg *domain.AggregationResponse) error {
	mn := "Service.SendAggregatedMessage"
	err := s.newsAggregatorProducer.SendMessage(ctx, domain.ConvertAggregateNewsdResponseToAny(msg.AggregatedNews)...)
	if err != nil {
		s.log.Errorf("%s: failed to send aggregated news message: %v", mn, err)
		return err
	}
	return nil
}
