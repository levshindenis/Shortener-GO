package file

import (
	"encoding/json"
	"os"

	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

// DeleteData нужа для "удаления" переданных сокращенных URL из файла.
// Сначала берутся данные из файла.
// После этого открывается файл с параметрами очистки (файл становится пустым).
// В цикле берется каждый короткий URL с UserID и сравнивается с данными в файле.
// Если данные совпадают, то deleted = true.
// После этого измененные данные в формате JSON записываются в файл.
func (fs *File) DeleteData(delValues []models.DeleteValue) error {
	jsonData, err := tools.ReadFile(fs.Path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fs.Path, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, elem := range delValues {
		for ind, jd := range jsonData {
			if elem.Value == jd.Key && elem.Userid == jd.UserID {
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
