package middleware

import (
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"net/http"
)

type cookieWriter struct {
	http.ResponseWriter
	status int
}

func (cw *cookieWriter) WriteHeader(status int) {
	cw.status = status
	cw.ResponseWriter.WriteHeader(status)
}

func WithCookie(next http.HandlerFunc, serv handlers.HStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cw := &cookieWriter{ResponseWriter: w, status: http.StatusCreated}
		cookVal, cookFlag, err := serv.GetCookieStorage().CheckCookie(r)
		if err != nil {
			cw.status = http.StatusBadRequest
		}
		if !cookFlag {
			if r.Method == http.MethodPost {
				mystr := "UserID=" + cookVal
				r.Header.Set("Cookie", mystr)
				http.SetCookie(cw, &http.Cookie{
					Name:  "UserID",
					Value: cookVal,
				})
			} else {
				cw.status = http.StatusUnauthorized
			}
		}
		next.ServeHTTP(cw, r)
	}
}
