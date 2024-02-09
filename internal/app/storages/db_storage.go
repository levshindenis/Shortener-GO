package storages

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type DBStorage struct {
	address string
}

func (dbs *DBStorage) GetAddress() string {
	return dbs.address
}

func (dbs *DBStorage) GetData(value string, param string, userid string) (string, []bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
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

	var items []MSItem
	for rows.Next() {
		var item MSItem
		if err = rows.Scan(&item.key, &item.value, &item.userid, &item.deleted); err != nil {
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
		mystr += elem.key + "*" + elem.value + "*"
		mybool = append(mybool, elem.deleted)
	}
	if mystr != "" {
		return mystr[:len(mystr)-1], mybool, nil
	}
	return "", nil, nil
}

func (dbs *DBStorage) SetData(key string, value string, userid string) error {
	db, err := sql.Open("pgx", dbs.GetAddress())
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

func (dbs *DBStorage) MakeDB() {
	db, err := sql.Open("pgx", dbs.GetAddress())
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
