// Package config нужен для сбора флагов и переменных окружения.
package config

import (
	"flag"
	"os"
)

// ServerConfig - структура для хранения флагов и переменных оркужения.
type ServerConfig struct {
	startAddress string
	shortBaseURL string
	filePath     string
	dbAddress    string
}

// GetStartAddress - возвращает адрес запуска HTTP-сервера.
func (sa *ServerConfig) GetStartAddress() string {
	return sa.startAddress
}

// GetShortBaseURL - возвращает базовый адрес результирующего сокращённого URL.
func (sa *ServerConfig) GetShortBaseURL() string {
	return sa.shortBaseURL
}

// GetFilePath - возвращает путь с именем файла, где хранятся данные.
func (sa *ServerConfig) GetFilePath() string {
	return sa.filePath
}

// GetDBAddress - возращает адрес БД.
func (sa *ServerConfig) GetDBAddress() string {
	return sa.dbAddress
}

// SetStartAddress - устанавливает значение value для startAddress.
func (sa *ServerConfig) SetStartAddress(value string) {
	sa.startAddress = value
}

// SetShortBaseURL - устанавливает значение value для shortBaseURL.
func (sa *ServerConfig) SetShortBaseURL(value string) {
	sa.shortBaseURL = value
}

// SetFilePath - устанавливает значение value для filePath.
func (sa *ServerConfig) SetFilePath(value string) {
	sa.filePath = value
}

// SetDBAddress - устанавливает значение value для dbAddress.
func (sa *ServerConfig) SetDBAddress(value string) {
	sa.dbAddress = value
}

// ParseFlags - берет значения из флагов или переменных окружения и устанавливает значения в структуру ServerConfig.
func (sa *ServerConfig) ParseFlags() {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run shortener")
	flag.StringVar(&sa.shortBaseURL, "b", "http://localhost:8080", "address and port for base short URL")
	flag.StringVar(&sa.filePath, "f", "/tmp/short-url-db.json", "storage file path")
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
