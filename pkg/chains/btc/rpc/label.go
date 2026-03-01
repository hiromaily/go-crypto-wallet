package rpc

import (
	"encoding/json"
	"fmt"
)

// SetLabel calls setlabel to assign a label to an address.
func (c *rpcClient) SetLabel(addr, label string) error {
	bAddr, err := json.Marshal(addr)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(addr): %w", err)
	}
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(label): %w", err)
	}
	_, err = c.client.RawRequest("setlabel", []json.RawMessage{bAddr, bLabel})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(setlabel): %w", err)
	}
	return nil
}
