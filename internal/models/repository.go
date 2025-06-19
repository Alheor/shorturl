package models

// BatchEl - элемент сокращенного URL при обработке массовой вставки.
type BatchEl struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
}

// UniqueErr - тип ошибки, обозначающий, что вставляем URL уже существует.
type UniqueErr struct {
	ShortKey string
	Err      error
}

func (e *UniqueErr) Error() string {
	return e.Err.Error()
}

// HistoryNotFoundErr - тип ошибки, обозначающий, что URL-ы отсутствуют. Используется при запросе всех URL пользователя.
type HistoryNotFoundErr struct {
	error
}

// HistoryEl - сокращенный url, используется при запросе всех URL пользователя.
type HistoryEl struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
