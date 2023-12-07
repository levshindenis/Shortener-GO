package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStorage_PostHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		requestBody  string
		expectedCode int
		emptyBody    bool
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			requestBody:  "https://yandex.ru/",
			expectedCode: http.StatusCreated,
			emptyBody:    false,
		},
		{
			name:         "Bad method",
			method:       http.MethodGet,
			requestBody:  "https://yandex.ru/",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
		{
			name:         "Bad url",
			method:       http.MethodPost,
			requestBody:  "Hello",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var storage Storage
			storage.EmptyStorage()
			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			storage.PostHandler(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if !tt.emptyBody {
				assert.Contains(t, w.Body.String(), "http://localhost:8080/",
					"Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestStorage_GetHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedBody string
		emptyBody    bool
	}{
		{
			name:         "Good test",
			method:       http.MethodGet,
			url:          "GyuRe0",
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "https://yandex.ru/",
			emptyBody:    false,
		},
		{
			name:         "Bad method",
			method:       http.MethodPost,
			url:          "GyuRe0",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			emptyBody:    true,
		},
		{
			name:         "No url",
			method:       http.MethodGet,
			url:          "GyuAe0",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			emptyBody:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var storage Storage
			storage.EmptyStorage()
			storage["GyuRe0"] = "https://yandex.ru/"
			reqURL := "/" + tt.url
			r := httptest.NewRequest(tt.method, reqURL, nil)
			w := httptest.NewRecorder()
			storage.GetHandler(w, r)

			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if !tt.emptyBody {
				assert.Equal(t, w.Header().Get("Location"), tt.expectedBody,
					"Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
