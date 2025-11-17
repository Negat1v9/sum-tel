package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

func NewPostgresConn(host string, port int, user, password, dbname string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dbx, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return dbx, nil
}
