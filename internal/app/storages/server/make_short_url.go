package server

import "github.com/levshindenis/sprint1/internal/app/tools"

func (serv *Server) MakeShortURL(longURL string) (string, bool, error) {
	value, _, err := serv.st.GetData(longURL, "Value", "")
	if err != nil {
		return "", false, err
	}
	if value != "" {
		return value, true, nil
	}
	shortKey := tools.GenerateShortKey()
	for {
		result, _, err := serv.GetStorage().GetData(shortKey, "key", "")
		if err != nil {
			return "", false, err
		}
		if result == "" {
			return shortKey, false, nil
		}
		shortKey = tools.GenerateShortKey()
	}
}
