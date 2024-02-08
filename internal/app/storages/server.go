package storages

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type ServerStorage struct {
	st []Storage
	sa config.ServerConfig
	co []string
}

func (serv *ServerStorage) Init() {
	serv.sa.ParseFlags()
	if serv.GetConfigParameter("db") != "" {
		serv.MakeDB()
	} else if serv.GetConfigParameter("file") != "" {
		serv.MakeFile()
	}
}

func (serv *ServerStorage) MakeFile() {
	configFile := serv.GetConfigParameter("file")
	if _, err := os.Stat(filepath.Dir(configFile)); err != nil {
		os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
	}

	if _, err := os.Stat(configFile); err != nil {
		file, err1 := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err1 != nil {
			panic(err)
		}
		defer file.Close()

		file.Write([]byte("[]"))
	}
}

func (serv *ServerStorage) MakeDB() {
	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
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

func (serv *ServerStorage) GetConfigParameter(param string) string {
	switch param {
	case "address":
		return serv.sa.GetStartAddress()
	case "baseURL":
		return serv.sa.GetShortBaseURL()
	case "file":
		return serv.sa.GetFilePath()
	case "db":
		return serv.sa.GetDBAddress()
	default:
		return ""
	}
}

func (serv *ServerStorage) SetConfigParameter(value string, param string) {
	switch param {
	case "address":
		serv.sa.SetStartAddress(value)
	case "baseURL":
		serv.sa.SetShortBaseURL(value)
	case "file":
		serv.sa.SetFilePath(value)
	case "db":
		serv.sa.SetDBAddress(value)
	default:
		break
	}
}

func (serv *ServerStorage) MakeShortURL(longURL string) (string, bool, error) {
	value, _, err := serv.Get(longURL, "value", "")
	if err != nil {
		return "", false, err
	}
	if value != "" {
		return value, true, nil
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			result, _, err := serv.Get(shortKey, "key", "")
			if err != nil {
				return "", false, err
			}
			if result == "" {
				return shortKey, false, nil
			}
			shortKey = tools.GenerateShortKey()
		}
	}
}

// get

func (serv *ServerStorage) Get(value string, param string, userid string) (string, []bool, error) {
	if serv.GetConfigParameter("db") != "" {
		return serv.GetDBData(value, param, userid)
	}
	if serv.GetConfigParameter("file") != "" {
		return serv.GetFileData(value, param, userid)
	}
	return serv.GetStorageData(value, param, userid)
}

func (serv *ServerStorage) GetDBData(value string, param string, userid string) (string, []bool, error) {
	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
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
	} else {
		var stors []Storage
		for rows.Next() {
			var stor Storage
			if err = rows.Scan(&stor.key, &stor.value, &stor.userid, &stor.deleted); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return "", nil, nil
				} else {
					return "", nil, err
				}
			}
			stors = append(stors, stor)
		}
		if err = rows.Err(); err != nil {
			return "", nil, err
		}
		mystr := ""
		var mybool []bool
		for _, elem := range stors {
			mystr += elem.key + "*" + elem.value + "*"
			mybool = append(mybool, elem.deleted)
		}
		if mystr != "" {
			return mystr[:len(mystr)-1], mybool, nil
		}
		return "", nil, nil
	}
}

func (serv *ServerStorage) GetFileData(value string, param string, userid string) (string, []bool, error) {
	jsonData, err := tools.ReadFile(serv.GetConfigParameter("file"))
	if err != nil {
		return "", nil, err
	}

	if param == "key" {
		for _, elem := range jsonData {
			if elem.Key == value {
				return elem.Value, []bool{elem.Deleted}, nil
			}
		}
	} else if param == "value" {
		for _, elem := range jsonData {
			if elem.Value == value {
				return elem.Key, []bool{elem.Deleted}, nil
			}
		}
	} else if param == "all" {
		mystr := ""
		var mybool []bool
		for _, elem := range jsonData {
			if elem.UserID == userid {
				mystr += elem.Key + "*" + elem.Value + "*"
				mybool = append(mybool, elem.Deleted)
			}
		}
		if mystr != "" {
			return mystr[:len(mystr)-1], mybool, nil
		}
	} else {
		return "", nil, errors.New("unknown param")
	}
	return "", nil, nil
}

func (serv *ServerStorage) GetStorageData(value string, param string, userid string) (string, []bool, error) {
	if param == "key" {
		for _, elem := range serv.st {
			if elem.key == value {
				return elem.value, []bool{elem.deleted}, nil
			}
		}
	} else if param == "value" {
		for _, elem := range serv.st {
			if elem.value == value {
				return elem.key, []bool{elem.deleted}, nil
			}
		}
	} else if param == "all" {
		mystr := ""
		var mybool []bool
		for _, elem := range serv.st {
			if elem.userid == userid {
				mystr += elem.key + "*" + elem.value + "*"
				mybool = append(mybool, elem.deleted)
			}
		}
		if mystr != "" {
			return mystr[:len(mystr)-1], mybool, nil
		}
	} else {
		return "", nil, errors.New("unknown param")
	}
	return "", nil, nil
}

// set

func (serv *ServerStorage) Save(key string, value string, userid string) error {
	if serv.GetConfigParameter("db") != "" {
		if err := serv.SetDBData(key, value, userid); err != nil {
			return err
		}
	} else if serv.GetConfigParameter("file") != "" {
		if err := serv.SetFileData(key, value, userid); err != nil {
			return err
		}
	} else {
		serv.SetStorage(key, value, userid)
	}
	return nil
}

func (serv *ServerStorage) SetFileData(key string, value string, userid string) error {
	jsonData, err := tools.ReadFile(serv.GetConfigParameter("file"))
	if err != nil {
		return err
	}

	file, err := os.OpenFile(serv.GetConfigParameter("file"), os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData = append(jsonData, tools.JSONData{UUID: len(jsonData) + 1, Key: key, Value: value, UserID: userid})
	toFileData, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		return err
	}

	if _, err = file.Write(toFileData); err != nil {
		return err
	}
	return nil
}

func (serv *ServerStorage) SetDBData(key string, value string, userid string) error {
	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
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

func (serv *ServerStorage) SetStorage(key string, value string, userid string) {
	serv.st = append(serv.st, *NewStorage(key, value, userid))
}

// cookie

func (serv *ServerStorage) CountCookies() int {
	return len(serv.co)
}

func (serv *ServerStorage) SetCookie(value string) {
	serv.co = append(serv.co, value)
}

func (serv *ServerStorage) InCookies(value string) bool {
	for ind := range serv.co {
		if serv.co[ind] == value {
			return true
		}
	}
	return false
}

func (serv *ServerStorage) CheckCookie(r *http.Request) (string, bool, error) {
	cookie, err := r.Cookie("UserID")
	if err != nil || !serv.InCookies(cookie.Value) {
		gen, err1 := tools.GenerateCookie(serv.CountCookies() + 1)
		if err1 != nil {
			return "", false, err
		}
		serv.SetCookie(gen)
		if r.Method == http.MethodGet {
			return "", false, nil
		}
		return gen, false, nil
	}
	return cookie.Value, true, nil
}
