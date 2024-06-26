// Package middleware - пакет с используемыми middleware
package middleware

import (
	"net/http"

	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

// CheckCookie - middleware для проверки куки у клиента.
// Если куки у клиента отсутствует или такого куки нет в хранилище, то
// Если это Get или Delete запрос, то возвращается StatusUnauthorized.
// Если это был Post запрос, то создается куки и передается дальше по запросу.
func CheckCookie(next http.HandlerFunc, hs *handlers.HStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("UserID")
		if err != nil || !hs.GetCookieStorage().InCookies(cookie.Value) {
			if r.Method == http.MethodGet || r.Method == http.MethodDelete {
				http.Error(w, "Not cookie", http.StatusUnauthorized)
				return
			}
			value, err1 := tools.GenerateCookie(len(hs.GetCookieStorage().GetArr()) + 1)
			if err1 != nil {
				http.Error(w, "Something bad with GenerateCookie", http.StatusBadRequest)
				return
			}

			r.AddCookie(&http.Cookie{Name: "UserID", Value: value})
			hs.GetCookieStorage().SetValue(value)
		}
		next.ServeHTTP(w, r)
	}
}
