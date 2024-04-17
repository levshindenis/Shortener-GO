// Package db нужен для работы с БД, когда БД выбрана хранилищем.
package db

import (
	"context"
	"database/sql"
	"time"
)

// Database - структура для работы с БД.
// Address - поле, которое хранит адрес подключения к БД.
type Database struct {
	DB *sql.DB
}

// MakeDB - создает таблицу "shortener" по адресу.
func (dbs *Database) MakeDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := dbs.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS shortener(short_url text, long_url text, user_id text, deleted boolean)`)
	if err != nil {
		panic(err)
	}
}
