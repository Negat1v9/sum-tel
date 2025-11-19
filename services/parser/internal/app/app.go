package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/Negat1v9/sum-tel/services/parser/internal/parser/service"
	tgparser "github.com/Negat1v9/sum-tel/services/parser/internal/parser/tgParser"
	"github.com/Negat1v9/sum-tel/services/parser/internal/server"
	storage "github.com/Negat1v9/sum-tel/services/parser/internal/store"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(cfg *config.GrpcServerConfig, db *sqlx.DB) *App {

	tgParser := tgparser.NewTgParser()

	msgsService := service.NewParserService(tgParser, storage.NewStorage(db))

	gRPCServer := grpc.NewServer()

	server.Register(gRPCServer, msgsService)

	return &App{
		log:        slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		gRPCServer: gRPCServer,
		port:       cfg.GPRPCPort,
	}
}

// run grpc server on specified port
func (a *App) Run() error {
	// start listen tcp socket
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	a.log.Info("gRPC server started", slog.Int("port", a.port))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	return nil
}

// stop grpc server
func (a *App) Stop() {
	a.log.Info("gRPC server stopped", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
