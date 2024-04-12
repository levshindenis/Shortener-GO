package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/levshindenis/sprint1/internal/app/models"
	"net/http"
)

// BatchURLs - нужен для обработки запроса от клиента по адресу /api/shorten/batch.
// Сначала проверяются входящие данные на JSON формат.
// При успешной проверке длинные URL считываются из request.Body и преобразуются из формата JSON.
// Создаются короткие URL и все данные (длинный URL, короткий URL, cookie) сохраняются.
// Созданные короткие URL возвращаются клиенту в формате JSON.
// При успешной обработке запроса устанавливается StatusCreated.
func (serv *HStorage) BatchURLs(w http.ResponseWriter, r *http.Request) {
	var (
		dec []models.BatchDecoder
		buf bytes.Buffer
	)

	cookie, _ := r.Cookie("UserID")
	http.SetCookie(w, &http.Cookie{Name: "UserID", Value: cookie.Value})

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "Something bad with read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buf.Bytes(), &dec); err != nil {
		http.Error(w, "Something bad with decoding JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := make([]models.BatchEncoder, len(dec))
	for ind := range dec {
		short, flag, err := serv.MakeShortURL(dec[ind].LongURL)
		if err != nil {
			http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
			return
		}

		if !flag {
			if err = serv.GetStorage().SetData(short, dec[ind].LongURL, cookie.Value); err != nil {
				http.Error(w, "Somethings bad with Save", http.StatusBadRequest)
				return
			}
		}

		short = serv.GetServerConfig().GetShortBaseURL() + "/" + short
		enc[ind] = models.BatchEncoder{ID: dec[ind].ID, ShortURL: short}
	}

	resp, err := json.Marshal(enc)
	if err != nil {
		http.Error(w, "Something bad with encoding JSON", http.StatusBadRequest)
		return
	}

	if _, err = w.Write(resp); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}
