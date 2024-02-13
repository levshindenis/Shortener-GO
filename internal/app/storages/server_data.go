package storages

import "context"

type BaseFuncs interface {
	SetData(key string, value string, userid string) error
	GetData(value string, param string, userid string) (string, []bool, error)
	DeleteData(delValues []DeleteValue) error
}

type ServerData struct {
	data BaseFuncs
}

type DeleteValue struct {
	Value  string
	Userid string
}

type ChanData struct {
	ch     chan DeleteValue
	ctx    context.Context
	cancel context.CancelFunc
}
