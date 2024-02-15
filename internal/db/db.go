package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db DB

type SQLOperations interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DB interface {
	SQLOperations
	Close() error
	Ping() error
}

type AppDB struct {
	*sql.DB
	valid bool
}

func InitDB() DB {
	return InitDBWithURL(
		os.Getenv("DATABASE_URL"),
	)
}

func InitDBWithURL(databaseURL string) DB {
	if databaseURL == "" {
		fmt.Println("database url is required")
	}

	dB, err := sql.Open("postgres", databaseURL)
	if err != nil {
		fmt.Printf("sql.Open failed: %v", err)
	}

	db = &AppDB{
		DB:    dB,
		valid: true,
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("db ping failed: %v", err)
	}

	return db
}
