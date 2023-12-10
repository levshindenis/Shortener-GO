package config

import "flag"

type ServerAddress struct {
	startAddress    string
	shortURLAddress string
}

func (sa *ServerAddress) GetStartAddress() string {
	return sa.startAddress
}

func (sa *ServerAddress) GetShortURLAddress() string {
	return sa.shortURLAddress
}

func (sa *ServerAddress) SetStartAddress(value string) {
	sa.startAddress = value
}

func (sa *ServerAddress) SetShortURLAddress(value string) {
	sa.shortURLAddress = value
}

func ParseFlags(sa *ServerAddress) {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&sa.shortURLAddress, "b", "localhost:8080", "address and port for base short URL")

	flag.Parse()
}
