package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/levshindenis/sprint1/internal/app/models"
)

func (serv *HStorage) BatchURLs(w http.ResponseWriter, r *http.Request) {
	var (
		enc []models.BatchEncoder
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

	for _, elem := range dec {
		short, flag, err := serv.MakeShortURL(elem.LongURL)
		if err != nil {
			http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
			return
		}

		if !flag {
			if err = serv.GetStorage().SetData(short, elem.LongURL, cookie.Value); err != nil {
				http.Error(w, "Something bad with Save", http.StatusBadRequest)
				return
			}
		}

		short = serv.GetServerConfig("baseURL") + "/" + short
		enc = append(enc, models.BatchEncoder{ID: elem.ID, ShortURL: short})
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
