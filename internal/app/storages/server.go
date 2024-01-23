package storages

import (
	"encoding/json"
	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
	"io"
	"os"
	"path/filepath"
)

type ServerStorage struct {
	st Storage
	sa config.ServerConfig
}

func (serv *ServerStorage) Init() {
	serv.st.EmptyStorage()
	serv.sa.ParseFlags()
	if serv.GetFilePath() != "" {
		serv.MakeDir()
		serv.GetFileData()
	}
}

func (serv *ServerStorage) GetStorage() Storage {
	return serv.st
}

func (serv *ServerStorage) SetStorage(key string, value string) {
	serv.GetStorage()[key] = value
}

func (serv *ServerStorage) ValueInStorage(value string) (string, bool) {
	return serv.st.ValueIn(value)
}

func (serv *ServerStorage) InitStorage() {
	serv.st.EmptyStorage()
}

func (serv *ServerStorage) GetStartSA() string {
	return serv.sa.GetStartAddress()
}

func (serv *ServerStorage) GetBaseSA() string {
	return serv.sa.GetShortBaseURL()
}

func (serv *ServerStorage) SetBaseSA(value string) {
	serv.sa.SetShortBaseURL(value)
}

func (serv *ServerStorage) GetFilePath() string {
	return serv.sa.GetFilePath()
}

func (serv *ServerStorage) SetFilePath(value string) {
	serv.sa.SetFilePath(value)
}

func (serv *ServerStorage) MakeDir() {
	serv.SetFilePath(serv.GetFilePath()[1:])
	if _, err := os.Stat(serv.GetFilePath()); err != nil {
		os.MkdirAll(filepath.Dir(serv.GetFilePath()), os.ModePerm)
	}
}

func (serv *ServerStorage) GetFileData() {
	type JSONData struct {
		UUID  int    `json:"uuid"`
		Key   string `json:"short_url"`
		Value string `json:"original_url"`
	}

	file, err := os.OpenFile(serv.GetFilePath(), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(serv.GetFilePath())
	if err != nil {
		panic(err)
	}

	if fileInfo.Size() == 0 {
		file.Write([]byte("[]"))
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var jsonData []JSONData
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		panic(err)
	}

	for _, elem := range jsonData {
		serv.SetStorage(elem.Key, elem.Value)
	}
}

func (serv *ServerStorage) GetAddress(str string) (string, error) {
	addr := serv.GetBaseSA() + "/"
	if value, ok := serv.ValueInStorage(str); ok {
		return addr + value, nil
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			if _, in := serv.GetStorage()[shortKey]; !in {
				serv.SetStorage(shortKey, str)
				break
			}
			shortKey = tools.GenerateShortKey()
		}
		if err := serv.Save(shortKey, str); err != nil {
			return "", err
		}
		return addr + shortKey, nil
	}
}

func (serv *ServerStorage) Save(key string, value string) error {
	if serv.GetFilePath() == "" {
		return nil
	}

	type JSONData struct {
		UUID  int    `json:"uuid"`
		Key   string `json:"short_url"`
		Value string `json:"original_url"`
	}

	file, err := os.OpenFile(serv.GetFilePath(), os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	fromFileData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var jsonData []JSONData
	if err := json.Unmarshal(fromFileData, &jsonData); err != nil {
		return err
	}

	jsonData = append(jsonData, JSONData{UUID: len(serv.GetStorage()), Key: key, Value: value})
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	toFileData, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		return err
	}
	if _, err := file.Write(toFileData); err != nil {
		return err
	}

	return nil
}
