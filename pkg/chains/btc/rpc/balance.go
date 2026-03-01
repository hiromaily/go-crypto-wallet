package rpc

import (
	"encoding/json"
	"fmt"
)

// GetBalance calls getbalance with wildcard account and the given minimum confirmations.
// Returns the balance in BTC as a float64.
func (c *rpcClient) GetBalance(confirmationBlock int) (float64, error) {
	bWildcard, err := json.Marshal("*")
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marshal(wildcard): %w", err)
	}
	bConf, err := json.Marshal(confirmationBlock)
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marshal(confirmationBlock): %w", err)
	}
	rawResult, err := c.client.RawRequest("getbalance", []json.RawMessage{bWildcard, bConf})
	if err != nil {
		return 0, fmt.Errorf("fail to call RawRequest(getbalance): %w", err)
	}
	var amount float64
	if err := json.Unmarshal(rawResult, &amount); err != nil {
		return 0, fmt.Errorf("fail to call json.Unmarshal(getbalance): %w", err)
	}
	return amount, nil
}
