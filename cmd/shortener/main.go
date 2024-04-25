// Package shortener - основной пакет
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"golang.org/x/crypto/acme/autocert"

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

	if err := conf.ParseFlags(); err != nil {
		panic(err)
	}
	server.Init(conf)

	if !conf.GetHTTPS() {
		if err := http.ListenAndServe(conf.GetStartAddress(), router.MyRouter(&server)); err != nil {
			panic(err)
		}
	} else {
		manager := &autocert.Manager{
			Cache:  autocert.DirCache("cache-dir"),
			Prompt: autocert.AcceptTOS,
		}
		HTTPSServer := &http.Server{
			Addr:      conf.GetStartAddress(),
			Handler:   router.MyRouter(&server),
			TLSConfig: manager.TLSConfig(),
		}
		if err := HTTPSServer.ListenAndServeTLS("", ""); err != nil {
			panic(err)
		}
	}

	server.Cancel()
}
