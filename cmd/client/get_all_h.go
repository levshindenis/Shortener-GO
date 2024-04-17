package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// GetAllH используется для отправки запроса через хендлер GetURLs.
// В ответ клиент получает все свои сокращенные URL.
func (s *Server) GetAllH() {
	var (
		resp *resty.Response
		err  error
	)

	if s.cookie != "" {
		resp, err = s.client.R().SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).
			Get(s.address + "api/user/urls")
	} else {
		resp, err = s.client.R().Get(s.address + "api/user/urls")
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
