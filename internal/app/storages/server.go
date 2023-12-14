package storages

import "github.com/levshindenis/sprint1/internal/app/config"

type ServerStorage struct {
	st Storage
	sa config.ServerAddress
}

func (serv *ServerStorage) Init() {
	serv.st.EmptyStorage()
	serv.sa.ParseFlags()
}

func (serv *ServerStorage) InitStorage() {
	serv.st.EmptyStorage()
}

func (serv *ServerStorage) GetStartSA() string {
	return serv.sa.GetStartAddress()
}

func (serv *ServerStorage) GetBaseSA() string {
	return serv.sa.GetShortBaseURL()
}

func (serv *ServerStorage) GetStorage() Storage {
	return serv.st
}

func (serv *ServerStorage) SetStorage(key string, value string) {
	serv.st[key] = value
}

func (serv *ServerStorage) SetBaseSA(value string) {
	serv.sa.SetShortBaseURL(value)
}
