package cookie

// InCookies  проверяет, есть ли value в массиве куки.
// Если есть, то возвращает true, иначе - false.
func (co *UserCookie) InCookies(value string) bool {
	for ind := range co.Arr {
		if co.Arr[ind] == value {
			return true
		}
	}
	return false
}
