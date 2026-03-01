package rpc

import (
	"encoding/json"
	"fmt"
)

// ImportAddress calls importaddress to import an address or script.
func (c *Client) ImportAddress(address, label string, rescan bool) error {
	bAddress, err := json.Marshal(address)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(address): %w", err)
	}
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(label): %w", err)
	}
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(rescan): %w", err)
	}
	_, err = c.client.RawRequest("importaddress", []json.RawMessage{bAddress, bLabel, bRescan})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(importaddress): %w", err)
	}
	return nil
}

// ImportPrivKey calls importprivkey to import a WIF-encoded private key.
func (c *Client) ImportPrivKey(wifKey, label string, rescan bool) error {
	bKey, err := json.Marshal(wifKey)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(wifKey): %w", err)
	}
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(label): %w", err)
	}
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return fmt.Errorf("fail to call json.Marshal(rescan): %w", err)
	}
	_, err = c.client.RawRequest("importprivkey", []json.RawMessage{bKey, bLabel, bRescan})
	if err != nil {
		return fmt.Errorf("fail to call RawRequest(importprivkey): %w", err)
	}
	return nil
}
