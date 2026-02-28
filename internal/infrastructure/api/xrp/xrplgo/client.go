// Package xrplgo provides a wrapper around the xrpl-go library for direct XRPL node communication.
// This package implements the small interfaces defined in application/ports/api/xrp.
package xrplgo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	xrpl "github.com/xrpscan/xrpl-go"
)

// ClientConfig holds configuration for the XRPL client.
type ClientConfig struct {
	// WebSocketURL is the XRPL node WebSocket URL (e.g., "wss://s.altnet.rippletest.net:51233").
	WebSocketURL string

	// ReadTimeout is the timeout for read operations. Default: 60 seconds.
	ReadTimeout time.Duration

	// WriteTimeout is the timeout for write operations. Default: 60 seconds.
	WriteTimeout time.Duration

	// HeartbeatInterval is the interval between heartbeat messages. Default: 5 seconds.
	HeartbeatInterval time.Duration

	// ValidationTimeout is the timeout for waiting ledger validation. Default: 5 minutes.
	ValidationTimeout time.Duration
}

// DefaultConfig returns a ClientConfig with default values.
func DefaultConfig(wsURL string) ClientConfig {
	return ClientConfig{
		WebSocketURL:      wsURL,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		HeartbeatInterval: 5 * time.Second,
		ValidationTimeout: 5 * time.Minute,
	}
}

// Client wraps the xrpl-go client and implements the XRP interfaces.
type Client struct {
	client *xrpl.Client
	config ClientConfig
	mu     sync.RWMutex
}

// NewClient creates a new Client with the given configuration.
func NewClient(config ClientConfig) (*Client, error) {
	if config.WebSocketURL == "" {
		return nil, errors.New("websocket URL is required")
	}

	xrplConfig := xrpl.ClientConfig{
		URL: config.WebSocketURL,
	}

	if config.ReadTimeout > 0 {
		xrplConfig.ReadTimeout = config.ReadTimeout
	}
	if config.WriteTimeout > 0 {
		xrplConfig.WriteTimeout = config.WriteTimeout
	}
	if config.HeartbeatInterval > 0 {
		xrplConfig.HeartbeatInterval = config.HeartbeatInterval
	}

	// Set default validation timeout if not specified
	if config.ValidationTimeout == 0 {
		config.ValidationTimeout = 5 * time.Minute
	}

	client := xrpl.NewClient(xrplConfig)

	return &Client{
		client: client,
		config: config,
	}, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		return err
	}
	return nil
}

// request sends a request to the XRPL node and returns the response.
// The context is currently not used for cancellation but included for future compatibility.
func (c *Client) request(_ context.Context, req xrpl.BaseRequest) (xrpl.BaseResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.client == nil {
		return nil, errors.New("client is closed")
	}

	return c.client.Request(req)
}

// extractResult extracts the "result" field from a response and decodes it into the target.
func extractResult(resp xrpl.BaseResponse, target any) error {
	result, hasResult := resp["result"]
	if !hasResult {
		return errors.New("response missing 'result' field")
	}

	// Check for error in result
	if resultMap, isMap := result.(map[string]any); isMap {
		if errMsg, hasError := resultMap["error"].(string); hasError {
			var errCode float64
			var errMessage string

			if code, ok := resultMap["error_code"].(float64); ok {
				errCode = code
			}
			if msg, ok := resultMap["error_message"].(string); ok {
				errMessage = msg
			}
			if errMessage == "" {
				errMessage = errMsg
			}
			return fmt.Errorf("XRPL error %v: %s", errCode, errMessage)
		}
	}

	// Use mapstructure for efficient decoding from map to struct
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           target,
		TagName:          "json",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(result); err != nil {
		return fmt.Errorf("failed to decode result: %w", err)
	}

	return nil
}
