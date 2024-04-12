package server

import (
	"context"
	"database/sql"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/models"
	"github.com/levshindenis/sprint1/internal/app/storages/cookie"
	"github.com/levshindenis/sprint1/internal/app/storages/db"
	"github.com/levshindenis/sprint1/internal/app/storages/file"
	"github.com/levshindenis/sprint1/internal/app/storages/memory"
)

// Init используется для инициализации переменных хранилища.
func (serv *Server) Init(conf config.ServerConfig) {
	serv.sc = conf
	serv.cs = cookie.UserCookie{Arr: make([]string, 0)}
	serv.ctx, serv.cancel = context.WithCancel(context.Background())
	serv.InitStorage()
	serv.ch = make(chan models.DeleteValue, 1024)
	go serv.DeleteItems(serv.ctx)
}

// InitStorage используется для выбора хранилища в зависимости от аргументов командной строки.
func (serv *Server) InitStorage() {
	if serv.GetServerConfig().GetDBAddress() != "" {
		serv.InitDB()
	} else if serv.GetServerConfig().GetFilePath() != "" {
		serv.InitFile()
	} else {
		serv.InitMemory()
	}
}

// InitDB используется для определения БД как основного хранилища и создания БД.
func (serv *Server) InitDB() {
	newDB, err := sql.Open("pgx", serv.GetServerConfig().GetDBAddress())
	if err != nil {
		panic(err)
	}
	serv.db = newDB

	dbs := db.Database{DB: newDB}
	dbs.MakeDB()
	serv.st = &dbs
}

// InitFile используется для определения файла как основного хранилища и создания файла-хранилища.
func (serv *Server) InitFile() {
	fl := file.File{Path: serv.GetServerConfig().GetFilePath()}
	fl.MakeFile()
	serv.st = &fl
}

// InitMemory используется для определения памяти как основного хранилища.
func (serv *Server) InitMemory() {
	mem := memory.Memory{Arr: []models.MSItem{}}
	serv.st = &mem
}
