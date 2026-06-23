package transport

// ErrorResp is the standard JSON error response.
type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
