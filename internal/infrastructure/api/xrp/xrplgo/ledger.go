package xrplgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	xrpl "github.com/xrpscan/xrpl-go"
)

// ledgerClosedMessage represents a ledger stream message.
type ledgerClosedMessage struct {
	Type             string `json:"type"`
	LedgerIndex      uint64 `json:"ledger_index"`
	LedgerHash       string `json:"ledger_hash"`
	LedgerTime       uint64 `json:"ledger_time"`
	TxnCount         uint64 `json:"txn_count"`
	ValidatedLedgers string `json:"validated_ledgers"`
}

// ledgerResponse represents the response from ledger command.
type ledgerResponse struct {
	LedgerIndex uint64 `json:"ledger_index"`
	LedgerHash  string `json:"ledger_hash"`
	Validated   bool   `json:"validated"`
}

// WaitValidation waits until the specified ledger version is validated.
// Implements LedgerWaiter interface.
func (c *Client) WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error) {
	// First, check current validated ledger
	currentLedger, err := c.getCurrentLedger(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current ledger: %w", err)
	}

	if currentLedger >= targetLedgerVersion {
		return currentLedger, nil
	}

	// Subscribe to ledger stream and wait for target version
	return c.waitForLedger(ctx, targetLedgerVersion)
}

// getCurrentLedger gets the current validated ledger index.
func (c *Client) getCurrentLedger(ctx context.Context) (uint64, error) {
	req := xrpl.BaseRequest{
		"command":      "ledger",
		"ledger_index": "validated",
	}

	resp, err := c.request(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to get ledger: %w", err)
	}

	var result struct {
		Ledger ledgerResponse `json:"ledger"`
	}
	if err := extractResult(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to parse ledger response: %w", err)
	}

	return result.Ledger.LedgerIndex, nil
}

// waitForLedger subscribes to the ledger stream and waits for the target version.
func (c *Client) waitForLedger(ctx context.Context, targetLedgerVersion uint64) (uint64, error) {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil {
		return 0, errors.New("client is closed")
	}

	// Subscribe to ledger stream
	_, err := client.Subscribe([]string{xrpl.StreamTypeLedger})
	if err != nil {
		return 0, fmt.Errorf("failed to subscribe to ledger stream: %w", err)
	}
	defer func() {
		_, _ = client.Unsubscribe([]string{xrpl.StreamTypeLedger})
	}()

	// Wait for target ledger version with timeout
	timeout := 5 * time.Minute
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()

		case ledgerData, ok := <-client.StreamLedger:
			if !ok {
				return 0, errors.New("ledger stream closed unexpectedly")
			}
			var msg ledgerClosedMessage
			if err := json.Unmarshal(ledgerData, &msg); err != nil {
				continue // Skip malformed messages
			}

			if msg.LedgerIndex >= targetLedgerVersion {
				return msg.LedgerIndex, nil
			}

		case <-timer.C:
			// Timeout - try to get current ledger one more time
			currentLedger, err := c.getCurrentLedger(ctx)
			if err != nil {
				return 0, fmt.Errorf("timeout waiting for ledger %d: %w", targetLedgerVersion, err)
			}
			if currentLedger >= targetLedgerVersion {
				return currentLedger, nil
			}
			return 0, fmt.Errorf("timeout waiting for ledger %d, current: %d", targetLedgerVersion, currentLedger)
		}
	}
}

// SubscribeLedger subscribes to the ledger stream and returns a channel that receives ledger indices.
// The caller is responsible for closing the context to stop the subscription.
func (c *Client) SubscribeLedger(ctx context.Context) (<-chan uint64, error) {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()

	if client == nil {
		return nil, errors.New("client is closed")
	}

	// Subscribe to ledger stream
	_, err := client.Subscribe([]string{xrpl.StreamTypeLedger})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to ledger stream: %w", err)
	}

	ch := make(chan uint64, 10)

	go func() {
		defer close(ch)
		defer func() {
			_, _ = client.Unsubscribe([]string{xrpl.StreamTypeLedger})
		}()

		for {
			select {
			case <-ctx.Done():
				return

			case ledgerData, ok := <-client.StreamLedger:
				if !ok {
					// Stream closed, exit goroutine
					return
				}
				var msg ledgerClosedMessage
				if err := json.Unmarshal(ledgerData, &msg); err != nil {
					continue
				}

				select {
				case ch <- msg.LedgerIndex:
				case <-ctx.Done():
					return
				default:
					// Channel full, skip this ledger
				}
			}
		}
	}()

	return ch, nil
}
