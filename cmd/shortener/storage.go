package main

type Storage map[string]string

func (storage *Storage) EmptyStorage() {
	*storage = make(map[string]string)
}

// ValueIn проверяет наличие значения в map
func (storage *Storage) ValueIn(s string) string {
	for key, value := range *storage {
		if value == s {
			return key
		}
	}
	return ""
}
