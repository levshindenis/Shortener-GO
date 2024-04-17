package handlers

import (
	"net/http"
)

// GetLongURL нужен для обработки запроса от клиента по адресу /{id}, где id - сокращенный URL.
// Значение короткого URL берется из последней части URL.
// Сначала проверяется, есть ли такой короткий URL в хранилище. Если нет, то возвращается ошибка.
// Затем идет проверка на удаление URL. Если URL удален, то устанавливается StatusGone.
// При успешной обработке запроса в заголовке Location устанавливается длинный URL и
// устанавливется StatusTemporaryRedirect.
func (serv *HStorage) GetLongURL(w http.ResponseWriter, r *http.Request) {
	result, isdeleted, err := serv.GetStorage().GetData(r.URL.Path[1:], "key", "")
	if err != nil {
		http.Error(w, "Something bad with GetLongURL", http.StatusBadRequest)
		return
	}
	if result == "" {
		http.Error(w, "There is no such shortUrl", http.StatusBadRequest)
		return
	}
	if isdeleted[0] {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Add("Location", result)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
