package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type DB struct {
	conn *sqlx.DB
}

func NewDB() (*DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found in InitDB")
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	sslmode := os.Getenv("db_sslmode")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, dbHost, dbPort, dbName, sslmode)

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigrations(connStr); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	if _, err := db.Exec("SET search_path TO subscriptions, public"); err != nil {
		return nil, fmt.Errorf("failed to set search_path: %w", err)
	}

	log.Println("DB connection is established =)")

	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetDB() *sqlx.DB {
	return db.conn
}
