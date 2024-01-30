package storages

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type ServerStorage struct {
	st Storage
	sa config.ServerConfig
}

func (serv *ServerStorage) Init() {
	serv.sa.ParseFlags()
	if serv.GetConfigParameter("db") != "" {
		serv.MakeDB()
	} else if serv.GetConfigParameter("file") != "" {
		serv.MakeFile()
	} else {
		serv.st.EmptyStorage()
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

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS shortener(short_url text, long_url text)`)
	if err != nil {
		panic(err)
	}
}

func (serv *ServerStorage) InitStorage() {
	serv.st.EmptyStorage()
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
	value, err := serv.Get(longURL, "value")
	if err != nil {
		return "", false, err
	}
	if value != "" {
		return value, true, nil
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			result, err := serv.Get(shortKey, "key")
			if err != nil {
				fmt.Println("2")
				return "", false, err
			}
			if result == "" {
				return shortKey, false, nil
			}
			shortKey = tools.GenerateShortKey()
		}
	}
}

//

func (serv *ServerStorage) Get(value string, param string) (string, error) {
	if serv.GetConfigParameter("db") != "" {
		return serv.GetDBData(value, param)
	}
	if serv.GetConfigParameter("file") != "" {
		return serv.GetFileData(value, param)
	}
	return serv.GetStorageData(value, param)
}

func (serv *ServerStorage) GetDBData(value string, param string) (string, error) {
	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
	if err != nil {
		return "", err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbAnswer *sql.Row

	if param == "value" {
		dbAnswer = db.QueryRowContext(ctx, `SELECT short_url FROM shortener WHERE long_url = $1`, value)
	} else if param == "key" {
		dbAnswer = db.QueryRowContext(ctx, `SELECT long_url FROM shortener WHERE short_url = $1`, value)
	} else {
		return "", errors.New("unknown param")
	}

	var result string
	err = dbAnswer.Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		} else {
			return "", err
		}
	}

	return result, nil
}

func (serv *ServerStorage) GetFileData(value string, param string) (string, error) {
	jsonData, err := tools.ReadFile(serv.GetConfigParameter("file"))
	if err != nil {
		return "", err
	}

	if param == "value" {
		for _, elem := range jsonData {
			if elem.Value == value {
				return elem.Key, nil
			}
		}
	}
	if param == "key" {
		for _, elem := range jsonData {
			if elem.Key == value {
				return elem.Value, nil
			}
		}
	}
	return "", errors.New("unknown param")
}

func (serv *ServerStorage) GetStorageData(value string, param string) (string, error) {
	return serv.st.GetStorageData(value, param)
}

//

func (serv *ServerStorage) Save(key string, value string) error {
	if serv.GetConfigParameter("db") != "" {
		if err := serv.SetDBData(key, value); err != nil {
			return err
		}
	}
	if serv.GetConfigParameter("file") != "" {
		if err := serv.SetFileData(key, value); err != nil {
			return err
		}
	}
	serv.SetStorage(key, value)
	return nil
}

func (serv *ServerStorage) SetFileData(key string, value string) error {
	jsonData, err := tools.ReadFile(serv.GetConfigParameter("file"))
	if err != nil {
		return err
	}

	file, err := os.OpenFile(serv.GetConfigParameter("file"), os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData = append(jsonData, tools.JSONData{UUID: len(jsonData) + 1, Key: key, Value: value})
	toFileData, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		return err
	}

	if _, err = file.Write(toFileData); err != nil {
		return err
	}
	return nil
}

func (serv *ServerStorage) SetDBData(key string, value string) error {
	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `INSERT INTO shortener (short_url, long_url) values ($1, $2)`, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (serv *ServerStorage) SetStorage(key string, value string) {
	serv.st.SetStorage(key, value)
}
