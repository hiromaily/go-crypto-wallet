package xrplgo

import (
	"context"
	"testing"
	"time"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// Ensure mockLogger satisfies logger.Logger at compile time.
var _ logger.Logger = (*mockLogger)(nil)

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
				ValidationTimeout: 10 * time.Minute,
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
	if config.ValidationTimeout != 5*time.Minute {
		t.Errorf("DefaultConfig().ValidationTimeout = %v, want %v", config.ValidationTimeout, 5*time.Minute)
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
					"error_code":    float64(19),
					"error_message": "Account not found.",
				},
			},
			wantErr: true,
		},
		{
			name: "error with missing error_message",
			resp: map[string]any{
				"result": map[string]any{
					"error":      "actNotFound",
					"error_code": float64(19),
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

func TestParseAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   any
		wantInfo amountInfo
	}{
		{
			name:   "XRP amount as string",
			amount: "1000000",
			wantInfo: amountInfo{
				Value:    "1",
				Currency: "XRP",
			},
		},
		{
			name: "token amount as object",
			amount: map[string]any{
				"value":    "100.5",
				"currency": "USD",
				"issuer":   "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			},
			wantInfo: amountInfo{
				Value:    "100.5",
				Currency: "USD",
				Issuer:   "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			},
		},
		{
			name:   "nil amount",
			amount: nil,
			wantInfo: amountInfo{
				Value:    "0",
				Currency: "XRP",
			},
		},
		{
			name:   "unexpected type",
			amount: 12345,
			wantInfo: amountInfo{
				Value:    "0",
				Currency: "XRP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseAmount(tt.amount)
			if got.Value != tt.wantInfo.Value {
				t.Errorf("parseAmount().Value = %v, want %v", got.Value, tt.wantInfo.Value)
			}
			if got.Currency != tt.wantInfo.Currency {
				t.Errorf("parseAmount().Currency = %v, want %v", got.Currency, tt.wantInfo.Currency)
			}
			if got.Issuer != tt.wantInfo.Issuer {
				t.Errorf("parseAmount().Issuer = %v, want %v", got.Issuer, tt.wantInfo.Issuer)
			}
		})
	}
}

// mockLogger implements logger.Logger interface for testing.
type mockLogger struct {
	debugCalls []string
	warnCalls  []string
	errorCalls []string
}

func (m *mockLogger) Debug(msg string, _ ...any) { m.debugCalls = append(m.debugCalls, msg) }
func (*mockLogger) Info(_ string, _ ...any)      {}
func (m *mockLogger) Warn(msg string, _ ...any)  { m.warnCalls = append(m.warnCalls, msg) }
func (m *mockLogger) Error(msg string, _ ...any) { m.errorCalls = append(m.errorCalls, msg) }

func TestClientWithGlobalLogger(t *testing.T) {
	t.Parallel()

	// Set a mock global logger and verify NewClient succeeds.
	logger.SetGlobal(&mockLogger{})
	t.Cleanup(func() { logger.SetGlobal(logger.NewNoopLogger()) })

	config := ClientConfig{
		WebSocketURL: "wss://s.altnet.rippletest.net:51233",
	}
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer func() { _ = client.Close() }()
}

func TestNoopLogger(t *testing.T) {
	t.Parallel()

	// Ensure NoopLogger doesn't panic
	l := logger.NewNoopLogger()
	l.Debug("test")
	l.Warn("test")
	l.Error("test")
}

func TestDropsInXRPConstant(t *testing.T) {
	t.Parallel()

	if dropsInXRP != 1_000_000 {
		t.Errorf("dropsInXRP = %d, want %d", dropsInXRP, 1_000_000)
	}
}
