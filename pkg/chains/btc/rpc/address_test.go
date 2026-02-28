package rpc

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── GetAddressInfo ───────────────────────────────────────────────────────────

func TestGetAddressInfo_HappyPath(t *testing.T) {
	t.Parallel()
	resp := GetAddressInfoResult{
		Address:      "bc1qtest",
		ScriptPubKey: "0014abcd",
		IsMine:       true,
		Labels:       FlexibleLabels{"deposit"},
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, params []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getaddressinfo", method)
			require.Len(t, params, 1)
			var addr string
			_ = json.Unmarshal(params[0], &addr)
			assert.Equal(t, "bc1qtest", addr)
			return b, nil
		},
	}
	res, err := GetAddressInfo(caller, "bc1qtest")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "bc1qtest", res.Address)
	assert.Equal(t, FlexibleLabels{"deposit"}, res.Labels)
}

func TestGetAddressInfo_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("rpc error")
		},
	}
	_, err := GetAddressInfo(caller, "bad_addr")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getaddressinfo")
}

// ─── ValidateAddress ──────────────────────────────────────────────────────────

func TestValidateAddress_HappyPath(t *testing.T) {
	t.Parallel()
	resp := ValidateAddressResult{IsValid: true, Address: "bc1qtest", ScriptPubKey: "0014abcd"}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, params []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "validateaddress", method)
			return b, nil
		},
	}
	res, err := ValidateAddress(caller, "bc1qtest")
	require.NoError(t, err)
	assert.True(t, res.IsValid)
}

func TestValidateAddress_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("rpc error")
		},
	}
	_, err := ValidateAddress(caller, "bad_addr")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validateaddress")
}

// ─── GetAddressesByLabel ──────────────────────────────────────────────────────

func TestGetAddressesByLabel_HappyPath(t *testing.T) {
	t.Parallel()
	resp := map[string]Purpose{
		"bc1qtest": {Purpose: "receive"},
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, params []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getaddressesbylabel", method)
			var label string
			_ = json.Unmarshal(params[0], &label)
			assert.Equal(t, "deposit", label)
			return b, nil
		},
	}
	res, err := GetAddressesByLabel(caller, "deposit")
	require.NoError(t, err)
	assert.Contains(t, res, "bc1qtest")
	assert.Equal(t, "receive", res["bc1qtest"].Purpose)
}

func TestGetAddressesByLabel_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("label not found")
		},
	}
	_, err := GetAddressesByLabel(caller, "unknown")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getaddressesbylabel")
}
