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
	var result string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := dbs.DB.Begin()
	if err != nil {
		return err
	}

	for ind := range delValues {
		row := dbs.DB.QueryRowContext(ctx, `SELECT user_id FROM shortener WHERE short_url = $1`, delValues[ind].Value)
		err = row.Scan(&result)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}
		if result != delValues[ind].Userid {
			continue
		}
		_, err = tx.ExecContext(ctx,
			`UPDATE shortener SET deleted = $1 WHERE short_url = $2`, true, delValues[ind].Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
