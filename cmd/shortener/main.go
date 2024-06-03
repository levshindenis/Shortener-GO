// Package shortener - основной пакет
package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/levshindenis/sprint1/cmd/proto/shortener"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/router"
	"github.com/levshindenis/sprint1/internal/app/tools"
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
		gs     *grpc.Server
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

	listen, err := net.Listen("tcp", conf.GetStartAddressG())
	if err != nil {
		log.Fatal(err)
	}

	if !conf.GetGTLS() {
		gs = grpc.NewServer()
	} else {
		tlsCreds, err := tools.GenerateTLSCreds()
		if err != nil {
			panic(err)
		}
		gs = grpc.NewServer(grpc.Creds(tlsCreds))
	}

	reflection.Register(gs)
	pb.RegisterShortenerServer(gs, &ShortenerServer{serv: &server.Server})

	go func() {
		if err = gs.Serve(listen); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
		gs.GracefulStop()
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
		if err = srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}

	<-exitProgram

	server.Cancel()
}
