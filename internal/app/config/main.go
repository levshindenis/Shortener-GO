// Package config нужен для сбора флагов и переменных окружения.
package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"strconv"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// ServerConfig - структура для хранения флагов и переменных оркужения.
type ServerConfig struct {
	startAddress   string
	shortBaseURL   string
	filePath       string
	dbAddress      string
	https          bool
	configFilePath string
	trustedSubnet  string
	startAddressG  string
	gTls           bool
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

// GetHTTPS - Возвращает TLS
func (sa *ServerConfig) GetHTTPS() bool {
	return sa.https
}

// GetConfigFilePath - Возвращает путь config файла
func (sa *ServerConfig) GetConfigFilePath() string {
	return sa.configFilePath
}

// GetTrustedSubnet - возвращает разрешенные IP
func (sa *ServerConfig) GetTrustedSubnet() string {
	return sa.trustedSubnet
}

// GetStartAddressG - возвращает адрес запуска gRPC-сервера.
func (sa *ServerConfig) GetStartAddressG() string {
	return sa.startAddressG
}

// GetGTLS - Возвращает TLS gRPC
func (sa *ServerConfig) GetGTLS() bool {
	return sa.gTls
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

// SetHTTPS - устанавливает значение value для https.
func (sa *ServerConfig) SetHTTPS(value bool) {
	sa.https = value
}

// SetConfigFilePath - устанавливает значение value для configFilePath.
func (sa *ServerConfig) SetConfigFilePath(value string) {
	sa.configFilePath = value
}

// SetTrustedSubnet - устанавливает значение value для разрешенных IP
func (sa *ServerConfig) SetTrustedSubnet(value string) {
	sa.trustedSubnet = value
}

// SetStartAddressG - устанавливает значение value для startAddressG.
func (sa *ServerConfig) SetStartAddressG(value string) {
	sa.startAddressG = value
}

// SetGTLS - устанавливает значение value для gTLS.
func (sa *ServerConfig) SetGTLS(value bool) {
	sa.gTls = value
}

// ParseFlags - берет значения из флагов или переменных окружения и устанавливает значения в структуру ServerConfig.
func (sa *ServerConfig) ParseFlags() error {
	flag.StringVar(&sa.startAddress, "a", "localhost:8080", "address and port to run shortener")
	flag.StringVar(&sa.shortBaseURL, "b", "http://localhost:8080", "address and port for base short URL")
	flag.StringVar(&sa.filePath, "f", "/tmp/short-url-db.json", "storage file path")
	flag.StringVar(&sa.dbAddress, "d", "", "db address")
	flag.BoolVar(&sa.https, "s", false, "tls")
	flag.StringVar(&sa.configFilePath, "c", "", "config file path")
	flag.StringVar(&sa.trustedSubnet, "t", "", "IPs")
	flag.StringVar(&sa.startAddressG, "ga", ":3200", "address and port to run gRPC shortener")
	flag.BoolVar(&sa.gTls, "gs", false, "tls gRPC")

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

	if envHTTPS := os.Getenv("ENABLE_HTTPS"); envHTTPS != "" {
		sa.https, _ = strconv.ParseBool(envHTTPS)
	}

	if envConfigFile, in := os.LookupEnv("CONFIG"); in {
		sa.SetConfigFilePath(envConfigFile)
	}

	if envTrustedSubnet, in := os.LookupEnv("TRUSTED_SUBNET"); in {
		sa.SetTrustedSubnet(envTrustedSubnet)
	}

	if envStartAddressG, in := os.LookupEnv("SERVER_ADDRESS_GRPC"); in {
		sa.SetStartAddressG(envStartAddressG)
	}

	if envGTLS := os.Getenv("ENABLE_TLS_GRPC"); envGTLS != "" {
		sa.gTls, _ = strconv.ParseBool(envGTLS)
	}

	if sa.GetConfigFilePath() != "" {
		if err := sa.ReadConfigFile(); err != nil {
			return err
		}
	}
	return nil
}

// ReadConfigFile - устанавливает новые значения для переменных окружения из config файла
func (sa *ServerConfig) ReadConfigFile() error {
	var jsonData models.ConfigFileData

	file, err := os.OpenFile(sa.GetConfigFilePath(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	fromFileData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(fromFileData, &jsonData); err != nil {
		return err
	}

	if sa.GetStartAddress() == "" {
		sa.SetStartAddress(jsonData.ServerAddress)
	}

	if sa.GetShortBaseURL() == "" {
		sa.SetShortBaseURL(jsonData.BaseURL)
	}

	if sa.GetFilePath() == "" {
		sa.SetFilePath(jsonData.FileStoragePath)
	}

	if sa.GetDBAddress() == "" {
		sa.SetDBAddress(jsonData.DatabaseDSN)
	}

	if !sa.GetHTTPS() {
		sa.SetHTTPS(jsonData.EnableHTTPS)
	}

	if sa.GetTrustedSubnet() == "" {
		sa.SetTrustedSubnet(jsonData.TrustedSubnet)
	}

	if sa.GetStartAddressG() == "" {
		sa.SetStartAddressG(jsonData.ServerAddressG)
	}

	if !sa.GetGTLS() {
		sa.SetGTLS(jsonData.GTls)
	}

	return nil
}
