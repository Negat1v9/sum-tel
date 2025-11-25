package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Negat1v9/sum-tel/shared/logger"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	topic string

	log      *logger.Logger
	producer *kafka.Writer
}

func NewProducer(log *logger.Logger, brokers []string, batchSize int, topic string) *Producer {
	kConfig := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		BatchSize:    batchSize,
		BatchTimeout: 100 * time.Millisecond,
		MaxAttempts:  3,
		RequiredAcks: int(kafka.RequireOne),
	}
	producer := kafka.NewWriter(kConfig)

	return &Producer{
		producer: producer,
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) SendMessage(ctx context.Context, msgs ...any) error {
	kafkaMsgs := make([]kafka.Message, 0, len(msgs))
	for _, msg := range msgs {
		value, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		kafkaMsgs = append(kafkaMsgs, kafka.Message{Value: value})
	}

	return p.producer.WriteMessages(ctx, kafkaMsgs...)
}
