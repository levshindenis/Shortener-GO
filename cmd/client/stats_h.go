package main

import (
	"context"
	"fmt"
)

// StatsH - функция для отправки клиентом запроса к серверу по статистике
func (s *Server) StatsH() {
	var (
		ip string
	)

	fmt.Println("Введите IP:")
	fmt.Scanf("%s\n", &ip)

	resp, err := s.client.R().SetHeader("X-Real-IP", "").Get(s.address + "api/internal/stats")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
