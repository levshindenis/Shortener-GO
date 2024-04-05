package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipWriter - основная структура для сжатия исходящих данных
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write - метод, который перезаписывает базовый Write
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WithCompression - middleware для сжатия данных, передаваемых клиенту.
// Сначала проверяется, есть ли в заголовке gzip.
// Если нет, то запрос передается дальше.
// Если есть, то переписывается ResponseWriter на сжатие данных.
func WithCompression(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") ||
			!strings.Contains(r.Header.Get("Content-Type"), "application/json") &&
				!strings.Contains(r.Header.Get("Content-Type"), "text/html") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, "Something bad with gzip", http.StatusBadRequest)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
}
