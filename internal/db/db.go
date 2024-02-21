package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var db DB

//go:embed migrations/*.sql
var embedMigrations embed.FS

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
		fmt.Printf("sql open failed: %v", err)
	}

	db = &AppDB{
		DB:    dB,
		valid: true,
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("db ping failed: %v", err)
	}

	runDBMigration(dB)

	return db
}

func runDBMigration(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("cannot set goose dialect: ", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Print("db migrated successfully")
}
