package rpc

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── EstimateSmartFee ─────────────────────────────────────────────────────────

func TestEstimateSmartFee_HappyPath(t *testing.T) {
	t.Parallel()
	resp := EstimateSmartFeeResult{FeeRate: 0.00012345, Blocks: 6}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, params []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "estimatesmartfee", method)
			require.Len(t, params, 1)
			var blocks int
			_ = json.Unmarshal(params[0], &blocks)
			assert.Equal(t, 6, blocks)
			return b, nil
		},
	}
	feeRate, err := EstimateSmartFee(caller, 6)
	require.NoError(t, err)
	assert.InDelta(t, 0.00012345, feeRate, 1e-10)
}

func TestEstimateSmartFee_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("fee estimation unavailable")
		},
	}
	_, err := EstimateSmartFee(caller, 6)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "estimatesmartfee")
}

// ─── GetBalance ───────────────────────────────────────────────────────────────

func TestGetBalance_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(1.23456789)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getbalance", method)
			return b, nil
		},
	}
	amount, err := GetBalance(caller, 6)
	require.NoError(t, err)
	assert.InDelta(t, 1.23456789, amount, 1e-10)
}

func TestGetBalance_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("wallet locked")
		},
	}
	_, err := GetBalance(caller, 6)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getbalance")
}

// ─── SetLabel ─────────────────────────────────────────────────────────────────

func TestSetLabel_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(nil)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, params []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "setlabel", method)
			var addr, label string
			_ = json.Unmarshal(params[0], &addr)
			_ = json.Unmarshal(params[1], &label)
			assert.Equal(t, "bc1qtest", addr)
			assert.Equal(t, "deposit", label)
			return b, nil
		},
	}
	err = SetLabel(caller, "bc1qtest", "deposit")
	require.NoError(t, err)
}

func TestSetLabel_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("label error")
		},
	}
	err := SetLabel(caller, "bad_addr", "label")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "setlabel")
}

// ─── AddMultisigAddress ───────────────────────────────────────────────────────

func TestAddMultisigAddress_HappyPath(t *testing.T) {
	t.Parallel()
	resp := AddMultisigAddressResult{Address: "2Nmsig", RedeemScript: "52210...ae"}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "addmultisigaddress", method)
			return b, nil
		},
	}
	res, err := AddMultisigAddress(caller, 2, []string{"addr1", "addr2"}, "payment", "p2sh-segwit")
	require.NoError(t, err)
	assert.Equal(t, "2Nmsig", res.Address)
}

func TestAddMultisigAddress_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("multisig error")
		},
	}
	_, err := AddMultisigAddress(caller, 2, []string{"a1", "a2"}, "account", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "addmultisigaddress")
}

// ─── Logging ──────────────────────────────────────────────────────────────────

func TestLogging_HappyPath(t *testing.T) {
	t.Parallel()
	resp := LoggingResult{Net: true, RPC: true}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "logging", method)
			return b, nil
		},
	}
	res, err := Logging(caller)
	require.NoError(t, err)
	assert.True(t, res.Net)
	assert.True(t, res.RPC)
}

func TestLogging_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("logging error")
		},
	}
	_, err := Logging(caller)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "logging")
}

// ─── ImportAddress ────────────────────────────────────────────────────────────

func TestImportAddress_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(nil)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "importaddress", method)
			return b, nil
		},
	}
	err = ImportAddress(caller, "bc1qtest", "deposit", false)
	require.NoError(t, err)
}

func TestImportAddress_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("import error")
		},
	}
	err := ImportAddress(caller, "bad", "", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "importaddress")
}

// ─── ImportPrivKey ────────────────────────────────────────────────────────────

func TestImportPrivKey_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(nil)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "importprivkey", method)
			return b, nil
		},
	}
	err = ImportPrivKey(caller, "cTWIF...", "deposit", false)
	require.NoError(t, err)
}

func TestImportPrivKey_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("import error")
		},
	}
	err := ImportPrivKey(caller, "bad_wif", "", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "importprivkey")
}

// ─── ImportDescriptors ───────────────────────────────────────────────────────

