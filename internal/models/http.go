package models

// APIRequest - тело запроса при добавлении URL пользователя.
type APIRequest struct {
	URL string `json:"url"`
}

// APIResponse - тело ответа сервиса.
type APIResponse struct {
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"-"`
}

// APIBatchRequestEl - тело запроса при массовом добавлении URL пользователя.
type APIBatchRequestEl struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// APIBatchResponseEl - тело ответа при массовом добавлении URL пользователя.
type APIBatchResponseEl struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
