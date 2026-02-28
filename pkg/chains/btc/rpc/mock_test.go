package rpc

import "encoding/json"

// mockRPCCaller is a test double for RPCCaller shared across all test files.
type mockRPCCaller struct {
	rawRequestFn func(method string, params []json.RawMessage) (json.RawMessage, error)
}

func (m *mockRPCCaller) RawRequest(method string, params []json.RawMessage) (json.RawMessage, error) {
	return m.rawRequestFn(method, params)
}
