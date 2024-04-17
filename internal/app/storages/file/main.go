// Package file нужен для работы с файлом, когда файл выбран хранилищем.
package file

import (
	"os"
	"path/filepath"
)

// File - структура для работы с файлом-хранилищем.
// Path - поле, которое хранит путь к файлу-хранилищу.
type File struct {
	Path string
}

// MakeFile - используется для создания файла-хранилища.
// Сначала создаются папки, которые нужны для полного пути к файлу.
// Если файл существует, то ничего не происходит.
// Если файл ещё не создан, то он создается.
func (fs *File) MakeFile() {
	configFile := fs.Path
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
