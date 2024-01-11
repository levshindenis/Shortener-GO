package storages

type Storage map[string]string

func (storage *Storage) EmptyStorage() {
	*storage = make(map[string]string)
}

func (storage *Storage) ValueIn(str string) (string, bool) {
	for key, value := range *storage {
		if value == str {
			return key, true
		}
	}
	return "", false
}
