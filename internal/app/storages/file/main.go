package file

import (
	"os"
	"path/filepath"
)

type File struct {
	Path string
}

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
