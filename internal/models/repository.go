package models

type BatchEl struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
}

type UniqueErr struct {
	ShortKey string
	Err      error
}

type HistoryNotFoundErr struct {
	error
}

type HistoryEl struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func (e *UniqueErr) Error() string {
	return e.Err.Error()
}
