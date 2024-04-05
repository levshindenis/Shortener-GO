// Package shortener - основной пакет
package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/router"
)

func main() {
	var (
		server handlers.HStorage
		conf   config.ServerConfig
	)

	conf.ParseFlags()
	server.Init(conf)
	if err := http.ListenAndServe(conf.GetStartAddress(), router.MyRouter(&server)); err != nil {
		panic(err)
	}

	server.CancelCh()
}
