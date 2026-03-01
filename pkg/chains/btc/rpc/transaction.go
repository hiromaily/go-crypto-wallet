package rpc

import (
	"encoding/json"
	"fmt"
)

// GetTransactionResult is the wire-format response of the gettransaction command.
type GetTransactionResult struct {
	Amount            float64                `json:"amount"`
	Fee               float64                `json:"fee"`
	Confirmations     uint64                 `json:"confirmations"`
	BlockHash         string                 `json:"blockhash"`
	BlockHeight       uint64                 `json:"blockheight"`
	BlockIndex        uint64                 `json:"blockindex"`
	BlockTime         uint64                 `json:"blocktime"`
	Txid              string                 `json:"txid"`
	WalletConflicts   []any                  `json:"walletconflicts"`
	MempoolConflicts  []string               `json:"mempoolconflicts,omitempty"`
	Time              int64                  `json:"time"`
	TimeReceived      int64                  `json:"timereceived"`
	Bip125Replaceable string                 `json:"bip125-replaceable"`
	Details           []GetTransactionDetail `json:"details"`
	Hex               string                 `json:"hex"`
}

// GetTransactionDetail is a detail entry in GetTransactionResult.
type GetTransactionDetail struct {
	Address   string  `json:"address"`
	Category  string  `json:"category"`
	Amount    float64 `json:"amount"`
	Label     string  `json:"label"`
	Vout      uint32  `json:"vout"`
	Fee       float64 `json:"fee,omitempty"`
	Abandoned bool    `json:"abandoned,omitempty"`
}

// TxRawResult is the wire-format response of the decoderawtransaction command.
type TxRawResult struct {
	Txid     string      `json:"txid"`
	Hash     string      `json:"hash"`
	Version  uint32      `json:"version"`
	Size     int32       `json:"size"`
	Vsize    int32       `json:"vsize"`
	Weight   int32       `json:"weight"`
	LockTime uint32      `json:"locktime"`
	Vin      []TxRawVin  `json:"vin"`
	Vout     []TxRawVout `json:"vout"`
}

