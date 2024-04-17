package db

import (
	"context"
	"time"
)

// SetData - нужна для записи значений в БД.
func (dbs *Database) SetData(key string, value string, userid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := dbs.DB.ExecContext(ctx,
		`INSERT INTO shortener (short_url, long_url, user_id, deleted) values ($1, $2, $3, $4)`,
		key, value, userid, false)
	if err != nil {
		return err
	}

	return nil
}
