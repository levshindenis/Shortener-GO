package main

import (
	"net/http"

	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/routers"
)

func main() {
	var server handlers.HStorage
	server.Init()
	if err := run(server); err != nil {
		panic(err)
	}
	server.CancelCh()
}

func run(server handlers.HStorage) error {
	return http.ListenAndServe(server.GetServerConfig("address"), routers.MyRouter(server))
}
