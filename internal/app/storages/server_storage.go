package storages

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type ServerStorage struct {
	sc config.ServerConfig
	cs CookieStorage
	SD ServerData
}

func (serv *ServerStorage) Init() {
	serv.ParseFlags()
	serv.InitStorage()
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
	serv.SD = ServerData{data: &db}
}

func (serv *ServerStorage) InitFile() {
	file := FileStorage{path: serv.GetServerConfig("file")}
	file.MakeFile()
	serv.SD = ServerData{data: &file}
}

func (serv *ServerStorage) InitMemory() {
	memory := MemoryStorage{arr: []MSItem{}}
	serv.SD = ServerData{data: &memory}
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

func (serv *ServerStorage) GetStorageData() GetterSetter {
	return serv.SD.data
}

func (serv *ServerStorage) MakeShortURL(longURL string) (string, bool, error) {
	value, _, err := serv.GetStorageData().GetData(longURL, "value", "")
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
