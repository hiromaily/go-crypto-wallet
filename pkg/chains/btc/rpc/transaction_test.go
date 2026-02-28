package rpc

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── GetTransaction ───────────────────────────────────────────────────────────

func TestGetTransaction_HappyPath(t *testing.T) {
	t.Parallel()
	resp := GetTransactionResult{
		Txid:          "abc123",
		Confirmations: 6,
		Amount:        0.01,
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, params []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "gettransaction", method)
			var txid string
			_ = json.Unmarshal(params[0], &txid)
			assert.Equal(t, "abc123", txid)
			return b, nil
		},
	}
	res, err := GetTransaction(caller, "abc123")
	require.NoError(t, err)
	assert.Equal(t, "abc123", res.Txid)
	assert.Equal(t, uint64(6), res.Confirmations)
}

func TestGetTransaction_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("tx not found")
		},
	}
	_, err := GetTransaction(caller, "abc123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gettransaction")
}

// ─── DecodeRawTransaction ────────────────────────────────────────────────────

func TestDecodeRawTransaction_HappyPath(t *testing.T) {
	t.Parallel()
	resp := TxRawResult{Txid: "abc123", Version: 2}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "decoderawtransaction", method)
			return b, nil
		},
	}
	res, err := DecodeRawTransaction(caller, "deadbeef")
	require.NoError(t, err)
	assert.Equal(t, "abc123", res.Txid)
}

func TestDecodeRawTransaction_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("decode error")
		},
	}
	_, err := DecodeRawTransaction(caller, "bad_hex")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decoderawtransaction")
}

// ─── FundRawTransaction ───────────────────────────────────────────────────────

func TestFundRawTransaction_HappyPath(t *testing.T) {
	t.Parallel()
	resp := FundRawTransactionResult{Hex: "funded_hex", Fee: 1000, Changepos: 1}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "fundrawtransaction", method)
			return b, nil
		},
	}
	opts := &FundRawTransactionOptions{FeeRate: 0.0001}
	res, err := FundRawTransaction(caller, "raw_hex", opts)
	require.NoError(t, err)
	assert.Equal(t, "funded_hex", res.Hex)
	assert.Equal(t, int64(1000), res.Fee)
}

func TestFundRawTransaction_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("insufficient funds")
		},
	}
	_, err := FundRawTransaction(caller, "raw_hex", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fundrawtransaction")
}

// ─── SignRawTransactionWithWallet ─────────────────────────────────────────────

func TestSignRawTransactionWithWallet_HappyPath(t *testing.T) {
	t.Parallel()
	resp := SignRawTransactionResult{Hex: "signed_hex", Complete: true}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "signrawtransactionwithwallet", method)
			return b, nil
		},
	}
	res, err := SignRawTransactionWithWallet(caller, "raw_hex", []PrevTx{})
	require.NoError(t, err)
	assert.True(t, res.Complete)
	assert.Equal(t, "signed_hex", res.Hex)
}

func TestSignRawTransactionWithWallet_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("sign error")
		},
	}
	_, err := SignRawTransactionWithWallet(caller, "raw_hex", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "signrawtransactionwithwallet")
}

// ─── SignRawTransactionWithKey ────────────────────────────────────────────────

func TestSignRawTransactionWithKey_HappyPath(t *testing.T) {
	t.Parallel()
	resp := SignRawTransactionResult{Hex: "signed_hex", Complete: false}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "signrawtransactionwithkey", method)
			return b, nil
		},
	}
	res, err := SignRawTransactionWithKey(caller, "raw_hex", []string{"wif_key"}, []PrevTx{})
	require.NoError(t, err)
	assert.Equal(t, "signed_hex", res.Hex)
}

func TestSignRawTransactionWithKey_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("sign error")
		},
	}
	_, err := SignRawTransactionWithKey(caller, "raw_hex", []string{"wif"}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "signrawtransactionwithkey")
}

// ─── CreateRawTransaction ────────────────────────────────────────────────────

func TestCreateRawTransaction_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal("raw_hex_result")
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "createrawtransaction", method)
			return b, nil
		},
	}
	inputs := []TxInput{{TxID: "prev_txid", Vout: 0}}
	outputs := map[string]float64{"bc1qrecipient": 0.001}
	hexTx, err := CreateRawTransaction(caller, inputs, outputs)
	require.NoError(t, err)
	assert.Equal(t, "raw_hex_result", hexTx)
}

func TestCreateRawTransaction_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("create error")
		},
	}
	_, err := CreateRawTransaction(caller, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "createrawtransaction")
}

// ─── SendRawTransaction ───────────────────────────────────────────────────────

func TestSendRawTransaction_HappyPath(t *testing.T) {
	t.Parallel()
	b, err := json.Marshal("txid_result")
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "sendrawtransaction", method)
			return b, nil
		},
	}
	txid, err := SendRawTransaction(caller, "signed_hex")
	require.NoError(t, err)
	assert.Equal(t, "txid_result", txid)
}

func TestSendRawTransaction_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("broadcast error")
		},
	}
	_, err := SendRawTransaction(caller, "bad_hex")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "sendrawtransaction")
}

// ─── GetTxOut ────────────────────────────────────────────────────────────────

func TestGetTxOut_HappyPath(t *testing.T) {
	t.Parallel()
	resp := GetTxOutResult{Confirmations: 100, Value: 0.65, Coinbase: false}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	caller := &mockRPCCaller{
		rawRequestFn: func(method string, _ []json.RawMessage) (json.RawMessage, error) {
			assert.Equal(t, "gettxout", method)
			return b, nil
		},
	}
	res, err := GetTxOut(caller, "txid123", 0, false)
	require.NoError(t, err)
	assert.Equal(t, int64(100), res.Confirmations)
	assert.InDelta(t, 0.65, res.Value, 1e-10)
}

func TestGetTxOut_Error(t *testing.T) {
	t.Parallel()
	caller := &mockRPCCaller{
		rawRequestFn: func(_ string, _ []json.RawMessage) (json.RawMessage, error) {
			return nil, errors.New("txout error")
		},
	}
	_, err := GetTxOut(caller, "bad_txid", 0, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gettxout")
}
