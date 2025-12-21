package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func New() (*DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := "host=" + host +
		" port=" + port +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" sslmode=disable"

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) Close() error {
	return db.conn.Close()
}
