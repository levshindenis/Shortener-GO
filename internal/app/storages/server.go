package storages

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type GetterSetter interface {
	SetData(key string, value string, userid string) error
	GetData(value string, param string, userid string) (string, []bool, error)
}

type ServerData struct {
	data GetterSetter
}

type ServerStorage struct {
	sc config.ServerConfig
	cs CookieStorage
	sd ServerData
}

func (serv *ServerStorage) Init() {
	serv.ParseFlags()
	serv.InitStorage()
}

func (serv *ServerStorage) InitStorage() {
	if serv.GetSC("db") != "" {
		serv.InitDB()
	} else if serv.GetSC("file") != "" {
		serv.InitFile()
	} else {
		serv.InitMemory()
	}
}

func (serv *ServerStorage) ParseFlags() {
	serv.sc.ParseFlags()
}

func (serv *ServerStorage) InitDB() {
	db := DBStorage{address: serv.GetSC("db")}
	db.MakeDB()
	serv.sd = ServerData{data: &db}
}

func (serv *ServerStorage) InitFile() {
	file := FileStorage{path: serv.GetSC("file")}
	file.MakeFile()
	serv.sd = ServerData{data: &file}
}

func (serv *ServerStorage) InitMemory() {
	memory := MemoryStorage{arr: []MSItem{}}
	serv.sd = ServerData{data: &memory}
}

func (serv *ServerStorage) GetSC(param string) string {
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

func (serv *ServerStorage) SetSC(value string, param string) {
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

func (serv *ServerStorage) GetCS() CookieStorage {
	return serv.cs
}

func (serv *ServerStorage) GetSD() GetterSetter {
	return serv.sd.data
}

func (serv *ServerStorage) MakeShortURL(longURL string) (string, bool, error) {
	value, _, err := serv.GetSD().GetData(longURL, "value", "")
	if err != nil {
		return "", false, err
	}
	if value != "" {
		return value, true, nil
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			result, _, err := serv.GetSD().GetData(shortKey, "key", "")
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
