package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// GetData - нужна для получение каких-либо данных из БД.
// Открывается БД.
// Если param == key, то будет возвращен длинный URL и параметр deleted.
// Если param == value, то будет возвращен короткий URL и параметр deleted.
// Если param == all, то будут возвращены все записи по полученному UserID.
func (dbs *Database) GetData(value string, param string, userid string) (string, []bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		row  *sql.Row
		rows *sql.Rows
		err  error
	)

	if param == "key" {
		row = dbs.DB.QueryRowContext(ctx, `SELECT long_url, deleted FROM shortener WHERE short_url = $1`, value)
	}
	if param == "value" {
		row = dbs.DB.QueryRowContext(ctx, `SELECT short_url, deleted FROM shortener WHERE long_url = $1`, value)
	}
	if param == "all" {
		rows, err = dbs.DB.QueryContext(ctx, `SELECT * FROM shortener WHERE user_id = $1`, userid)
		if err != nil {
			return "", nil, nil
		}
		defer rows.Close()
	}

	if row != nil {
		var (
			result  string
			deleted bool
		)

		err = row.Scan(&result, &deleted)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", nil, nil
			}
			return "", nil, err
		}
		return result, []bool{deleted}, nil
	}

	var (
		item   models.MSItem
		myBool []bool
		myStr  strings.Builder
	)

	for rows.Next() {
		if err = rows.Scan(&item.Key, &item.Value, &item.UserID, &item.Deleted); err != nil {
			return "", nil, err
		}
		myStr.WriteString(item.Key + "*" + item.Value + "*")
		myBool = append(myBool, item.Deleted)
	}
	if err = rows.Err(); err != nil {
		return "", nil, err
	}

	return myStr.String(), myBool, nil
}
