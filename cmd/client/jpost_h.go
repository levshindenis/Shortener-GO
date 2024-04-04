package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/levshindenis/sprint1/internal/app/models"
)

func (s *Server) JPostH() {
	var (
		longURL  string
		shortURL models.JsonEncoder
		resp     *resty.Response
		err      error
	)

	fmt.Println("Введите длинный URL:")
	fmt.Scanf("%s\n", &longURL)

	marsh, err := json.Marshal(models.JsonDecoder{LongURL: longURL})
	if err != nil {
		panic(err)
	}

	if s.cookie != "" {
		resp, err = s.client.R().SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).
			SetBody(bytes.NewBuffer(marsh)).Post(s.address + "api/shorten")
	} else {
		resp, err = s.client.R().SetBody(bytes.NewBuffer(marsh)).Post(s.address + "api/shorten")
	}

	if err != nil {
		panic(err)
	}

	if s.cookie == "" {
		s.cookie = resp.Cookies()[0].Value
	}

	if err = json.Unmarshal(resp.Body(), &shortURL); err != nil {
		panic(err)
	}

	fmt.Println("Ответ:\n", shortURL.ShortURL)

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
