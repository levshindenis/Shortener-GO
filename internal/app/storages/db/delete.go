package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// DeleteData нужа для "удаления" переданных сокращенных URL из БД.
// Открывается БД.
// В цикле берется каждый короткий URL и по нему сравнивается UserID из поступивших данных и значение из БД.
// Если значения не совпадают, то удаление не происходит.
// Если значения совпали, то меняется значение "deleted" на true.
func (dbs *Database) DeleteData(delValues []models.DeleteValue) error {
	db, err := sql.Open("pgx", dbs.Address)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, elem := range delValues {
		row := db.QueryRowContext(ctx, `SELECT user_id FROM shortener WHERE short_url = $1`, elem.Value)
		var result string
		err = row.Scan(&result)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			} else {
				return err
			}
		}
		if result != elem.Userid {
			continue
		}
		_, err = tx.ExecContext(ctx,
			`UPDATE shortener SET deleted = $1 WHERE short_url = $2`, true, elem.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
