package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData - хранит в себе статус и размер ответа на запрос
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter - структура для перезаписывания Writer
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// MyLogger - мой логгер
type MyLogger struct {
	loggerSugar zap.SugaredLogger
}

// Init нужен для создания логгера
func (ml *MyLogger) Init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	ml.loggerSugar = *logger.Sugar()
}

// Write - перезаписывает основной одноименный метод
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader - перезаписывает заголовок ответа
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging - middleware для логирования поступающих запросов.
func (ml *MyLogger) WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		ml.loggerSugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
}
