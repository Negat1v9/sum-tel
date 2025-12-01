package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	aihttpclient "github.com/Negat1v9/sum-tel/services/narrator/internal/infrastructure/aiHttpClient"
	"github.com/Negat1v9/sum-tel/services/narrator/internal/service/narrator"

	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/kafka/consumer"
	"github.com/Negat1v9/sum-tel/shared/kafka/producer"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

func main() {

	cfg, err := config.LoadNarratorConfig("./config/narrator-config")
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.AppConfig.Env)

	rawMsgConsumer := consumer.NewConsumer(log, []string{cfg.RawMessageConsumerCfg.Broker}, cfg.RawMessageConsumerCfg.Topic, cfg.RawMessageConsumerCfg.GroupID, false)

	AggregatorProducer := producer.NewProducer(log, []string{cfg.NarrationProducerCfg.Broker}, 1, cfg.NarrationProducerCfg.Topic)

	shutdown, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wg := &sync.WaitGroup{}

	aiClient := aihttpclient.NewClient(&cfg.AiClientCfg, false)
	service := narrator.NewService(log, aiClient, rawMsgConsumer, AggregatorProducer, wg)

	go rawMsgConsumer.ProcessMessages(shutdown, service.RawMessagesHandler(), 15)

	<-shutdown.Done()
	log.Infof("Shutting down narrator service...")

	wg.Wait()

	log.Infof("Narrator service stopped")

	// if err := rawMsgConsumer.Close(); err != nil {
	// 	log.Errorf("Error closing raw message consumer: %v", err)
	// }
}
