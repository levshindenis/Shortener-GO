package main

import (
	"context"
	"fmt"
	"strings"
)

func (s *Server) GetH() {
	var (
		shortURL string
	)

	fmt.Println("Введите короткий URL:")
	fmt.Scanf("%s\n", &shortURL)

	resp, err := s.client.R().Get(s.address + shortURL)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), "\"")[1])
	} else {
		fmt.Println(resp.StatusCode())
	}

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
