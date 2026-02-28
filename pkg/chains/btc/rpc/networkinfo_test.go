package rpc

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── GetNetworkInfo ───────────────────────────────────────────────────────────

func TestGetNetworkInfo_HappyPath(t *testing.T) {
	t.Parallel()
	resp := GetNetworkInfoResult{
		Version:    280000,
		Subversion: "/Satoshi:28.0.0/",
		RelayFee:   0.00001,
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getnetworkinfo", method)
			return b, nil
		},
	}
	res, err := GetNetworkInfo(caller)
	require.NoError(t, err)
	assert.Equal(t, int32(280000), res.Version)
	assert.InDelta(t, 0.00001, res.RelayFee, 1e-10)
}

func TestGetNetworkInfo_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("node error")
		},
	}
	_, err := GetNetworkInfo(caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getnetworkinfo")
}

// ─── GetBlockchainInfo ────────────────────────────────────────────────────────

func TestGetBlockchainInfo_HappyPath(t *testing.T) {
	t.Parallel()
	resp := GetBlockchainInfoResult{
		Chain:  "main",
		Blocks: 800000,
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getblockchaininfo", method)
			return b, nil
		},
	}
	res, err := GetBlockchainInfo(caller)
	require.NoError(t, err)
	assert.Equal(t, "main", res.Chain)
	assert.Equal(t, uint32(800000), res.Blocks)
}

func TestGetBlockchainInfo_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("node error")
		},
	}
	_, err := GetBlockchainInfo(caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getblockchaininfo")
}

// ─── GetBlockCount ────────────────────────────────────────────────────────────

func TestGetBlockCount_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(int64(800000))
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getblockcount", method)
			return b, nil
		},
	}
	count, err := GetBlockCount(caller)
	require.NoError(t, err)
	assert.Equal(t, int64(800000), count)
}

func TestGetBlockCount_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("node error")
		},
	}
	_, err := GetBlockCount(caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getblockcount")
}