// TxRawVin is an input entry in TxRawResult.
type TxRawVin struct {
	Txid      string    `json:"txid"`
	Vout      uint32    `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  uint32    `json:"sequence"`
}

// ScriptSig is the script signature in TxRawVin.
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// TxRawVout is an output entry in TxRawResult.
type TxRawVout struct {
	Value        float64      `json:"value"`
	N            uint32       `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// ScriptPubKey is the script pubkey in TxRawVout.
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   uint32   `json:"reqSigs"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

// FundRawTransactionOptions configures the fundrawtransaction call.
type FundRawTransactionOptions struct {
	FeeRate       float64 `json:"feeRate,omitempty"`
	MaxTxWeight   *int    `json:"max_tx_weight,omitempty"`
	ChangeType    string  `json:"changeType,omitempty"`
	IncludeUnsafe bool    `json:"includeUnsafe,omitempty"`
}

// FundRawTransactionResult is the wire-format response of fundrawtransaction.
type FundRawTransactionResult struct {
	Hex       string `json:"hex"`
	Fee       int64  `json:"fee"`
	Changepos int64  `json:"changepos"`
}

// SignRawTransactionResult is the wire-format response of signrawtransaction*.
type SignRawTransactionResult struct {
	Hex      string                    `json:"hex"`
	Complete bool                      `json:"complete"`
	Errors   []SignRawTransactionError `json:"errors"`
}

// SignRawTransactionError is an error entry in SignRawTransactionResult.
type SignRawTransactionError struct {
	Txid      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	ScriptSig string `json:"scriptSig"`
	Sequence  uint32 `json:"sequence"`
	Error     string `json:"error"`
}

// PrevTx is the wire format for a previous transaction input used in signing.
type PrevTx struct {
	Txid          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	RedeemScript  string  `json:"redeemScript"`
	WitnessScript string  `json:"witnessScript"`
	Amount        float64 `json:"amount"`
}

// TxInput is a transaction input for createrawtransaction.
type TxInput struct {
	TxID     string `json:"txid"`
	Vout     uint32 `json:"vout"`
	Sequence int64  `json:"sequence,omitempty"`
}

// GetTxOutResult is the wire-format response of gettxout.
type GetTxOutResult struct {
	BestBlock     string       `json:"bestblock"`
	Confirmations int64        `json:"confirmations"`
	Value         float64      `json:"value"`
	ScriptPubKey  ScriptPubKey `json:"scriptPubKey"`
	Coinbase      bool         `json:"coinbase"`
}

// GetTransaction calls gettransaction and returns the raw wire response.
func (c *Client) GetTransaction(txID string) (*GetTransactionResult, error) {
	input, err := json.Marshal(txID)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(txID): %w", err)
	}
	rawResult, err := c.client.RawRequest("gettransaction", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(gettransaction): %w", err)
	}
	var result GetTransactionResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(gettransaction): %w", err)
	}
	return &result, nil
}

// DecodeRawTransaction calls decoderawtransaction and returns the raw wire response.
func (c *Client) DecodeRawTransaction(hexTx string) (*TxRawResult, error) {
	input, err := json.Marshal(hexTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(hexTx): %w", err)
	}
	rawResult, err := c.client.RawRequest("decoderawtransaction", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(decoderawtransaction): %w", err)
	}
	var result TxRawResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(decoderawtransaction): %w", err)
	}
	return &result, nil
}

// FundRawTransaction calls fundrawtransaction with the given options.
// If opts is nil, an empty options object is sent.
func (c *Client) FundRawTransaction(
	hexTx string, opts *FundRawTransactionOptions,
) (*FundRawTransactionResult, error) {
	bHex, err := json.Marshal(hexTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(hexTx): %w", err)
	}
	if opts == nil {
		opts = &FundRawTransactionOptions{}
	}
	bOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(opts): %w", err)
	}
	rawResult, err := c.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex, bOpts})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(fundrawtransaction): %w", err)
	}
	var result FundRawTransactionResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(fundrawtransaction): %w", err)
	}
	return &result, nil
}

// SignRawTransactionWithWallet calls signrawtransactionwithwallet.
func (c *Client) SignRawTransactionWithWallet(
	hexTx string, prevTxs []PrevTx,
) (*SignRawTransactionResult, error) {
	bHex, err := json.Marshal(hexTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(hexTx): %w", err)
	}
	bPrevTxs, err := json.Marshal(prevTxs)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(prevTxs): %w", err)
	}
	rawResult, err := c.client.RawRequest("signrawtransactionwithwallet", []json.RawMessage{bHex, bPrevTxs})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(signrawtransactionwithwallet): %w", err)
	}
	var result SignRawTransactionResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(signrawtransactionwithwallet): %w", err)
	}
	return &result, nil
}

// SignRawTransactionWithKey calls signrawtransactionwithkey.
func (c *Client) SignRawTransactionWithKey(
	hexTx string, privKeys []string, prevTxs []PrevTx,
) (*SignRawTransactionResult, error) {
	bHex, err := json.Marshal(hexTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(hexTx): %w", err)
	}
	bKeys, err := json.Marshal(privKeys)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(privKeys): %w", err)
	}
	bPrevTxs, err := json.Marshal(prevTxs)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(prevTxs): %w", err)
	}
	rawResult, err := c.client.RawRequest("signrawtransactionwithkey", []json.RawMessage{bHex, bKeys, bPrevTxs})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(signrawtransactionwithkey): %w", err)
	}
	var result SignRawTransactionResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(signrawtransactionwithkey): %w", err)
	}
	return &result, nil
}

// CreateRawTransaction calls createrawtransaction.
// inputs is the list of UTXOs to spend; outputs maps recipient address to BTC amount.
// Returns the hex-encoded unsigned transaction.
func (c *Client) CreateRawTransaction(inputs []TxInput, outputs map[string]float64) (string, error) {
	bInputs, err := json.Marshal(inputs)
	if err != nil {
		return "", fmt.Errorf("fail to call json.Marshal(inputs): %w", err)
	}
	bOutputs, err := json.Marshal(outputs)
	if err != nil {
		return "", fmt.Errorf("fail to call json.Marshal(outputs): %w", err)
	}
	rawResult, err := c.client.RawRequest("createrawtransaction", []json.RawMessage{bInputs, bOutputs})
	if err != nil {
		return "", fmt.Errorf("fail to call RawRequest(createrawtransaction): %w", err)
	}
	var hexTx string
	if err := json.Unmarshal(rawResult, &hexTx); err != nil {
		return "", fmt.Errorf("fail to call json.Unmarshal(createrawtransaction): %w", err)
	}
	return hexTx, nil
}

// SendRawTransaction calls sendrawtransaction and returns the transaction ID.
func (c *Client) SendRawTransaction(hexTx string) (string, error) {
	bHex, err := json.Marshal(hexTx)
	if err != nil {
		return "", fmt.Errorf("fail to call json.Marshal(hexTx): %w", err)
	}
	rawResult, err := c.client.RawRequest("sendrawtransaction", []json.RawMessage{bHex})
	if err != nil {
		return "", fmt.Errorf("fail to call RawRequest(sendrawtransaction): %w", err)
	}
	var txid string
	if err := json.Unmarshal(rawResult, &txid); err != nil {
		return "", fmt.Errorf("fail to call json.Unmarshal(sendrawtransaction): %w", err)
	}
	return txid, nil
}

// GetTxOut calls gettxout and returns the unspent transaction output info.
func (c *Client) GetTxOut(txID string, index uint32, mempool bool) (*GetTxOutResult, error) {
	bTxID, err := json.Marshal(txID)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(txID): %w", err)
	}
	bIndex, err := json.Marshal(index)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(index): %w", err)
	}
	bMempool, err := json.Marshal(mempool)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marshal(mempool): %w", err)
	}
	rawResult, err := c.client.RawRequest("gettxout", []json.RawMessage{bTxID, bIndex, bMempool})
	if err != nil {
		return nil, fmt.Errorf("fail to call RawRequest(gettxout): %w", err)
	}
	var result GetTxOutResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(gettxout): %w", err)
	}
	return &result, nil
}
