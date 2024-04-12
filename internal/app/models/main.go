// Package models - пакет с различными моделями, которые импользуются в разных пакетах
package models

// DeleteValue  - структура для хранения данных, которые отправляются в канал.
type DeleteValue struct {
	Value  string
	Userid string
}

// MSItem - структура хранения данных, если для хранилища используется память "компьютера".
type MSItem struct {
	Key     string
	Value   string
	UserID  string
	Deleted bool
}

// BatchDecoder - структура для возвращаемых данных из хендлера BatchURLs.
type BatchDecoder struct {
	ID      string `json:"correlation_id"`
	LongURL string `json:"original_url"`
}

type BatchDecoderArray struct {
	Items []BatchDecoder
}

// BatchEncoder - структура для получаемых данных в хендлере BatchURLs.
type BatchEncoder struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type BatchEncoderArray struct {
	Items []BatchEncoder
}

// JSONAllEncoder - структура для возвращаемых данных из хендлера GetURLs.
type JSONAllEncoder struct {
	Key   string `json:"short_url"`
	Value string `json:"original_url"`
}

// JSONDecoder - структура для получаемых даных в хендлере SetJSONLongURL.
type JSONDecoder struct {
	LongURL string `json:"url"`
}

// JSONEncoder - структура для возвращаемых данных из хендлера SetJSONLongURL.
type JSONEncoder struct {
	ShortURL string `json:"result"`
}

// JSONData - структура, которая используется для получения или записи данных в файл-хранилище.
type JSONData struct {
	UUID    int    `json:"uuid"`
	Key     string `json:"short_url"`
	Value   string `json:"original_url"`
	UserID  string `json:"user_id"`
	Deleted bool   `json:"deleted"`
}
