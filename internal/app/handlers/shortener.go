package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/levshindenis/sprint1/internal/app/storages"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type HStorage struct {
	storages.ServerStorage
}

func (serv *HStorage) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	var cookVal string
	cookie, err := r.Cookie("UserID")
	if err != nil {
		gen, err1 := tools.GenerateCookie(serv.CountCookies() + 1)
		if err1 != nil {
			http.Error(w, "Something bad with cookies", http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: gen,
		})
		serv.SetCookie(gen)
		cookVal = gen
		fmt.Println(gen)
	} else {
		if !serv.InCookies(cookie.Value) {
			http.Error(w, "Failed UserID", http.StatusUnauthorized)
			return
		}
		cookVal = cookie.Value
	}

	var body []byte
	if r.Header.Get("Content-Type") == "application/x-gzip" {
		body, err = tools.Unpacking(r.Body)
		if err != nil {
			http.Error(w, "Something bad with compression", http.StatusBadRequest)
			return
		}
	} else {
		body, _ = io.ReadAll(r.Body)
		if _, err = url.ParseRequestURI(string(body)); err != nil {
			http.Error(w, "There is not url", http.StatusBadRequest)
			return
		}
	}
	defer r.Body.Close()

	address, flag, err := serv.MakeShortURL(string(body))
	if err != nil {
		http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
		return
	}

	if flag {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
		if err = serv.Save(address, string(body), cookVal); err != nil {
			http.Error(w, "Something bad with Save", http.StatusBadRequest)
			return
		}
	}

	address = serv.GetConfigParameter("baseURL") + "/" + address
	if _, err := w.Write([]byte(address)); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}

func (serv *HStorage) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
	}

	if result, err := serv.Get(r.URL.Path[1:], "key", ""); err == nil && result != "" {
		w.Header().Add("Location", result)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "There is no such shortUrl", http.StatusBadRequest)
		return
	}
}

func (serv *HStorage) JSONPostHandler(w http.ResponseWriter, r *http.Request) {
	type Decoder struct {
		LongURL string `json:"url"`
	}
	type Encoder struct {
		ShortURL string `json:"result"`
	}
	var enc Encoder
	var dec Decoder
	var buf bytes.Buffer
	var err error

	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	var cookVal string
	cookie, err := r.Cookie("UserID")
	if err != nil {
		gen, err1 := tools.GenerateCookie(serv.CountCookies() + 1)
		if err1 != nil {
			http.Error(w, "Something bad with cookies", http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: gen,
		})
		serv.SetCookie(gen)
		cookVal = gen
		fmt.Println(gen)
	} else {
		if !serv.InCookies(cookie.Value) {
			http.Error(w, "Failed UserID", http.StatusUnauthorized)
			return
		}
		cookVal = cookie.Value
	}

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

	var flag bool
	enc.ShortURL, flag, err = serv.MakeShortURL(dec.LongURL)
	if err != nil {
		http.Error(w, "Something bad with MakeShortURL", http.StatusBadRequest)
		return
	}
	if flag {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
		if err = serv.Save(enc.ShortURL, dec.LongURL, cookVal); err != nil {
			http.Error(w, "Something bad with Save", http.StatusBadRequest)
			return
		}
	}

	enc.ShortURL = serv.GetConfigParameter("baseURL") + "/" + enc.ShortURL

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

func (serv *HStorage) GetPingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
	}

	db, err := sql.Open("pgx", serv.GetConfigParameter("db"))
	if err != nil {
		http.Error(w, "Something bad with open db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		http.Error(w, "Something bad with ping", http.StatusInternalServerError)
		return
	}
}

func (serv *HStorage) BatchPostHandler(w http.ResponseWriter, r *http.Request) {
	type Decoder struct {
		ID      string `json:"correlation_id"`
		LongURL string `json:"original_url"`
	}
	type Encoder struct {
		ID       string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}
	var enc []Encoder
	var dec []Decoder
	var buf bytes.Buffer

	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	var cookVal string
	cookie, err := r.Cookie("UserID")
	if err != nil {
		gen, err1 := tools.GenerateCookie(serv.CountCookies() + 1)
		if err1 != nil {
			http.Error(w, "Something bad with cookies", http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: gen,
		})
		serv.SetCookie(gen)
		cookVal = gen
		fmt.Println(gen)
	} else {
		if !serv.InCookies(cookie.Value) {
			http.Error(w, "Failed UserID", http.StatusUnauthorized)
			return
		}
		cookVal = cookie.Value
	}

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
			if err = serv.Save(short, elem.LongURL, cookVal); err != nil {
				http.Error(w, "Something bad with Save", http.StatusBadRequest)
				return
			}
		}

		short = serv.GetConfigParameter("baseURL") + "/" + short
		enc = append(enc, Encoder{ID: elem.ID, ShortURL: short})
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

func (serv *HStorage) GetURLS(w http.ResponseWriter, r *http.Request) {
	type JSONstr struct {
		Key   string `json:"short_url"`
		Value string `json:"original_url"`
	}

	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("UserID")
	if err != nil {
		gen, err1 := tools.GenerateCookie(serv.CountCookies() + 1)
		if err1 != nil {
			http.Error(w, "Something bad with cookies", http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: gen,
		})
		serv.SetCookie(gen)
		http.Error(w, "Failed UserID", http.StatusUnauthorized)
		return
	}
	cookVal := cookie.Value

	mystr, err := serv.Get("", "all", cookVal)
	if err != nil {
		http.Error(w, "Something bad with GetAllURLS", http.StatusBadRequest)
		return
	}
	if mystr == "" {
		http.Error(w, "No data", http.StatusNoContent)
		return
	}

	myarr := strings.Split(mystr, "*")
	var jo []JSONstr
	for i := 0; i < len(myarr); i += 2 {
		jo = append(jo, JSONstr{Key: serv.GetConfigParameter("baseURL") + "/" + myarr[i], Value: myarr[i+1]})
	}

	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(jo)
	if err != nil {
		http.Error(w, "Something bad with encoding JSON", http.StatusBadRequest)
		return
	}

	if _, err = w.Write(resp); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}
