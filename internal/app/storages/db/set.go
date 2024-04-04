package db

import (
	"context"
	"database/sql"
	"time"
)

func (dbs *Database) SetData(key string, value string, userid string) error {
	db, err := sql.Open("pgx", dbs.Address)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx,
		`INSERT INTO shortener (short_url, long_url, user_id, deleted) values ($1, $2, $3, $4)`,
		key, value, userid, false)
	if err != nil {
		return err
	}

	return nil
}
