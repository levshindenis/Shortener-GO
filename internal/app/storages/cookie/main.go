package cookie

type UserCookie struct {
	Arr []string
}

func (co *UserCookie) GetArr() []string {
	return co.Arr
}

func (co *UserCookie) SetValue(value string) {
	co.Arr = append(co.Arr, value)
}
