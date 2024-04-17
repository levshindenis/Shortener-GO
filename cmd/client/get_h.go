package main

import (
	"context"
	"fmt"
	"strings"
)

// GetH используется для отправки запроса через хендлер GetLongURL.
// Клиент вводит короткий URL, который добавляется к адресу запроса для его отправки.
// В ответ клиент получает длинный URL.
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
