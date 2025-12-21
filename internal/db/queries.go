package db

import (
	"context"
	"strconv"
	"time"
)

type Stats struct {
	TotalItems      int
	TotalCategories int
	TotalPrice      int
}

func (db *DB) InsertPrices(ctx context.Context, records [][]string) (Stats, error) {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return Stats{}, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO prices (name, category, price, create_date)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return Stats{}, err
	}
	defer stmt.Close()

	for _, row := range records[1:] {
		price, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return Stats{}, err
		}

		if _, err := stmt.ExecContext(
			ctx,
			row[1],
			row[2],
			price,
			row[4],
		); err != nil {
			return Stats{}, err
		}
	}

	var stats Stats
	err = tx.QueryRowContext(ctx, `
		SELECT
			COUNT(*) AS total_items,
			COUNT(DISTINCT category) AS total_categories,
			COALESCE(SUM(price), 0)::INT AS total_price
		FROM prices
	`).Scan(
		&stats.TotalItems,
		&stats.TotalCategories,
		&stats.TotalPrice,
	)
	if err != nil {
		return Stats{}, err
	}

	if err := tx.Commit(); err != nil {
		return Stats{}, err
	}

	return stats, nil
}

type Price struct {
	ID         int
	Name       string
	Category   string
	Price      float64
	CreateDate time.Time
}

func (db *DB) GetPrices(ctx context.Context) ([]Price, error) {
	rows, err := db.conn.QueryContext(ctx, `
		SELECT id, name, category, price, create_date
		FROM prices
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Price

	for rows.Next() {
		var p Price
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Category,
			&p.Price,
			&p.CreateDate,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
