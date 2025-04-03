package models

type APIRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"-"`
}
