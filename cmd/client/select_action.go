package main

import (
	"context"
	"fmt"
)

// SelectAction используется для выбора запроса.
func (s *Server) SelectAction() {
	fmt.Println("Cookie: ", s.cookie)
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1) Из лонг в шорт")
		fmt.Println("2) Из шорт в лонг")
		fmt.Println("3) Пинг")
		fmt.Println("4) Из лонг в шорт (JSON)")
		fmt.Println("5) Из лонг в шорт (Batch)")
		fmt.Println("6) Вернуть все URL")
		fmt.Println("7) Удалить URLs")
		fmt.Println("8) Статистика")
		fmt.Println("===========================")
		fmt.Print("Ввод:    ")
		fmt.Scanf("%s", &s.choice)
		fmt.Println(s.choice)
		if s.choice == "1" || s.choice == "2" || s.choice == "3" || s.choice == "4" || s.choice == "5" ||
			s.choice == "6" || s.choice == "7" || s.choice == "8" {
			break
		}
		fmt.Println("Bad answer. Please repeat!")
	}

	if err := s.f.Event(context.Background(), s.m[s.choice]); err != nil {
		panic(err)
	}
}
