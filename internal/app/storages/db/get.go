package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// GetData - нужна для получение каких-либо данных из БД.
// Открывается БД.
// Если param == key, то будет возвращен длинный URL и параметр deleted.
// Если param == value, то будет возвращен короткий URL и параметр deleted.
// Если param == all, то будут возвращены все записи по полученному UserID.
func (dbs *Database) GetData(value string, param string, userid string) (string, []bool, error) {
	db, err := sql.Open("pgx", dbs.Address)
	if err != nil {
		return "", nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var row *sql.Row
	var rows *sql.Rows

	if param == "key" {
		row = db.QueryRowContext(ctx, `SELECT long_url, deleted FROM shortener WHERE short_url = $1`, value)
	} else if param == "value" {
		row = db.QueryRowContext(ctx, `SELECT short_url, deleted FROM shortener WHERE long_url = $1`, value)
	} else if param == "all" {
		rows, err = db.QueryContext(ctx, `SELECT * FROM shortener WHERE user_id = $1`, userid)
		if err != nil {
			return "", nil, nil
		}
		defer rows.Close()
	} else {
		return "", nil, errors.New("unknown param")
	}

	if row != nil {
		var result string
		var deleted bool
		err = row.Scan(&result, &deleted)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", nil, nil
			} else {
				return "", nil, err
			}
		}
		return result, []bool{deleted}, nil
	}

	var items []models.MSItem
	for rows.Next() {
		var item models.MSItem
		if err = rows.Scan(&item.Key, &item.Value, &item.UserID, &item.Deleted); err != nil {
			return "", nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return "", nil, err
	}
	mystr := ""
	var mybool []bool
	for _, elem := range items {
		mystr += elem.Key + "*" + elem.Value + "*"
		mybool = append(mybool, elem.Deleted)
	}
	if mystr != "" {
		return mystr[:len(mystr)-1], mybool, nil
	}
	return "", nil, nil
}
