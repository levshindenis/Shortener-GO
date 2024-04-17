// Package handlers используется для обработки запросов к серверу.
package handlers

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/storages/server"
)

// HStorage - основная структура для хендлеров
type HStorage struct {
	server.Server
}
