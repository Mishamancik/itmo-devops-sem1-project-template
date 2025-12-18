package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

// Пока в явном виде в коде
func New() (*DB, error) {
	dsn := fmt.Sprintf(
		"host=localhost port=5432 user=validator password=val1dat0r dbname=project-sem-1 sslmode=disable",
	)

	conn, err := sql.Open("postgres", dsn)
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
