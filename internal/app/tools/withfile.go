package tools

import (
	"encoding/json"
	"io"
	"os"
)

type JSONData struct {
	UUID    int    `json:"uuid"`
	Key     string `json:"short_url"`
	Value   string `json:"original_url"`
	UserID  string `json:"user_id"`
	Deleted bool   `json:"deleted"`
}

func ReadFile(filepath string) ([]JSONData, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fromFileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonData []JSONData
	if err = json.Unmarshal(fromFileData, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}
