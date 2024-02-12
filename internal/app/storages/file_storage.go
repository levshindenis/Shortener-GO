package storages

import (
	"encoding/json"
	"errors"
	"github.com/levshindenis/sprint1/internal/app/tools"
	"os"
	"path/filepath"
)

type FileStorage struct {
	path string
}

func (fs *FileStorage) GetPath() string {
	return fs.path
}

func (fs *FileStorage) GetData(value string, param string, userid string) (string, []bool, error) {
	jsonData, err := tools.ReadFile(fs.GetPath())
	if err != nil {
		return "", nil, err
	}

	if param == "key" {
		for _, elem := range jsonData {
			if elem.Key == value {
				return elem.Value, []bool{elem.Deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "Value" {
		for _, elem := range jsonData {
			if elem.Value == value {
				return elem.Key, []bool{elem.Deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "all" {
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
		return "", nil, nil
	}
	return "", nil, errors.New("unknown param")
}

func (fs *FileStorage) SetData(key string, value string, userid string) error {
	jsonData, err := tools.ReadFile(fs.GetPath())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fs.GetPath(), os.O_TRUNC|os.O_WRONLY, os.ModePerm)
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

func (fs *FileStorage) DeleteData(delValues []DeleteValue) error {
	jsonData, err := tools.ReadFile(fs.GetPath())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fs.GetPath(), os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, elem := range delValues {
		for ind, jd := range jsonData {
			if elem.Value == jd.Value && elem.Userid == jd.UserID {
				jsonData[ind].Deleted = true
			}
		}
	}

	toFileData, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		return err
	}

	if _, err = file.Write(toFileData); err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) MakeFile() {
	configFile := fs.GetPath()
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
