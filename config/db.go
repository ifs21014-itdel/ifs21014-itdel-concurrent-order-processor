package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Postgres driver
)

func NewDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// cek koneksi
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
