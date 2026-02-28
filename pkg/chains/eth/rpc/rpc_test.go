package rpc

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRPCCaller is a simple test double for RPCCaller.
// callContextFn receives the same arguments as CallContext and is responsible for
// populating *result via a type assertion.
type mockRPCCaller struct {
	callContextFn func(ctx context.Context, result any, method string, args ...any) error
}

func (m *mockRPCCaller) CallContext(ctx context.Context, result any, method string, args ...any) error {
	return m.callContextFn(ctx, result, method, args...)
}

// setString is a helper that assigns a string value through the result pointer.
func setString(result any, value string) {
	if p, ok := result.(*string); ok {
		*p = value
	}
}

// setBool is a helper that assigns a bool value through the result pointer.
func setBool(result any, value bool) {
	if p, ok := result.(*any); ok {
		*p = value
	}
}

// setMap is a helper that assigns a map[string]any value through the result pointer.
func setMapAny(result any, value map[string]any) {
	if p, ok := result.(*any); ok {
		*p = value
	}
}

// setStringSlice is a helper that assigns a []string value through the result pointer.
func setStringSlice(result any, value []string) {
	if p, ok := result.(*[]string); ok {
		*p = value
	}
}

// ─── Syncing ──────────────────────────────────────────────────────────────────

func TestSyncing_NotSyncing(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_syncing", method)
			setBool(result, false)
			return nil
		},
	}
	syncing, isSyncing, err := Syncing(context.Background(), caller)
	require.NoError(t, err)
	assert.Nil(t, syncing)
	assert.False(t, isSyncing)
}

func TestSyncing_IsSyncing(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_syncing", method)
			setMapAny(result, map[string]any{
				"currentBlock":  "0x1",
				"highestBlock":  "0xa",
				"startingBlock": "0x0",
			})
			return nil
		},
	}
	syncing, isSyncing, err := Syncing(context.Background(), caller)
	require.NoError(t, err)
	require.NotNil(t, syncing)
	assert.True(t, isSyncing)
	assert.Equal(t, int64(1), syncing.CurrentBlock)
	assert.Equal(t, int64(10), syncing.HighestBlock)
}

func TestSyncing_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("connection refused")
		},
	}
	_, _, err := Syncing(context.Background(), caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "eth_syncing")
}

// ─── BlockNumber ──────────────────────────────────────────────────────────────

func TestBlockNumber_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_blockNumber", method)
			setString(result, "0x4b7")
			return nil
		},
	}
	num, err := BlockNumber(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0x4b7), num)
}

func TestBlockNumber_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("rpc error")
		},
	}
	_, err := BlockNumber(context.Background(), caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "eth_blockNumber")
}

func TestBlockNumber_InvalidHex(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, _ string, _ ...any) error {
			setString(result, "not-hex")
			return nil
		},
	}
	_, err := BlockNumber(context.Background(), caller)
	require.Error(t, err)
}

// ─── GetBalance ───────────────────────────────────────────────────────────────

func TestGetBalance_HappyPath(t *testing.T) {
	t.Parallel()
	const addr = "0xabcdef1234567890abcdef1234567890abcdef12"
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, args ...any) error {
			require.Equal(t, "eth_getBalance", method)
			require.Equal(t, addr, args[0])
			require.Equal(t, "latest", args[1])
			setString(result, "0xde0b6b3a7640000") // 1 ETH in wei
			return nil
		},
	}
	balance, err := GetBalance(context.Background(), caller, addr, QuantityTagLatest)
	require.NoError(t, err)
	expected, _ := new(big.Int).SetString("de0b6b3a7640000", 16)
	assert.Equal(t, expected, balance)
}

func TestGetBalance_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("node unavailable")
		},
	}
	_, err := GetBalance(context.Background(), caller, "0x0", QuantityTagLatest)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "eth_getBalance")
}

// ─── GetTransactionCount ──────────────────────────────────────────────────────

func TestGetTransactionCount_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_getTransactionCount", method)
			setString(result, "0x5")
			return nil
		},
	}
	count, err := GetTransactionCount(context.Background(), caller, "0x1234", QuantityTagPending)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(5), count)
}

// ─── GasPrice ─────────────────────────────────────────────────────────────────

func TestGasPrice_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_gasPrice", method)
			setString(result, "0x3b9aca00") // 1 Gwei
			return nil
		},
	}
	price, err := GasPrice(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1_000_000_000), price)
}

func TestGasPrice_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("rpc error")
		},
	}
	_, err := GasPrice(context.Background(), caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "eth_gasPrice")
}

// ─── GetTransactionByHash ─────────────────────────────────────────────────────

