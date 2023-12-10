package config

import "flag"

type ServerAddress struct {
	startAddress    string
	shortUrlAddress string
}

func (sa *ServerAddress) GetStartAddress() string {
	return sa.startAddress
}

func (sa *ServerAddress) GetShortUrlAddress() string {
	return sa.shortUrlAddress
}

func (sa *ServerAddress) SetStartAddress(value string) {
	sa.startAddress = value
}

func (sa *ServerAddress) SetShortUrlAddress(value string) {
	sa.shortUrlAddress = value
}

func ParseFlags(sa *ServerAddress) {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&sa.shortUrlAddress, "b", "localhost:8080", "address and port for base short URL")

	flag.Parse()
}
