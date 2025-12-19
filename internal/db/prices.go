package db

import (
	"context"
	"strconv"
	"time"
)

// ===================== Query for POST =====================
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

	categorySet := map[string]struct{}{}
	totalPrice := 0
	inserted := 0

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

    _, err = stmt.ExecContext(
        ctx,
        row[1],
        row[2],
        price,
        row[4],
    )
    if err != nil {
        return Stats{}, err
    }

    inserted++
    totalPrice += int(price)
    categorySet[row[2]] = struct{}{}
}

	if err := tx.Commit(); err != nil {
		return Stats{}, err
	}

	return Stats{
		TotalItems:      inserted,
		TotalCategories: len(categorySet),
		TotalPrice:      totalPrice,
	}, nil
}

// ===================== Query for GET =====================
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

	return result, rows.Err()
}