func TestImportDescriptors_HappyPath(t *testing.T) {
	t.Parallel()
	resp := []ImportDescriptorsResponse{{Success: true}}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "importdescriptors", method)
			return b, nil
		},
	}
	reqs := []ImportDescriptorsRequest{{Descriptor: "wpkh(...)#checksum", Timestamp: "now", Active: true}}
	res, err := ImportDescriptors(caller, reqs)
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.True(t, res[0].Success)
}

func TestImportDescriptors_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("import error")
		},
	}
	_, err := ImportDescriptors(caller, []ImportDescriptorsRequest{{Descriptor: "bad", Timestamp: "now"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "importdescriptors")
}

// ─── ImportMulti ──────────────────────────────────────────────────────────────

func TestImportMulti_HappyPath(t *testing.T) {
	t.Parallel()
	resp := []ImportMultiResponse{{Success: true}}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "importmulti", method)
			return b, nil
		},
	}
	reqs := []ImportMultiRequest{{ScriptPubKey: map[string]string{"address": "bc1qtest"}, Timestamp: "now"}}
	res, err := ImportMulti(caller, reqs, nil)
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.True(t, res[0].Success)
}

func TestImportMulti_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("import error")
		},
	}
	_, err := ImportMulti(caller, []ImportMultiRequest{{ScriptPubKey: "x", Timestamp: "now"}}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "importmulti")
}

// ─── GetDescriptorInfo ───────────────────────────────────────────────────────

func TestGetDescriptorInfo_HappyPath(t *testing.T) {
	t.Parallel()
	resp := DescriptorInfo{Descriptor: "wpkh(...)#checksum", Checksum: "checksum", IsSolvable: true}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "getdescriptorinfo", method)
			return b, nil
		},
	}
	res, err := GetDescriptorInfo(caller, "wpkh(...)")
	require.NoError(t, err)
	assert.True(t, res.IsSolvable)
	assert.Equal(t, "checksum", res.Checksum)
}

func TestGetDescriptorInfo_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("descriptor error")
		},
	}
	_, err := GetDescriptorInfo(caller, "bad_descriptor")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getdescriptorinfo")
}

// ─── ListDescriptors ─────────────────────────────────────────────────────────

func TestListDescriptors_HappyPath(t *testing.T) {
	t.Parallel()
	resp := ListDescriptorsResult{
		WalletName:  "test_wallet",
		Descriptors: []DescriptorEntry{{Desc: "wpkh(...)#checksum", Active: true}},
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "listdescriptors", method)
			return b, nil
		},
	}
	res, err := ListDescriptors(caller, false)
	require.NoError(t, err)
	assert.Equal(t, "test_wallet", res.WalletName)
	require.Len(t, res.Descriptors, 1)
}

func TestListDescriptors_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("list error")
		},
	}
	_, err := ListDescriptors(caller, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "listdescriptors")
}

// ─── CreateWallet ─────────────────────────────────────────────────────────────

func TestCreateWallet_HappyPath(t *testing.T) {
	t.Parallel()
	resp := LoadWalletResult{WalletName: "new_wallet"}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "createwallet", method)
			return b, nil
		},
	}
	err = CreateWallet(caller, "new_wallet", &CreateWalletOptions{Descriptors: true})
	require.NoError(t, err)
}

func TestCreateWallet_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("wallet exists")
		},
	}
	err := CreateWallet(caller, "existing_wallet", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "createwallet")
}

// ─── LoadWallet ───────────────────────────────────────────────────────────────

func TestLoadWallet_HappyPath(t *testing.T) {
	t.Parallel()
	resp := LoadWalletResult{WalletName: "my_wallet"}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "loadwallet", method)
			return b, nil
		},
	}
	err = LoadWallet(caller, "my_wallet")
	require.NoError(t, err)
}

func TestLoadWallet_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("wallet not found")
		},
	}
	err := LoadWallet(caller, "missing_wallet")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "loadwallet")
}

// ─── UnloadWallet ─────────────────────────────────────────────────────────────

func TestUnloadWallet_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal(nil)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "unloadwallet", method)
			return b, nil
		},
	}
	err = UnloadWallet(caller, "my_wallet")
	require.NoError(t, err)
}

func TestUnloadWallet_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("unload error")
		},
	}
	err := UnloadWallet(caller, "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unloadwallet")
}
