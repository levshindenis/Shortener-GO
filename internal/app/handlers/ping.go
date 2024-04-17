package handlers

import (
	"context"
	"net/http"
	"time"
)

// Ping нужен для обработки запроса от клиента по адресу /ping.
// Хендлер проверяет, есть ли соединение с базой данных.
// При отсутствии соединения возвращается StatusInternalServerError.
func (serv *HStorage) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := serv.Server.GetDB().PingContext(ctx); err != nil {
		http.Error(w, "Something bad with ping", http.StatusInternalServerError)
		return
	}
}
