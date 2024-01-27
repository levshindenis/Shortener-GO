package storages

import "errors"

type Storage map[string]string

func (storage *Storage) EmptyStorage() {
	*storage = make(map[string]string)
}

func (storage *Storage) GetStorageData(value string, param string) (string, error) {
	if param == "key" {
		return (*storage)[value], nil
	} else if param == "value" {
		for k, v := range *storage {
			if v == value {
				return k, nil
			}
		}
		return "", nil
	} else {
		return "", errors.New("unknown param")
	}
}

func (storage *Storage) SetStorage(key string, value string) {
	(*storage)[key] = value
}
