package httphandler

type APIResponse struct {
	Result     string `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"-"`
}
