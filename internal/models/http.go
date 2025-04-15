package models

type APIRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"-"`
}

type APIBatchRequestEl struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type APIBatchResponseEl struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
