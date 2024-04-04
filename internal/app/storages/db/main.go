package db

import (
	"context"
	"database/sql"
	"time"
)

type Database struct {
	Address string
}

func (dbs *Database) MakeDB() {
	db, err := sql.Open("pgx", dbs.Address)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS shortener(short_url text, long_url text, user_id text, deleted boolean)`)
	if err != nil {
		panic(err)
	}
}
