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

	cookVal, cookFlag, err := serv.GetCookieStorage().CheckCookie(r)
	if err != nil {
		http.Error(w, "Something bad with check cookie", http.StatusBadRequest)
		return
	}
	if !cookFlag {
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: cookVal,
		})
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
		if err = serv.GetStorageData().SetData(address, string(body), cookVal); err != nil {
			http.Error(w, "Something bad with Save", http.StatusBadRequest)
			return
		}
	}

	address = serv.GetServerConfig("baseURL") + "/" + address
	if _, err := w.Write([]byte(address)); err != nil {
		http.Error(w, "Something bad with write address", http.StatusBadRequest)
		return
	}
}

func (serv *HStorage) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "There is not true method", http.StatusBadRequest)
	}

	result, isdeleted, err := serv.GetStorageData().GetData(r.URL.Path[1:], "key", "")
	if err != nil {
		http.Error(w, "Something bad with GetHandler", http.StatusBadRequest)
		return
	} else if isdeleted[0] {
		w.WriteHeader(http.StatusGone)
	} else if result != "" {
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

	if r.Method != http.MethodPost {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	cookVal, cookFlag, err := serv.GetCookieStorage().CheckCookie(r)
	if err != nil {
		http.Error(w, "Something bad with check cookie", http.StatusBadRequest)
		return
	}
	if !cookFlag {
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: cookVal,
		})
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
		if err = serv.GetStorageData().SetData(enc.ShortURL, dec.LongURL, cookVal); err != nil {
			http.Error(w, "Something bad with Save", http.StatusBadRequest)
			return
		}
	}

	enc.ShortURL = serv.GetServerConfig("baseURL") + "/" + enc.ShortURL

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

	db, err := sql.Open("pgx", serv.GetServerConfig("db"))
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

	cookVal, cookFlag, err := serv.GetCookieStorage().CheckCookie(r)
	if err != nil {
		http.Error(w, "Something bad with check cookie", http.StatusBadRequest)
		return
	}
	if !cookFlag {
		http.SetCookie(w, &http.Cookie{
			Name:  "UserID",
			Value: cookVal,
		})
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
			if err = serv.GetStorageData().SetData(short, elem.LongURL, cookVal); err != nil {
				http.Error(w, "Something bad with Save", http.StatusBadRequest)
				return
			}
		}

		short = serv.GetServerConfig("baseURL") + "/" + short
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

	cookVal, myflag, err := serv.GetCookieStorage().CheckCookie(r)
	fmt.Println(cookVal)
	if err != nil {
		http.Error(w, "Something bad with check cookie", http.StatusBadRequest)
		return
	}
	if !myflag {
		http.Error(w, "Failed UserID", http.StatusUnauthorized)
		return
	}

	mystr, _, err := serv.GetStorageData().GetData("", "all", cookVal)
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
		jo = append(jo, JSONstr{Key: serv.GetServerConfig("baseURL") + "/" + myarr[i],
			Value: myarr[i+1]})
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

func (serv *HStorage) DelURLS(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "There is not true method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "There is incorrect data format", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	var shortURLS []string

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

	cookie, err := r.Cookie("UserID")
	if err != nil {
		http.Error(w, "Something bad with cookie", http.StatusBadRequest)
		return
	}

	for _, elem := range shortURLS {
		serv.SetChan(storages.DeleteValue{Value: elem, Userid: cookie.Value})
	}

	w.WriteHeader(http.StatusAccepted)
}
