// пакеты исполняемых приложений должны называться main
package main

import (
	"fmt"
	"net/http"
)

type Storage map[string]string

//func Encode(s string) string {
//
//}

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	var storage Storage
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, PostShortenerHandler(&storage))
	mux.HandleFunc(`/{id}`, GetShortenerHandler(&storage))
	return http.ListenAndServe(`:8080`, mux)
}

func PostShortenerHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			// разрешаем только POST-запросы
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//if r.Header.Get("Content-Type") != "text/plain" {
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		//установим правильный заголовок для типа данных
		//w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Println("Host: ", r.Host)
		fmt.Println("URL: ", r.URL)
		fmt.Println("Body: ", r.Body)
		fmt.Println(storage)
		//fmt.Println("Psth: ", r.URL.Path)
		//fmt.Println("Storage: ", storage)
	}
}

func GetShortenerHandler(storage *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			// разрешаем только POST-запросы
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// установим правильный заголовок для типа данных
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(storage)
	}
}
