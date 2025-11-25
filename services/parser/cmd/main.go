package main

import (
	grpcapp "github.com/Negat1v9/sum-tel/services/parser/internal/app"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/postgres"
)

func main() {
	cfg, err := config.LoadParserConfig("./config/parser-config")
	if err != nil {
		panic(err)
	}

	db, err := postgres.NewPostgresConn(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	if err != nil {
		// log.Fatalf("failed to connect to database: %w`", err)
	}

	app := grpcapp.New(cfg, db)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
