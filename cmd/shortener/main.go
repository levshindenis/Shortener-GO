package main

import (
	"net/http"

	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/routers"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var server handlers.HStorage
	server.Init()

	return http.ListenAndServe(server.GetConfigParameter("address"), routers.MyRouter(server))
}

//new
