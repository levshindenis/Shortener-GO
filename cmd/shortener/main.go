// Package shortener - основной пакет
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

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
	srv := &http.Server{Addr: conf.GetStartAddress(), Handler: router.MyRouter(&server)}

	exitProgram := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
		close(exitProgram)
	}()

	if !conf.GetHTTPS() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	} else {
		manager := &autocert.Manager{
			Cache:  autocert.DirCache("cache-dir"),
			Prompt: autocert.AcceptTOS,
		}
		srv.TLSConfig = manager.TLSConfig()
		if err := srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}

	<-exitProgram

	server.Cancel()
}
