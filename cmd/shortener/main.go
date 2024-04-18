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
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	var (
		server handlers.HStorage
		conf   config.ServerConfig
	)
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s", buildVersion, buildDate, buildCommit)

	conf.ParseFlags()
	server.Init(conf)
	if err := http.ListenAndServe(conf.GetStartAddress(), router.MyRouter(&server)); err != nil {
		panic(err)
	}

	server.Cancel()
}
