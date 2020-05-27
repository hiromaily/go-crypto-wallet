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
