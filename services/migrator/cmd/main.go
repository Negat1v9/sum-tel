package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// load core config
	coreCfg, err := config.LoadCoreConfig("./config/core-config")
	if err != nil {
		log.Printf("[error] load core config: %v", err)
		os.Exit(1)
	}

	// load parser config
	parserCfg, err := config.LoadParserConfig("./config/parser-config")
	if err != nil {
		log.Printf("[error] load parser config: %v", err)
		os.Exit(1)
	}

	// core migrations
	coreDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		coreCfg.PostgresConfig.DbUser,
		coreCfg.PostgresConfig.DbPassword,
		coreCfg.PostgresConfig.DbHost,
		coreCfg.PostgresConfig.DbPort,
		coreCfg.PostgresConfig.DbName,
		coreCfg.PostgresConfig.DbSslMode,
	)
	log.Printf("[debug] core DSN: %s", maskedDSN(coreDSN))

	m, err := migrate.New("file://core-migrations", coreDSN)
	if err != nil {
		log.Printf("[error] core-migrations migrate.New: %v", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("[warn] core-migrations: no migrations to apply")
		} else {
			log.Printf("[error] core-migrations migrate.Up: %v", err)
			os.Exit(1)
		}
	}

	m.Close()
	log.Printf("[info] core-migrations applied successfully")

	// parser migrations
	parserDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		parserCfg.DbUser,
		parserCfg.DbPassword,
		parserCfg.DbHost,
		parserCfg.DbPort,
		parserCfg.DbName,
		parserCfg.DbSslMode,
	)
	log.Printf("[debug] parser DSN: %s", maskedDSN(parserDSN))

	m, err = migrate.New("file://parser-migrations", parserDSN)
	if err != nil {
		log.Printf("[error] parser-migrations migrate.New: %v", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("[warn] parser-migrations: no migrations to apply")
		} else {
			log.Printf("[error] parser-migrations migrate.Up: %v", err)
			os.Exit(1)
		}
	}

	m.Close()
	log.Printf("[info] parser-migrations applied successfully")
	os.Exit(0)
}

func maskedDSN(dsn string) string {
	if len(dsn) > 20 {
		return dsn[:20] + "***"
	}
	return dsn
}
