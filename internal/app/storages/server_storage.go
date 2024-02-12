package storages

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type ServerStorage struct {
	sc config.ServerConfig
	cs CookieStorage
	sd ServerData
	ch chan DeleteValue
}

func (serv *ServerStorage) Init() {
	serv.ParseFlags()
	serv.InitStorage()
	serv.ch = make(chan DeleteValue, 1024)
	go serv.DeleteItems()
}

func (serv *ServerStorage) InitStorage() {
	if serv.GetServerConfig("db") != "" {
		serv.InitDB()
	} else if serv.GetServerConfig("file") != "" {
		serv.InitFile()
	} else {
		serv.InitMemory()
	}
}

func (serv *ServerStorage) ParseFlags() {
	serv.sc.ParseFlags()
}

func (serv *ServerStorage) InitDB() {
	db := DBStorage{address: serv.GetServerConfig("db")}
	db.MakeDB()
	serv.sd = ServerData{data: &db}
}

func (serv *ServerStorage) InitFile() {
	file := FileStorage{path: serv.GetServerConfig("file")}
	file.MakeFile()
	serv.sd = ServerData{data: &file}
}

func (serv *ServerStorage) InitMemory() {
	memory := MemoryStorage{arr: []MSItem{}}
	serv.sd = ServerData{data: &memory}
}

func (serv *ServerStorage) GetServerConfig(param string) string {
	switch param {
	case "address":
		return serv.sc.GetStartAddress()
	case "baseURL":
		return serv.sc.GetShortBaseURL()
	case "file":
		return serv.sc.GetFilePath()
	case "db":
		return serv.sc.GetDBAddress()
	default:
		return ""
	}
}

func (serv *ServerStorage) SetServerConfig(value string, param string) {
	switch param {
	case "address":
		serv.sc.SetStartAddress(value)
	case "baseURL":
		serv.sc.SetShortBaseURL(value)
	case "file":
		serv.sc.SetFilePath(value)
	case "db":
		serv.sc.SetDBAddress(value)
	default:
		break
	}
}

func (serv *ServerStorage) GetCookieStorage() *CookieStorage {
	return &serv.cs
}

func (serv *ServerStorage) GetStorageData() BaseFuncs {
	return serv.sd.data
}

func (serv *ServerStorage) SetChan(delValue DeleteValue) {
	serv.ch <- delValue
}

func (serv *ServerStorage) MakeShortURL(longURL string) (string, bool, error) {
	value, _, err := serv.GetStorageData().GetData(longURL, "Value", "")
	if err != nil {
		return "", false, err
	}
	if value != "" {
		return value, true, nil
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			result, _, err := serv.GetStorageData().GetData(shortKey, "key", "")
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

func (serv *ServerStorage) DeleteItems() {
	ticker := time.NewTicker(2 * time.Second)

	var values []DeleteValue

	for {
		select {
		case value := <-serv.ch:
			// добавим сообщение в слайс для последующего сохранения
			values = append(values, value)
		case <-ticker.C:
			// подождём, пока придёт хотя бы одно сообщение
			if len(values) == 0 {
				continue
			}
			// сохраним все пришедшие сообщения одновременно
			err := serv.GetStorageData().DeleteData(values)
			if err != nil {
				panic(err)
			}
			// сотрём успешно отосланные сообщения
			values = nil
		}
	}
}
