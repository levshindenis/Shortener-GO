package file

import (
	"errors"

	"github.com/levshindenis/sprint1/internal/app/tools"
)

// GetData - нужна для получение каких-либо данных из файла-хранилища.
// Берутся данные из файла-хранилища.
// Если param == key, то будет возвращен длинный URL и параметр deleted.
// Если param == value, то будет возвращен короткий URL и параметр deleted.
// Если param == all, то будут возвращены все записи по полученному UserID.
func (fs *File) GetData(value string, param string, userid string) (string, []bool, error) {
	jsonData, err := tools.ReadFile(fs.Path)
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
	if param == "value" {
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
