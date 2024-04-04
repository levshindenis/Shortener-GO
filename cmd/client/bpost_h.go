package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"

	"github.com/levshindenis/sprint1/internal/app/models"
)

func (s *Server) BPostH() {
	var (
		id   int
		dec  []models.BatchDecoder
		resp *resty.Response
		err  error
	)

	fmt.Println("Введите первый ID: ")
	fmt.Scanf("%d\n", &id)

	for i := 0; i < 5; i++ {
		myStr := strconv.Itoa(i + id)
		dec = append(dec, models.BatchDecoder{ID: myStr, LongURL: "https://yandex" + myStr + ".ru/"})
	}

	marsh, err := json.Marshal(dec)
	if err != nil {
		panic(err)
	}

	if s.cookie != "" {
		resp, err = s.client.R().SetCookie(&http.Cookie{Name: "UserID", Value: s.cookie}).
			SetBody(bytes.NewBuffer(marsh)).Post(s.address + "api/shorten/batch")
	} else {
		resp, err = s.client.R().SetBody(bytes.NewBuffer(marsh)).Post(s.address + "api/shorten/batch")
	}

	if err != nil {
		panic(err)
	}

	if s.cookie == "" {
		s.cookie = resp.Cookies()[0].Value
	}

	fmt.Println(resp)

	if err = s.f.Event(context.Background(), "mainpage"); err != nil {
		panic(err)
	}
}
