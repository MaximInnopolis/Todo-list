package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func (db *Database) GetPool() *pgxpool.Pool {
	return db.pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func NewPool(dbUrl string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), dbUrl)
}
