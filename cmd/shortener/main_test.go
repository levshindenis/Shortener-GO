package main

import (
	"bytes"
	"encoding/json"
	"github.com/levshindenis/sprint1/internal/app/models"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/handlers"
)

var (
	serv      handlers.HStorage
	conf      config.ServerConfig
	arr1      []string
	arr2      []string
	cookieArr []string
)

func TestMain(m *testing.M) {
	cookieArr = []string{"abcdefg", "zxcvbn", "spotifyLU", "spotifyJLU", "spotifyBLU", "aaaaaa"}

	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	os.Exit(m.Run())
}

func TestHSStorage_SetLongURL(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		cookVal      string
	}{
		{
			name:         "Good test_1",
			requestBody:  "https://yandex.ru/",
			expectedCode: http.StatusCreated,
			cookVal:      "abcdefg",
		},
		{
			name:         "Good test_2",
			requestBody:  "https://yandex1.ru/",
			expectedCode: http.StatusCreated,
			cookVal:      "zxcvbn",
		},
		{
			name:         "Repeat",
			requestBody:  "https://yandex.ru/",
			expectedCode: http.StatusConflict,
			cookVal:      "abcdefg",
		},
		{
			name:         "Bad url",
			requestBody:  "Hello",
			expectedCode: http.StatusBadRequest,
			cookVal:      "abcdefg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookVal})
			w := httptest.NewRecorder()
			serv.SetLongURL(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if w.Code == 201 {
				arr1 = append(arr1, strings.Split(w.Body.String(), "/")[3])
			}
		})
	}
}

