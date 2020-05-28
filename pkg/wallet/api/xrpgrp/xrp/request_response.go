package xrp

// RequestCommand is only command request
type RequestCommand struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// ResponseError is common error
type ResponseError struct {
	ID      int         `json:"id"`
	Status  string      `json:"status"`
	Type    string      `json:"type"`
	Error   string      `json:"error"`
	Request interface{} `json:"request"`
}

// StatusCode is status code of response
type StatusCode string

// status_code
const (
	StatusCodeError   StatusCode = "error"
	StatusCodeSuccess StatusCode = "success"
)

// String converter
func (s StatusCode) String() string {
	return string(s)
}
