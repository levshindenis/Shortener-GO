package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	startAddress string
	shortBaseURL string
	filePath     string
	dbAddress    string
}

func (sa *ServerConfig) GetStartAddress() string {
	return sa.startAddress
}

func (sa *ServerConfig) GetShortBaseURL() string {
	return sa.shortBaseURL
}

func (sa *ServerConfig) GetFilePath() string {
	return sa.filePath
}

func (sa *ServerConfig) GetDBAddress() string {
	return sa.dbAddress
}

func (sa *ServerConfig) SetStartAddress(value string) {
	sa.startAddress = value
}

func (sa *ServerConfig) SetShortBaseURL(value string) {
	sa.shortBaseURL = value
}

func (sa *ServerConfig) SetFilePath(value string) {
	sa.filePath = value
}

func (sa *ServerConfig) SetDBAddress(value string) {
	sa.dbAddress = value
}

func (sa *ServerConfig) ParseFlags() {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run shortener")
	flag.StringVar(&sa.shortBaseURL, "b", "http://localhost:8080", "address and port for base short URL")
	flag.StringVar(&sa.filePath, "f", "", "storage file path")
	flag.StringVar(&sa.dbAddress, "d", "", "db address")

	flag.Parse()

	if envStartAddress, in := os.LookupEnv("SERVER_ADDRESS"); in {
		sa.SetStartAddress(envStartAddress)
	}

	if envShortBaseURL, in := os.LookupEnv("BASE_URL"); in {
		sa.SetShortBaseURL(envShortBaseURL)
	}

	if envFilePath, in := os.LookupEnv("FILE_STORAGE_PATH"); in {
		sa.SetFilePath(envFilePath)
	}

	if envDBAddress, in := os.LookupEnv("DATABASE_DSN"); in {
		sa.SetDBAddress(envDBAddress)
	}
}
