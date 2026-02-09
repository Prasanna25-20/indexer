package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(),
		"postgres://postgres:password@localhost:5432/indexer?sslmode=disable")
}
