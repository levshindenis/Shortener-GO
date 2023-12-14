package config

import (
	"flag"
	"os"
)

type ServerAddress struct {
	startAddress string
	shortBaseURL string
}

func (sa *ServerAddress) GetStartAddress() string {
	return sa.startAddress
}

func (sa *ServerAddress) GetShortBaseURL() string {
	return sa.shortBaseURL
}

func (sa *ServerAddress) SetStartAddress(value string) {
	sa.startAddress = value
}

func (sa *ServerAddress) SetShortBaseURL(value string) {
	sa.shortBaseURL = value
}

func (sa *ServerAddress) ParseFlags() {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run shortener")
	flag.StringVar(&sa.shortBaseURL, "b", "http://localhost:8080", "address and port for base short URL")

	flag.Parse()

	if envStartAddress := os.Getenv("SERVER_ADDRESS"); envStartAddress != "" {
		sa.SetStartAddress(envStartAddress)
	}

	if envShortBaseURL := os.Getenv("BASE_URL"); envShortBaseURL != "" {
		sa.SetShortBaseURL(envShortBaseURL)
	}

}
