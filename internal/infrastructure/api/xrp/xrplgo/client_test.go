package xrplgo

import (
	"context"
	"testing"
	"time"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
)

// TestClientImplementsInterfaces verifies that Client implements the required interfaces.
func TestClientImplementsInterfaces(t *testing.T) {
	t.Parallel()

	// This test verifies at compile time that Client implements the required interfaces.
	// If the interfaces are not implemented, this test will fail to compile.
	var _ apixrp.AccountInfoProvider = (*Client)(nil)
	var _ apixrp.TransactionSubmitter = (*Client)(nil)
	var _ apixrp.TransactionGetter = (*Client)(nil)
	var _ apixrp.LedgerWaiter = (*Client)(nil)
	var _ apixrp.BalanceChecker = (*Client)(nil)
	var _ apixrp.Closer = (*Client)(nil)
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ClientConfig{
				WebSocketURL: "wss://s.altnet.rippletest.net:51233",
			},
			wantErr: false,
		},
		{
			name:    "empty URL",
			config:  ClientConfig{},
			wantErr: true,
		},
		{
			name: "with timeouts",
			config: ClientConfig{
				WebSocketURL:      "wss://s.altnet.rippletest.net:51233",
				ReadTimeout:       30 * time.Second,
				WriteTimeout:      30 * time.Second,
				HeartbeatInterval: 10 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
			if client != nil {
				_ = client.Close()
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	url := "wss://s.altnet.rippletest.net:51233"
	config := DefaultConfig(url)

	if config.WebSocketURL != url {
		t.Errorf("DefaultConfig().WebSocketURL = %v, want %v", config.WebSocketURL, url)
	}
	if config.ReadTimeout != 60*time.Second {
		t.Errorf("DefaultConfig().ReadTimeout = %v, want %v", config.ReadTimeout, 60*time.Second)
	}
	if config.WriteTimeout != 60*time.Second {
		t.Errorf("DefaultConfig().WriteTimeout = %v, want %v", config.WriteTimeout, 60*time.Second)
	}
	if config.HeartbeatInterval != 5*time.Second {
		t.Errorf("DefaultConfig().HeartbeatInterval = %v, want %v", config.HeartbeatInterval, 5*time.Second)
	}
}

func TestDropsToXRP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		drops string
		want  string
	}{
		{
			name:  "1 XRP",
			drops: "1000000",
			want:  "1",
		},
		{
			name:  "0.1 XRP",
			drops: "100000",
			want:  "0.1",
		},
		{
			name:  "0.000001 XRP",
			drops: "1",
			want:  "0.000001",
		},
		{
			name:  "zero",
			drops: "0",
			want:  "0",
		},
		{
			name:  "large amount",
			drops: "100000000000",
			want:  "100000",
		},
		{
			name:  "invalid",
			drops: "invalid",
			want:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := dropsToXRP(tt.drops)
			if got != tt.want {
				t.Errorf("dropsToXRP(%q) = %q, want %q", tt.drops, got, tt.want)
			}
		})
	}
}

func TestExtractResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		resp    map[string]any
		wantErr bool
	}{
		{
			name: "valid result",
			resp: map[string]any{
				"result": map[string]any{
					"status": "success",
				},
			},
			wantErr: false,
		},
		{
			name:    "missing result",
			resp:    map[string]any{},
			wantErr: true,
		},
		{
			name: "error in result",
			resp: map[string]any{
				"result": map[string]any{
					"error":         "actNotFound",
					"error_code":    19,
					"error_message": "Account not found.",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var target map[string]any
			err := extractResult(tt.resp, &target)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientClose(t *testing.T) {
	t.Parallel()

	config := ClientConfig{
		WebSocketURL: "wss://s.altnet.rippletest.net:51233",
	}
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// First close should succeed
	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Subsequent operations should fail
	_, err = client.GetAccountInfo(context.Background(), "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")
	if err == nil {
		t.Error("GetAccountInfo() after Close() should return error")
	}
}
