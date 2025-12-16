package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

const (
	maxOpenConn     = 20
	connMaxLife     = 10
	maxIdleConns    = 6
	connMaxIdleTime = 2
)

func NewPostgresConn(host string, port int, user, password, dbname string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbx, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = dbx.Ping(); err != nil {
		return nil, err
	}

	dbx.SetMaxOpenConns(maxOpenConn)
	dbx.SetConnMaxLifetime(connMaxLife)
	dbx.SetMaxIdleConns(maxIdleConns)
	dbx.SetConnMaxIdleTime(connMaxIdleTime)

	return dbx, nil
}
