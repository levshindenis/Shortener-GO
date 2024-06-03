// Package client используется для отправки запросов со стороны клиента. Испульзуется finite state machine.
package main

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/looplab/fsm"
)

// Server - основная структура для хранения значений сервера и клиента
type Server struct {
	client  *resty.Client
	cookie  string
	address string
	choice  string
	m       map[string]string
	f       *fsm.FSM
}

// NewServer - функция для создания нового Server
func NewServer() *Server {
	client := resty.New()
	m := map[string]string{
		"1": "post",
		"2": "get",
		"3": "ping",
		"4": "jPost",
		"5": "bPost",
		"6": "getAll",
		"7": "del",
		"8": "stat",
	}
	return &Server{
		client:  client,
		cookie:  "",
		address: "http://localhost:8080/",
		choice:  "",
		m:       m,
	}
}

func main() {
	server := NewServer()
	server.f = fsm.NewFSM(
		"zero",
		fsm.Events{
			{Name: "go", Src: []string{"zero"}, Dst: "main"},
			{Name: "post", Src: []string{"main"}, Dst: "postH"},
			{Name: "get", Src: []string{"main"}, Dst: "getH"},
			{Name: "ping", Src: []string{"main"}, Dst: "pingH"},
			{Name: "jPost", Src: []string{"main"}, Dst: "jPostH"},
			{Name: "bPost", Src: []string{"main"}, Dst: "bPostH"},
			{Name: "getAll", Src: []string{"main"}, Dst: "getAllH"},
			{Name: "del", Src: []string{"main"}, Dst: "delH"},
			{Name: "stat", Src: []string{"main"}, Dst: "statH"},
			{Name: "mainpage",
				Src: []string{"postH", "getH", "pingH", "jPostH", "bPostH", "getAllH", "delH", "statH"},
				Dst: "main"},
		},
		fsm.Callbacks{
			"main":   func(_ context.Context, _ *fsm.Event) { server.SelectAction() },
			"post":   func(_ context.Context, _ *fsm.Event) { server.PostH() },
			"get":    func(_ context.Context, _ *fsm.Event) { server.GetH() },
			"ping":   func(_ context.Context, _ *fsm.Event) { server.PingH() },
			"jPost":  func(_ context.Context, _ *fsm.Event) { server.JPostH() },
			"bPost":  func(_ context.Context, _ *fsm.Event) { server.BPostH() },
			"getAll": func(_ context.Context, _ *fsm.Event) { server.GetAllH() },
			"del":    func(_ context.Context, _ *fsm.Event) { server.DelH() },
			"stat":   func(_ context.Context, _ *fsm.Event) { server.StatsH() },
		},
	)

	if err := server.f.Event(context.Background(), "go"); err != nil {
		panic(err)
	}
}
