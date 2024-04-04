package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/storages/server"
)

type HStorage struct {
	server.Server
}
