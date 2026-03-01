package rpc

import (
	"encoding/json"
	"fmt"
)

// EstimateSmartFeeResult is the wire-format response of the estimatesmartfee command.
type EstimateSmartFeeResult struct {
	FeeRate float64  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  uint64   `json:"blocks"`
}

// EstimateSmartFee calls estimatesmartfee with the given confirmation target
// and returns the fee rate in BTC/kB.
func (c *Client) EstimateSmartFee(confirmationBlock int) (float64, error) {
	input, err := json.Marshal(confirmationBlock)
	if err != nil {
		return 0, fmt.Errorf("fail to call json.Marshal(confirmationBlock): %w", err)
	}
	rawResult, err := c.client.RawRequest("estimatesmartfee", []json.RawMessage{input})
	if err != nil {
		return 0, fmt.Errorf("fail to call RawRequest(estimatesmartfee): %w", err)
	}
	var result EstimateSmartFeeResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return 0, fmt.Errorf("fail to call json.Unmarshal(estimatesmartfee): %w", err)
	}
	if len(result.Errors) != 0 {
		return 0, fmt.Errorf("estimatesmartfee response includes error: %s", result.Errors[0])
	}
	return result.FeeRate, nil
}
