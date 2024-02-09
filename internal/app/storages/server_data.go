package storages

type GetterSetter interface {
	SetData(key string, value string, userid string) error
	GetData(value string, param string, userid string) (string, []bool, error)
}

type ServerData struct {
	data GetterSetter
}
