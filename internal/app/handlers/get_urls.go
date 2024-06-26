package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/levshindenis/sprint1/internal/app/models"
)

// GetURLs нужен для обработки запроса от клиента по адресу /api/user/urls.
// Берется куки клиента и из хранилища достаются все сокращенные URL этого клиента.
// Если клиент не сократил ни одного URL, то возвращается StatusNoContent.
// Если сокращенные URL есть, то данные преобращуются в JSON и возвращаются клиенту.
func (serv *HStorage) GetURLs(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("UserID")

	mystr, _, err := serv.GetStorage().GetData("", "all", cookie.Value)
	if err != nil {
		http.Error(w, "Something bad with GetAllURLS", http.StatusBadRequest)
		return
	}
	if mystr == "" {
		http.Error(w, "No data", http.StatusNoContent)
		return
	}

	myarr := strings.Split(mystr[:len(mystr)-1], "*")

	jo := make([]models.JSONAllEncoder, len(myarr)/2)
	for i := 0; i < len(myarr); i += 2 {

		jo[i/2] = models.JSONAllEncoder{Key: serv.GetServerConfig().GetShortBaseURL() + "/" + myarr[i],
			Value: myarr[i+1]}

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
