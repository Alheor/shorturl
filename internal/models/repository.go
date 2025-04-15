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

func (e *UniqueErr) Error() string {
	return e.Err.Error()
}
