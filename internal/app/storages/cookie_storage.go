package storages

import (
	"github.com/levshindenis/sprint1/internal/app/tools"
	"net/http"
)

type CookieStorage struct {
	arr []string
}

func (co *CookieStorage) GetArr() []string {
	return co.arr
}

func (co *CookieStorage) CountCookies() int {
	return len(co.GetArr())
}

func (co *CookieStorage) SetCookie(value string) {
	co.arr = append(co.GetArr(), value)
}

func (co *CookieStorage) InCookies(value string) bool {
	for ind := range co.GetArr() {
		if co.GetArr()[ind] == value {
			return true
		}
	}
	return false
}

func (co *CookieStorage) CheckCookie(r *http.Request) (string, bool, error) {
	cookie, err := r.Cookie("UserID")
	if err != nil {
		return "non", false, nil
	}
	if !co.InCookies(cookie.Value) {
		gen, err1 := tools.GenerateCookie(co.CountCookies() + 1)
		if err1 != nil {
			return "", false, err
		}
		co.SetCookie(gen)
		return gen, false, nil
	}
	return cookie.Value, true, nil
}
