package db

import (
	"context"
	"strconv"
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
		price, _ := strconv.Atoi(row[3])

		_, err := stmt.ExecContext(
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
		totalPrice += price
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
