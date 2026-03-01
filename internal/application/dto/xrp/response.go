package xrp

// ResponseGenerateAddress represents the response from GenerateAddress API call.
type ResponseGenerateAddress struct {
	XAddress       string
	ClassicAddress string
	Address        string
	Secret         string
}

// ResponseGenerateXAddress represents the response from GenerateXAddress API call.
type ResponseGenerateXAddress struct {
	XAddress string
	Secret   string
}
