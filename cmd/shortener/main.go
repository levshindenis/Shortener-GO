// Package shortener - основной пакет
package main

import (
	"fmt"
	"github.com/levshindenis/sprint1/internal/app/router"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	_ "net/http/pprof"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/handlers"
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
	server.Init(conf)

	if conf.GetHTTPS() == "" {
		if err := http.ListenAndServe(conf.GetStartAddress(), router.MyRouter(&server)); err != nil {
			panic(err)
		}
	} else {
		manager := &autocert.Manager{
			Cache:  autocert.DirCache("cache-dir"),
			Prompt: autocert.AcceptTOS,
		}
		HTTPSServer := &http.Server{
			Addr:      ":443",
			Handler:   router.MyRouter(&server),
			TLSConfig: manager.TLSConfig(),
		}
		if err := HTTPSServer.ListenAndServeTLS("", ""); err != nil {
			panic(err)
		}
	}

	server.Cancel()
}
