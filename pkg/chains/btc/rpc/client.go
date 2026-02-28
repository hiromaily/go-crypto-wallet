package rpc

import "encoding/json"

// RPCCaller abstracts the raw JSON-RPC call method.
// The concrete *rpcclient.Client satisfies this interface without modification.
type RPCCaller interface {
	RawRequest(method string, params []json.RawMessage) (json.RawMessage, error)
}
