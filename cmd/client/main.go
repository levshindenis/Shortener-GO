package main

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/looplab/fsm"
)

type Server struct {
	client  *resty.Client
	cookie  string
	address string
	choice  string
	m       map[string]string
	f       *fsm.FSM
}

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
			{Name: "mainpage",
				Src: []string{"postH", "getH", "pingH", "jPostH", "bPostH", "getAllH", "delH"},
				Dst: "main"},
		},
		fsm.Callbacks{
			"main":    func(_ context.Context, _ *fsm.Event) { server.SelectAction() },
			"post":    func(_ context.Context, _ *fsm.Event) { server.PostH() },
			"get":     func(_ context.Context, _ *fsm.Event) { server.GetH() },
			"ping":    func(_ context.Context, _ *fsm.Event) { server.PingH() },
			"jPost":   func(_ context.Context, _ *fsm.Event) { server.JPostH() },
			"bPostH":  func(_ context.Context, _ *fsm.Event) { server.BPostH() },
			"getAllH": func(_ context.Context, _ *fsm.Event) { server.GetAllH() },
			"delH":    func(_ context.Context, _ *fsm.Event) { server.DelH() },
		},
	)

	if err := server.f.Event(context.Background(), "go"); err != nil {
		panic(err)
	}
}
