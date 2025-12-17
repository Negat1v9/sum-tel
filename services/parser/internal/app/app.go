package grpcapp

import (
	"context"
	"fmt"
	"net"

	server "github.com/Negat1v9/sum-tel/services/parser/internal/api/grpc"
	"github.com/Negat1v9/sum-tel/services/parser/internal/parser/service"
	tgparser "github.com/Negat1v9/sum-tel/services/parser/internal/parser/tgParser"
	storage "github.com/Negat1v9/sum-tel/services/parser/internal/store"
	proccessrawmessage "github.com/Negat1v9/sum-tel/services/parser/internal/workers/proccessRawMessage"
	"github.com/Negat1v9/sum-tel/services/parser/pkg/metrics"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/kafka/producer"
	"github.com/Negat1v9/sum-tel/shared/logger"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type App struct {
	log        *logger.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(cfg *config.ParserServiceConfig, db *sqlx.DB, metrics *metrics.PrometheusMetrics) *App {
	shutDownCtx := context.TODO()

	logger := logger.NewLogger(cfg.AppConfig.Env)

	rawMsgProducer := producer.NewProducer(logger, []string{cfg.RawMessageProducerConfig.Broker}, cfg.RawMessageProducerConfig.BatchSize, cfg.RawMessageProducerConfig.Topic)

	tgParser := tgparser.NewTgParser()

	msgsService := service.NewParserService(logger, tgParser, storage.NewStorage(db), metrics, rawMsgProducer)

	rawMessageWorker := proccessrawmessage.NewWorker(logger, msgsService)

	// check for new raw messages every minute and send to narration processing with kafka producer
	go rawMessageWorker.Start(shutDownCtx)

	gRPCServer := grpc.NewServer()

	server.Register(gRPCServer, msgsService)

	return &App{
		log:        logger,
		gRPCServer: gRPCServer,
		port:       cfg.GrpcServerConfig.GPRPCPort,
	}
}

// run grpc server on specified port
func (a *App) Run() error {
	// start listen tcp socket
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	a.log.Infof("gRPC server started port %d", a.port)

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	return nil
}

// stop grpc server
func (a *App) Stop() {
	a.log.Infof("gRPC server stopped")
	a.gRPCServer.GracefulStop()
}
