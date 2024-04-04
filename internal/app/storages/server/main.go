package server

import (
	"context"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/storages/cookie"
)

type IStorage interface {
	SetData(key string, value string, userid string) error
	GetData(value string, param string, userid string) (string, []bool, error)
	DeleteData(delValues []models.DeleteValue) error
}

type Server struct {
	sc     config.ServerConfig
	cs     cookie.UserCookie
	st     IStorage
	ch     chan models.DeleteValue
	ctx    context.Context
	cancel context.CancelFunc
}

func (serv *Server) GetCookieStorage() *cookie.UserCookie {
	return &serv.cs
}

func (serv *Server) GetStorage() IStorage {
	return serv.st
}

func (serv *Server) SetChan(delValue models.DeleteValue) {
	serv.ch <- delValue
}

func (serv *Server) CancelCh() {
	serv.cancel()
}
