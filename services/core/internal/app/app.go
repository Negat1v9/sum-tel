package app

import (
	channelservice "github.com/Negat1v9/sum-tel/services/core/internal/channel/service"
	grpcclient "github.com/Negat1v9/sum-tel/services/core/internal/grpc/client"
	"github.com/Negat1v9/sum-tel/services/core/internal/server"
	"github.com/Negat1v9/sum-tel/services/core/internal/store"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/logger"
	"github.com/Negat1v9/sum-tel/shared/postgres"
)

type App struct {
	cfg *config.CoreConfig
	log *logger.Logger
}

func NewApp(cfg *config.CoreConfig, log *logger.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run() error {
	// Application run logic goes here
	// initialize server, database connections, services, etc.
	db, err := postgres.NewPostgresConn(a.cfg.PostgresConfig.DbHost, a.cfg.PostgresConfig.DbPort, a.cfg.PostgresConfig.DbUser, a.cfg.PostgresConfig.DbPassword, a.cfg.PostgresConfig.DbName)
	if err != nil {
		return err
	}
	a.log.Infof("connect to postgres host: %s, port %d", a.cfg.PostgresConfig.DbHost, a.cfg.PostgresConfig.DbPort)

	storage := store.NewStorage(db)

	tgParsergRPCClient, err := grpcclient.NewTgParserClient(a.cfg.GrpcClientConfig.URL, a.log)
	if err != nil {
		return err
	}

	// services:
	channelService := channelservice.NewChannelService(storage, tgParsergRPCClient)

	server := server.New(a.cfg, a.log)
	server.MapHandlers(channelService)

	a.log.Infof("server is starting on %s", a.cfg.WebConfig.ListenAddress)
	// TODO: add graceful shutdown
	if err := server.Run(); err != nil {
		a.log.Errorf("server is stopped: %v", err)
		return err
	}

	return nil

}
