package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

func (serv *HStorage) Ping(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", serv.GetServerConfig("db"))
	if err != nil {
		http.Error(w, "Something bad with open db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		http.Error(w, "Something bad with ping", http.StatusInternalServerError)
		return
	}
}
