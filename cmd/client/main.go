package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	for {
		fmt.Println("Выберите действие:")
		fmt.Println("1) Ввести длинный URL")
		fmt.Println("2) Ввести короткий URL")
		fmt.Println("======================")
		fmt.Print("Ввод: ")

		reader := bufio.NewReader(os.Stdin)
		choice, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = strings.TrimSuffix(choice, "\n")

		switch choice {
		case "1":
			fmt.Println("\nВведите длинный URL: ")
		case "2":
			fmt.Println("\nВведите короткий URL: ")
		default:
			os.Exit(1)
		}

		value, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		value = strings.TrimSuffix(value, "\n")
		fmt.Println("\nОтвет:")
		client := resty.New()

		switch choice {
		case "1":
			resp, err := client.R().SetBody(value).Post("http://localhost:8080/")
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)
		case "2":
			myURL := "http://localhost:8080/" + value
			resp, err := client.R().Get(myURL)
			if err != nil {
				fmt.Println("Err")
				fmt.Println(strings.Split(err.Error(), "\"")[1])
			} else {
				if resp.RawResponse.Request.Referer() == myURL {
					fmt.Println("NN Referer")
					fmt.Println(resp.RawResponse.Request.URL)
				} else if resp.RawResponse.Request.Referer() != "" {
					fmt.Println("N Referer")
					fmt.Println(resp.RawResponse.Request.Referer())
				} else {
					fmt.Println("Невозможно преобразовать короткий URL")
				}
			}
		}

		fmt.Println("\nПродолжить?")
		fmt.Println("1) Да")
		fmt.Println("2) Нет")
		fmt.Println("===========")
		fmt.Print("Ответ: ")

		choice, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = strings.TrimSuffix(choice, "\n")

		switch choice {
		case "1":
		default:
			os.Exit(0)
		}
	}
}
