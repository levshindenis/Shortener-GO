// Package shortener - основной пакет
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/router"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	var (
		server handlers.HStorage
		conf   config.ServerConfig
	)

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	conf.ParseFlags()
	if conf.GetHttps() != "" {
		fmt.Println("TLS:  ", conf.GetHttps())
	}
	server.Init(conf)
	if err := http.ListenAndServe(conf.GetStartAddress(), router.MyRouter(&server)); err != nil {
		panic(err)
	}

	server.Cancel()
}
