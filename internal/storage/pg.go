package storage

import (
	"TestWebServer/internal/config"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	pg *sqlx.DB
}

func InitDatabase(cfg config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", cfg.DB)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewDatabase(conn *sqlx.DB) *Database {
	return &Database{pg: conn}
}
