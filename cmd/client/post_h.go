package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func (s *Server) PostH() {
	var (
		longURL string
		resp    *resty.Response
		err     error
	)

	fmt.Println("Введите длинный URL:")
	fmt.Scanf("%s\n", &longURL)

	if s.cookie != "" {
		resp, err = s.client.R().SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).
			SetBody(longURL).Post(s.address)
	} else {
		resp, err = s.client.R().SetBody(longURL).Post(s.address)
	}

	if err != nil {
		panic(err)
	}

	if s.cookie == "" {
		s.cookie = resp.Cookies()[0].Value
	}

	fmt.Println("Ответ:\n", resp)

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
