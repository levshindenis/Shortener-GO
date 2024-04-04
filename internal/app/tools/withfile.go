package tools

import (
	"encoding/json"
	"io"
	"os"

	"github.com/levshindenis/sprint1/internal/app/models"
)

func ReadFile(filepath string) ([]models.JSONData, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fromFileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonData []models.JSONData
	if err = json.Unmarshal(fromFileData, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}
