package cookie

func (co *UserCookie) InCookies(value string) bool {
	for ind := range co.Arr {
		if co.Arr[ind] == value {
			return true
		}
	}
	return false
}
