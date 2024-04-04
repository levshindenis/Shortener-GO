package server

import (
	"context"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/storages/cookie"
	"github.com/levshindenis/sprint1/internal/app/storages/db"
	"github.com/levshindenis/sprint1/internal/app/storages/file"
	"github.com/levshindenis/sprint1/internal/app/storages/memory"
)

func (serv *Server) Init(conf config.ServerConfig) {
	serv.sc = conf
	serv.cs = cookie.UserCookie{Arr: make([]string, 0)}
	serv.ctx, serv.cancel = context.WithCancel(context.Background())
	serv.InitStorage()
	serv.ch = make(chan models.DeleteValue, 1024)
	go serv.DeleteItems(serv.ctx)
}

func (serv *Server) InitStorage() {
	if serv.GetServerConfig("db") != "" {
		serv.InitDB()
	} else if serv.GetServerConfig("file") != "" {
		serv.InitFile()
	} else {
		serv.InitMemory()
	}
}

func (serv *Server) InitDB() {
	dbs := db.Database{Address: serv.GetServerConfig("db")}
	dbs.MakeDB()
	serv.st = &dbs
}

func (serv *Server) InitFile() {
	fl := file.File{Path: serv.GetServerConfig("file")}
	fl.MakeFile()
	serv.st = &fl
}

func (serv *Server) InitMemory() {
	mem := memory.Memory{Arr: []models.MSItem{}}
	serv.st = &mem
}
