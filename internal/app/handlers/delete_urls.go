package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// DelURLs нужен для обработки запроса от клиента по адресу /api/user/urls.
// Сначала проверяются входящие данные на JSON формат.
// При успешной проверке сокращенные URL из request.Body преобразуются из JSON.
// Берется куки клиента и в паре с каждым полученным сокращенным URL отправляется в канал.
// При успешной обработке запроса устанавливается StatusAccepted.
func (serv *HStorage) DelURLs(w http.ResponseWriter, r *http.Request) {
	var (
		buf       bytes.Buffer
		shortURLS []string
	)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := json.Unmarshal(buf.Bytes(), &shortURLS)
	if err != nil {
		http.Error(w, "Something bad with Unmarshal", http.StatusBadRequest)
		return
	}

	cookie, _ := r.Cookie("UserID")

	for ind := range shortURLS {
		serv.SetChan(models.DeleteValue{Value: shortURLS[ind], Userid: cookie.Value})
	}

	w.WriteHeader(http.StatusAccepted)
}
