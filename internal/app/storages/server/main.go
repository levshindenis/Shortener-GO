// Package server - хранилище для всех параметров, используемых приложением.
package server

import (
	"context"
	"database/sql"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/storages/cookie"
)

// IStorage - интерфейс для добавления, взятия или удаления данны из хранилища.
type IStorage interface {
	SetData(key string, value string, userid string) error
	GetData(value string, param string, userid string) (string, []bool, error)
	DeleteData(delValues []models.DeleteValue) error
}

// Server - структура для хранения всех параметров.
type Server struct {
	sc     config.ServerConfig
	cs     cookie.UserCookie
	st     IStorage
	db     *sql.DB
	ch     chan models.DeleteValue
	ctx    context.Context
	cancel context.CancelFunc
}

// GetCookieStorage возвращает указатель на хранилище куки клиентов.
func (serv *Server) GetCookieStorage() *cookie.UserCookie {
	return &serv.cs
}

// GetStorage возвращает интерфейс.
func (serv *Server) GetStorage() IStorage {
	return serv.st
}

// SetChan отправляет delValue в канал.
func (serv *Server) SetChan(delValue models.DeleteValue) {
	serv.ch <- delValue
}

// Cancel используется при завершении программы для завершения горутины и закрытия канала.
func (serv *Server) Cancel() {
	if serv.db != nil {
		serv.db.Close()
	}
	serv.cancel()
}

// GetServerConfig возвращает указатель на хранилище конфигураций системы.
func (serv *Server) GetServerConfig() *config.ServerConfig {
	return &serv.sc
}

// GetDB возвращает ссылку на БД-хранилище
func (serv *Server) GetDB() *sql.DB {
	return serv.db
}
