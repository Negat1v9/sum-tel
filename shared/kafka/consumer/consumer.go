package consumer

import (
	"context"
	"time"

	"github.com/Negat1v9/sum-tel/shared/logger"
	"github.com/segmentio/kafka-go"
)

type ProcFunc func(context.Context, []kafka.Message) bool

type Consumer struct {
	consumer *kafka.Reader
	// consumer *kafka.Consumer

	log            *logger.Logger
	pendingOffset  []kafka.Message
	autoCommit     bool          // commit after every readed message
	commitInterval time.Duration // commit
}

func NewConsumer(log *logger.Logger, brokers []string, topic string, groupID string, autoCommit bool) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		consumer:       r,
		log:            log,
		pendingOffset:  make([]kafka.Message, 0),
		autoCommit:     autoCommit,
		commitInterval: time.Minute,
	}
}

// run recieve all messages from kafka topic
func (c *Consumer) ProcessMessages(ctx context.Context, procFunc ProcFunc, batchSize int) {
	batch := make([]kafka.Message, 0, batchSize)

	c.autoCommitWorker(ctx)

	for {
		select {
		case <-ctx.Done():
			c.consumer.Close()
			return

		default:
			msg, err := c.consumer.FetchMessage(ctx)
			if err != nil {
				kerr, ok := err.(kafka.Error)
				// unexpected error
				if !ok {
					c.log.Debugf("Consumer.ProcessMessages error: %v", err)
					continue
				}
				if kerr.Timeout() {
					c.log.Debugf("Consumer.ProcessMessages timeout")
					// no messages it is normal
					continue
				}
				c.log.Debugf("Consumer.ProcessMessages error: %v", err)
				continue
			}
			batch = append(batch, msg)
			if len(batch) >= batchSize {
				// process mesages
				c.processMessage(ctx, batch, procFunc)
				batch = batch[:0]
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, messages []kafka.Message, procFunc ProcFunc) {

	if ok := procFunc(ctx, messages); !ok {
		// processing failed do not commit offsets
		c.log.Warnf("Consumer.processMessage: processing failed, skipping commit")
		return
	}

	if !c.autoCommit {
		// store offset
		c.storeOffsets(messages)
	}

}

func (c *Consumer) commitOffset() error {
	if len(c.pendingOffset) == 0 {
		return nil
	}
	c.log.Debugf("Consumer.commitOffset: committing %d messages", len(c.pendingOffset))

	if err := c.consumer.CommitMessages(context.TODO(), c.pendingOffset...); err != nil {
		return err
	}

	c.pendingOffset = c.pendingOffset[:0]

	return nil
}

func (c *Consumer) storeOffsets(msgs []kafka.Message) {
	c.pendingOffset = append(c.pendingOffset, msgs...)

}

func (c *Consumer) autoCommitWorker(ctx context.Context) {
	ticker := time.NewTicker(c.commitInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Warnf("Consumer.autoCommitWorker: recieve done context")
				if err := c.commitOffset(); err != nil {
					c.log.Errorf("Consumer.autoCommitWorker: %v", err)
				}
				return
			case <-ticker.C:
				if err := c.commitOffset(); err != nil {
					c.log.Errorf("Consumer.autoCommitWorker: %v", err)
				}
			}
		}
	}()
}
