// Package cookie - хранилище для куки клиентов.
package cookie

// UserCookie - структура для хранение куки клиентов.
type UserCookie struct {
	Arr []string
}

// GetArr - возвращает массив куки.
func (co *UserCookie) GetArr() []string {
	return co.Arr
}

// SetValue - добавляет значение в массив куки.
func (co *UserCookie) SetValue(value string) {
	co.Arr = append(co.Arr, value)
}
