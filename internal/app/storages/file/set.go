package file

import (
	"encoding/json"
	"os"

	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

// SetData - нужна для записи значений в файл-хранилище.
// Сначала берутся данные из файла.
// После этого файл отрывается с параметрами очистки (файл становится пустым).
// В jsonData записываются полученные значения.
// В файл в формате JSON записываются данные из jsobData.
func (fs *File) SetData(key string, value string, userid string) error {
	jsonData, err := tools.ReadFile(fs.Path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fs.Path, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData = append(jsonData, models.JSONData{UUID: len(jsonData) + 1, Key: key, Value: value, UserID: userid})
	toFileData, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		return err
	}

	if _, err = file.Write(toFileData); err != nil {
		return err
	}
	return nil
}