func TestGetTransactionByHash_HappyPath(t *testing.T) {
	t.Parallel()
	const txHash = "0xabc123"
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, args ...any) error {
			require.Equal(t, "eth_getTransactionByHash", method)
			require.Equal(t, txHash, args[0])
			if p, ok := result.(*map[string]string); ok {
				*p = map[string]string{
					"blockHash":        "0xblockhash",
					"blockNumber":      "0x1",
					"from":             "0xfrom",
					"gas":              "0x5208",
					"gasPrice":         "0x3b9aca00",
					"hash":             txHash,
					"input":            "0x",
					"nonce":            "0x0",
					"to":               "0xto",
					"transactionIndex": "0x0",
					"value":            "0xde0b6b3a7640000",
					"v":                "0x1",
					"r":                "0xr",
					"s":                "0xs",
				}
			}
			return nil
		},
	}
	tx, err := GetTransactionByHash(context.Background(), caller, txHash)
	require.NoError(t, err)
	require.NotNil(t, tx)
	assert.Equal(t, txHash, tx.Hash)
	assert.Equal(t, int64(1), tx.BlockNumber)
	assert.Equal(t, int64(0x5208), tx.Gas)
}

func TestGetTransactionByHash_EmptyResponse(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return nil // result stays zero-value (empty map)
		},
	}
	_, err := GetTransactionByHash(context.Background(), caller, "0xhash")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestGetTransactionByHash_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("network error")
		},
	}
	_, err := GetTransactionByHash(context.Background(), caller, "0xhash")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "eth_getTransactionByHash")
}

// ─── GetTxReceipt (not-found returns nil, nil) ────────────────────────────────

func TestGetTxReceipt_NotFound(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return nil // empty resMap → errReceiptNotFound inside GetTransactionReceipt
		},
	}
	receipt, err := GetTxReceipt(context.Background(), caller, "0xhash")
	require.NoError(t, err)
	assert.Nil(t, receipt)
}

func TestGetTxReceipt_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("connection reset")
		},
	}
	_, err := GetTxReceipt(context.Background(), caller, "0xhash")
	require.Error(t, err)
}

// ─── NetVersion / NetListening / NetPeerCount ─────────────────────────────────

func TestNetVersion_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "net_version", method)
			setString(result, "1")
			return nil
		},
	}
	v, err := NetVersion(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, uint16(1), v)
}

func TestNetVersion_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("rpc error")
		},
	}
	_, err := NetVersion(context.Background(), caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "net_version")
}

func TestNetListening_True(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "net_listening", method)
			if p, ok := result.(*bool); ok {
				*p = true
			}
			return nil
		},
	}
	listening, err := NetListening(context.Background(), caller)
	require.NoError(t, err)
	assert.True(t, listening)
}

func TestNetPeerCount_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "net_peerCount", method)
			setString(result, "0x3")
			return nil
		},
	}
	count, err := NetPeerCount(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(3), count)
}

// ─── ClientVersion / SHA3 ─────────────────────────────────────────────────────

func TestClientVersion_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "web3_clientVersion", method)
			setString(result, "Geth/v1.13.0-stable")
			return nil
		},
	}
	v, err := ClientVersion(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, "Geth/v1.13.0-stable", v)
}

func TestClientVersion_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, _ any, _ string, _ ...any) error {
			return errors.New("rpc error")
		},
	}
	_, err := ClientVersion(context.Background(), caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "web3_clientVersion")
}

func TestSHA3_HappyPath(t *testing.T) {
	t.Parallel()
	const input = "0x68656c6c6f20776f726c64"
	const expected = "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, args ...any) error {
			require.Equal(t, "web3_sha3", method)
			require.Equal(t, input, args[0])
			setString(result, expected)
			return nil
		},
	}
	hash, err := SHA3(context.Background(), caller, input)
	require.NoError(t, err)
	assert.Equal(t, expected, hash)
}

// ─── Coinbase / Accounts ──────────────────────────────────────────────────────

func TestCoinbase_HappyPath(t *testing.T) {
	t.Parallel()
	const addr = "0xabcdef1234567890abcdef1234567890abcdef12"
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_coinbase", method)
			setString(result, addr)
			return nil
		},
	}
	got, err := Coinbase(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, addr, got)
}

func TestAccounts_HappyPath(t *testing.T) {
	t.Parallel()
	want := []string{"0xaddr1", "0xaddr2"}
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_accounts", method)
			setStringSlice(result, want)
			return nil
		},
	}
	got, err := Accounts(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// ─── ProtocolVersion ──────────────────────────────────────────────────────────

func TestProtocolVersion_HappyPath(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		callContextFn: func(_ context.Context, result any, method string, _ ...any) error {
			require.Equal(t, "eth_protocolVersion", method)
			setString(result, "0x41") // 65
			return nil
		},
	}
	v, err := ProtocolVersion(context.Background(), caller)
	require.NoError(t, err)
	assert.Equal(t, uint64(0x41), v)
}

// ─── ValidateAddr ─────────────────────────────────────────────────────────────

func TestValidateAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{"valid", "0xabcdef1234567890abcdef1234567890abcdef12", false},
		{"invalid", "notanaddress", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateAddr(tt.addr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ─── QuantityTag ─────────────────────────────────────────────────────────────

func TestQuantityTag_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "latest", QuantityTagLatest.String())
	assert.Equal(t, "pending", QuantityTagPending.String())
}
