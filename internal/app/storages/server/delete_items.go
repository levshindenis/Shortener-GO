package server

import (
	"context"
	"time"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// DeleteItems - горутина, которая берет поступившие данные из канала и записывает их в массив values.
// Каждый тик таймера все данные из values идут на вход функции DeleteData (удаление данных).
// Если при тике массив пустой, то ничего не происходит.
func (serv *Server) DeleteItems(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)

	var values []models.DeleteValue

	for {
		select {
		case <-ctx.Done():
			serv.ch <- models.DeleteValue{}
			return
		case value := <-serv.ch:
			values = append(values, value)
		case <-ticker.C:
			if len(values) == 0 {
				continue
			}
			err := serv.st.DeleteData(values)
			if err != nil {
				panic(err)
			}
			values = nil
		}
	}
}
