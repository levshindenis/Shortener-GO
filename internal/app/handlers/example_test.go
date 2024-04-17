package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/models"
)

func ExampleHStorage_SetLongURL() {
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://yandex1.ru"))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	w := httptest.NewRecorder()
	serv.SetLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 201
}

func ExampleHStorage_SetLongURL_second() {
	// Результатом будет статус 409, потому что такой URL уже был сокращен.
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://yandex1.ru"))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	w := httptest.NewRecorder()
	serv.SetLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 409
}

func ExampleHStorage_GetLongURL() {
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()
	serv.GetLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 307
}

func ExampleHStorage_GetLongURL_second() {
	// Результатом будет статус 400, потому что такого короткого URL не существует.
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodGet, "/abcd", nil)
	w := httptest.NewRecorder()
	serv.GetLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 400
}

func ExampleHStorage_GetLongURL_third() {
	// Результатом будет статус 410, потому что данный URL уже был удален
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	//Удаляем запись
	delValue := models.DeleteValue{Userid: "spotify", Value: "abc"}
	serv.GetStorage().DeleteData([]models.DeleteValue{delValue})

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()
	serv.GetLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 410
}

func ExampleHStorage_SetJSONLongURL() {
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Переводим данные в JSON формат
	marsh, _ := json.Marshal(models.JSONDecoder{LongURL: "https://yandex1.ru"})

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(marsh))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	serv.SetJSONLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 201
}

func ExampleHStorage_SetJSONLongURL_second() {
	// Результатом будет статус 409, потому что такой URL уже был сокращен.
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	// Переводим данные в JSON формат
	marsh, _ := json.Marshal(models.JSONDecoder{LongURL: "https://yandex1.ru"})

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(marsh))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	serv.SetJSONLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 409
}

func ExampleHStorage_SetJSONLongURL_third() {
	// Результатом будет статус 400, потому что забыли указать Content-Type : application/json.
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Переводим данные в JSON формат
	marsh, _ := json.Marshal(models.JSONDecoder{LongURL: "https://yandex1.ru"})

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(marsh))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	w := httptest.NewRecorder()
	serv.SetJSONLongURL(w, r)

	fmt.Println(w.Code)

	// Output:
	// 400
}

func ExampleHStorage_BatchURLs() {
	var (
		conf config.ServerConfig
		serv HStorage
		dec  []models.BatchDecoder
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Создаем и переводим данные в JSON формат
	for i := 0; i < 5; i++ {
		myStr := strconv.Itoa(i + 2)
		dec = append(dec, models.BatchDecoder{ID: myStr, LongURL: "https://yandex" + myStr + ".ru/"})
	}
	marsh, _ := json.Marshal(dec)

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(marsh))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	serv.BatchURLs(w, r)

	fmt.Println(w.Code)

	// Output:
	// 201
}

func ExampleHStorage_BatchURLs_second() {
	// Результатом будет статус 400, потому что в request.Body положили некорректные данные
	var (
		conf config.ServerConfig
		serv HStorage
		dec  []string
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Создаем и переводим данные в JSON формат
	for i := 0; i < 5; i++ {
		myStr := strconv.Itoa(i + 2)
		dec = append(dec, "https://yandex"+myStr+".ru/")
	}
	marsh, _ := json.Marshal(dec)

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(marsh))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	serv.BatchURLs(w, r)

	fmt.Println(w.Code)

	// Output:
	// 400
}

func ExampleHStorage_GetURLs() {
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	w := httptest.NewRecorder()
	serv.GetURLs(w, r)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func ExampleHStorage_GetURLs_second() {
	// Результатом будет статус 204, потому что данный пользователь не сокращал URL.
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("abc", "https://yandex1.ru", "spotify")

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "music"})
	w := httptest.NewRecorder()
	serv.GetURLs(w, r)

	fmt.Println(w.Code)

	// Output:
	// 204
}

func ExampleHStorage_DelURLs() {
	var (
		conf config.ServerConfig
		serv HStorage
	)

	// Вместо conf.ParseFlags()
	conf.SetStartAddress("localhost:8080")
	conf.SetShortBaseURL("http://localhost:8080")
	conf.SetDBAddress("")
	conf.SetFilePath("")

	serv.Init(conf)

	// Добавляем запись в хранилище
	serv.GetStorage().SetData("aaa", "https://yandex1.ru", "spotify")

	marsh, _ := json.Marshal([]string{"aaa"})

	// Запрос к хендлеру
	r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(marsh))
	r.AddCookie(&http.Cookie{Name: "UserID", Value: "spotify"})
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	serv.DelURLs(w, r)

	fmt.Println(w.Code)

	// Output:
	// 202
}
