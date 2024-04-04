package models

type DeleteValue struct {
	Value  string
	Userid string
}

type MSItem struct {
	Key     string
	Value   string
	UserId  string
	Deleted bool
}

type BatchDecoder struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}
type BatchEncoder struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type JsonAllEncoder struct {
	Key   string `json:"short_url"`
	Value string `json:"original_url"`
}

type JsonDecoder struct {
	LongURL string `json:"url"`
}
type JsonEncoder struct {
	ShortURL string `json:"result"`
}

type JSONData struct {
	UUID    int    `json:"uuid"`
	Key     string `json:"short_url"`
	Value   string `json:"original_url"`
	UserID  string `json:"user_id"`
	Deleted bool   `json:"deleted"`
}
