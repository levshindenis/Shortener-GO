package storages

import (
	"errors"
)

type MemoryStorage struct {
	arr []MSItem
}

type MSItem struct {
	key     string
	value   string
	userid  string
	deleted bool
}

func (ms *MemoryStorage) GetArr() []MSItem {
	return ms.arr
}

func (ms *MemoryStorage) SetData(key string, value string, userid string) error {
	ms.arr = append(ms.GetArr(), MSItem{key: key, value: value, userid: userid, deleted: false})
	return nil
}

func (ms *MemoryStorage) GetData(value string, param string, userid string) (string, []bool, error) {
	if param == "key" {
		for _, elem := range ms.GetArr() {
			if elem.key == value {
				return elem.value, []bool{elem.deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "Value" {
		for _, elem := range ms.GetArr() {
			if elem.value == value {
				return elem.key, []bool{elem.deleted}, nil
			}
		}
		return "", nil, nil
	}
	if param == "all" {
		mystr := ""
		var mybool []bool
		for _, elem := range ms.GetArr() {
			if elem.userid == userid {
				mystr += elem.key + "*" + elem.value + "*"
				mybool = append(mybool, elem.deleted)
			}
		}
		if mystr != "" {
			return mystr[:len(mystr)-1], mybool, nil
		}
		return "", nil, nil
	}
	return "", nil, errors.New("unknown param")
}

func (ms *MemoryStorage) DeleteData(delValues []DeleteValue) error {
	for _, elem := range delValues {
		for _, msi := range ms.GetArr() {
			if msi.key == elem.Value && msi.userid == elem.Userid {
				msi.deleted = true
			}
		}
	}
	return nil
}