func TestHSStorage_DelURLs(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  []string
		expectedCode int
		cookVal      string
		jsonStatus   bool
	}{
		{
			name:         "Good test",
			requestBody:  []string{arr1[0]},
			expectedCode: http.StatusAccepted,
			cookVal:      "abcdefg",
			jsonStatus:   true,
		},
		{
			name:         "Bad test",
			requestBody:  []string{"Hello"},
			expectedCode: http.StatusBadRequest,
			cookVal:      "abcdefg",
			jsonStatus:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marsh, _ := json.Marshal(tt.requestBody)
			r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(marsh))
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookVal})
			if tt.jsonStatus {
				r.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			serv.DelURLs(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
	time.Sleep(3 * time.Second)
}

func TestHSStorage_GetLongURL(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		shortURL     string
	}{
		{
			name:         "Good test_1",
			expectedCode: http.StatusTemporaryRedirect,
			shortURL:     arr1[1],
		},
		{
			name:         "Good test_2",
			expectedCode: http.StatusGone,
			shortURL:     arr1[0],
		},
		{
			name:         "Bad url",
			expectedCode: http.StatusBadRequest,
			shortURL:     "aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
			w := httptest.NewRecorder()
			serv.GetLongURL(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHSStorage_SetJsonLongURL(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		cookVal      string
	}{
		{
			name:         "Good test_1",
			requestBody:  `{"url":"https://yandex2.ru/"}`,
			expectedCode: http.StatusCreated,
			cookVal:      "abcdefg",
		},
		{
			name:         "Good test_2",
			requestBody:  `{"url":"https://yandex4.ru/"}`,
			expectedCode: http.StatusCreated,
			cookVal:      "zxcvbn",
		},
		{
			name:         "Repeat",
			requestBody:  `{"url":"https://yandex.ru/"}`,
			expectedCode: http.StatusConflict,
			cookVal:      "abcdefg",
		},
		{
			name:         "Bad url",
			requestBody:  "Hello",
			expectedCode: http.StatusBadRequest,
			cookVal:      "abcdefg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(tt.requestBody)))
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookVal})
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			serv.SetJSONLongURL(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if w.Code == 201 {
				arr1 = append(arr1, strings.Split(w.Body.String(), "/")[3][:6])
			}
		})
	}
}

func TestHSStorage_BatchURLs(t *testing.T) {
	tests := []struct {
		name         string
		ids          []string
		requestBody  []string
		expectedCode int
		cookVal      string
		jsonStatus   bool
	}{
		{
			name:         "Good test_1",
			ids:          []string{"1", "2"},
			requestBody:  []string{"https://yandex3.ru/", "https://yandex4.ru/"},
			expectedCode: http.StatusCreated,
			cookVal:      "abcdefg",
			jsonStatus:   true,
		},
		{
			name:         "Good test_2",
			ids:          []string{"3", "4"},
			requestBody:  []string{"https://yandex5.ru", "https://yandex6.ru"},
			expectedCode: http.StatusCreated,
			cookVal:      "zxcvbn",
			jsonStatus:   true,
		},
		{
			name:         "Bad url",
			ids:          []string{"5"},
			requestBody:  []string{"https://yandex6.ru/"},
			expectedCode: http.StatusBadRequest,
			cookVal:      "abcdefg",
			jsonStatus:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bDec []models.BatchDecoder
			for i, v := range tt.ids {
				bDec = append(bDec, models.BatchDecoder{ID: v, LongURL: tt.requestBody[i]})
			}
			marsh, _ := json.Marshal(bDec)
			r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(marsh))
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookVal})
			if tt.jsonStatus {
				r.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			serv.BatchURLs(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func TestHSStorage_GetURLs(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		cookVal      string
	}{
		{
			name:         "Good test_1",
			expectedCode: http.StatusOK,
			cookVal:      "abcdefg",
		},
		{
			name:         "Good test_2",
			expectedCode: http.StatusOK,
			cookVal:      "zxcvbn",
		},
		{
			name:         "No URLs",
			expectedCode: http.StatusNoContent,
			cookVal:      "aaa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			r.AddCookie(&http.Cookie{Name: "UserID", Value: tt.cookVal})
			w := httptest.NewRecorder()
			serv.GetURLs(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func BenchmarkHSStorage_SetLongURL(b *testing.B) {
	cookval := "spotifyLU"
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reqBody := "https://yandex" + strconv.Itoa(i+100) + ".ru/"
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		r.AddCookie(&http.Cookie{Name: "UserID", Value: cookval})
		w := httptest.NewRecorder()
		b.StartTimer()

		serv.SetLongURL(w, r)

		b.StopTimer()
		if w.Code == 201 {
			arr2 = append(arr2, strings.Split(w.Body.String(), "/")[3])
		}
		b.StartTimer()
	}
}

func BenchmarkHSStorage_GetLongURL(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		n := rnd.Intn(len(arr2))
		r := httptest.NewRequest(http.MethodGet, "/"+arr2[n], nil)
		w := httptest.NewRecorder()
		b.StartTimer()

		serv.GetLongURL(w, r)
	}
}

func BenchmarkHSStorage_SetJsonLongURL(b *testing.B) {
	cookval := "spotifyJLU"

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reqBody := "https://yande" + strconv.Itoa(i+100) + "x.ru/"
		marsh, _ := json.Marshal(models.JSONDecoder{LongURL: reqBody})
		r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(marsh))
		r.AddCookie(&http.Cookie{Name: "UserID", Value: cookval})
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		b.StartTimer()

		serv.SetJSONLongURL(w, r)

		b.StopTimer()
		if w.Code == 201 {
			arr2 = append(arr2, strings.Split(w.Body.String(), "/")[3][:6])
		}
		b.StartTimer()
	}
}

func BenchmarkHSStorage_BatchURLs(b *testing.B) {
	cookval := "spotifyBLU"

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var dec []models.BatchDecoder
		for j := 0; j < 5; j++ {
			reqBody := "https://yand" + strconv.Itoa(i*5+j) + "ex.ru/"
			dec = append(dec, models.BatchDecoder{ID: strconv.Itoa(i*5 + j), LongURL: reqBody})
		}
		marsh, _ := json.Marshal(dec)
		r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(marsh))
		r.AddCookie(&http.Cookie{Name: "UserID", Value: cookval})
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		b.StartTimer()

		serv.BatchURLs(w, r)
	}
}

func BenchmarkHSStorage_GetURLs(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ind := rnd.Intn(len(cookieArr))
		cookval := cookieArr[ind]
		r := httptest.NewRequest(http.MethodGet, "/user/urls", nil)
		r.AddCookie(&http.Cookie{Name: "UserID", Value: cookval})
		w := httptest.NewRecorder()
		b.StartTimer()

		serv.GetURLs(w, r)
	}
}

func BenchmarkHSStorage_DelURLs(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		indC := rnd.Intn(len(cookieArr))
		cookval := cookieArr[indC]
		var newArr []string
		for j := 0; j < 5; j++ {
			incS := rnd.Intn(len(arr2))
			newArr = append(newArr, arr2[incS])
		}
		marsh, _ := json.Marshal(newArr)

		r := httptest.NewRequest(http.MethodDelete, "/user/urls", bytes.NewReader(marsh))
		r.AddCookie(&http.Cookie{Name: "UserID", Value: cookval})
		w := httptest.NewRecorder()
		b.StartTimer()

		serv.DelURLs(w, r)
	}
}
