package memory

import "github.com/levshindenis/sprint1/internal/app/models"

// DeleteData нужа для "удаления" переданных сокращенных URL из Memory.
// В цикле берется каждый короткий URL и по нему сравнивается UserID из поступивших данных и значение из Memory.
// Если значения не совпадают, то удаление не происходит.
// Если значения совпали, то меняется значение "deleted" на true.
func (ms *Memory) DeleteData(delValues []models.DeleteValue) error {
	for _, elem := range delValues {
		for ind, msi := range ms.Arr {
			if msi.Key == elem.Value && msi.UserID == elem.Userid {
				ms.Arr[ind].Deleted = true
			}
		}
	}
	return nil
}
