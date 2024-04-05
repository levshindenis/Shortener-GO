package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// DelH используется для отправки запроса через хендлер DelURLs.
// Сначала клиента просят ввести количество URL, которые нужно удалить.
// После этого короткие URL вводятся последовательно.
// Полученные URL записываются в массив arr и в формате JSON кладутся в request.Body для отправки запроса.
// В ответ клиент получает статус обработки запроса.
func (s *Server) DelH() {
	var (
		count    int
		shortURL string
		arr      []string
		resp     *resty.Response
		err      error
	)

	fmt.Println("Введите количество URL:")
	fmt.Scanf("%d\n", &count)

	for i := 0; i < count; i++ {
		fmt.Println("Введите URL:")
		fmt.Scanf("%s\n", &shortURL)
		arr = append(arr, shortURL)
	}

	marsh, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}

	if s.cookie != "" {
		resp, err = s.client.R().SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).
			SetBody(bytes.NewBuffer(marsh)).Delete(s.address + "api/user/urls")
	} else {
		resp, err = s.client.R().
			SetBody(bytes.NewBuffer(marsh)).Delete(s.address + "api/user/urls")
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode())

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
