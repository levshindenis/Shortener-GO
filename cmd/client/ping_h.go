package main

import (
	"context"
	"fmt"
)

func (s *Server) PingH() {
	resp, err := s.client.R().Get(s.address + "ping")
	if err != nil {
		panic(err)
	}

	fmt.Println("Ответ:\n", resp.StatusCode())

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
