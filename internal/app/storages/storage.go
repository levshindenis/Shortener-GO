package storages

type Storage struct {
	key     string
	value   string
	userid  string
	deleted bool
}

func NewStorage(key string, value string, userid string) *Storage {
	return &Storage{key: key, value: value, userid: userid, deleted: false}
}
