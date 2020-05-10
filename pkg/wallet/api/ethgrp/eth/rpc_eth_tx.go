package eth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

// ResponseGetTransaction response of eth_getTransactionByHash
type ResponseGetTransaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      int64  `json:"blockNumber"`
	From             string `json:"from"`
	Gas              int64  `json:"gas"`
	GasPrice         int64  `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            int64  `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex int64  `json:"transactionIndex"`
	Value            int64  `json:"value"`
	V                int64  `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

// ResponseGetTransactionReceipt response of eth_getTransactionReceipt
type ResponseGetTransactionReceipt struct {
	TransactionHash   string   `json:"transactionHash"`
	TransactionIndex  int64    `json:"transactionIndex"`
	BlockHash         string   `json:"blockHash"`
	BlockNumber       int64    `json:"blockNumber"`
	From              string   `json:"from"`
	To                string   `json:"to"`
	CumulativeGasUsed int64    `json:"cumulativeGasUsed"`
	GasUsed           int64    `json:"gasUsed"`
	ContractAddress   string   `json:"contractAddress"`
	Logs              []string `json:"logs"`
	LogsBloom         string   `json:"logsBloom"`
	Status            int64    `json:"status"`
}

// Sign calculates an Ethereum specific signature with:
//  sign(keccak256("\x19Ethereum Signed Message:\n" + len(message) + message)))
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (e *Ethereum) Sign(hexAddr, message string) (string, error) {

	var signature string
	err := e.rpcClient.CallContext(e.ctx, &signature, "eth_sign", hexAddr, message)
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_sign)")
	}

	return signature, nil
}

// SendTransaction トランザクションを送信し、トランザクションhashを返す
// FIXME: Invalid params: Invalid bytes format. Expected a 0x-prefixed hex string with even length
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendtransaction
func (e *Ethereum) SendTransaction(msg ethereum.CallMsg) (string, error) {

	var txHash string
	err := e.rpcClient.CallContext(e.ctx, &txHash, "eth_sendTransaction", toCallArg(msg))
	if err != nil {
		//FIXME: Invalid params: Invalid bytes format. Expected a 0x-prefixed hex string with even length.
		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_sendTransaction)")
	}

	return txHash, nil
}

// SendRawTransaction creates new message call transaction or a contract creation for signed transactions
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func (e *Ethereum) SendRawTransaction(signedTx string) (string, error) {

	var txHash string
	err := e.rpcClient.CallContext(e.ctx, &txHash, "eth_sendRawTransaction", signedTx)
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_sendTransaction)")
	}

	return txHash, nil
}

// SendRawTransactionWithTypesTx call SendRawTransaction() by types.Transaction
func (e *Ethereum) SendRawTransactionWithTypesTx(tx *types.Transaction) (string, error) {
	encodedTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return "", errors.Wrap(err, "fail to call rlp.EncodeToBytes(tx)")
	}
	return e.SendRawTransaction(hexutil.Encode(encodedTx))
}

// Call executes a new message call immediately without creating a transaction on the block chain
// FIXME: check is not done yet
//func (e *Ethereum) Call(msg ethereum.CallMsg) (string, error) {
//	var txHash string
//	err := e.rpcClient.CallContext(e.ctx, &txHash, "eth_call", msg)
//	if err != nil {
//		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_call)")
//	}
//
//	return txHash, nil
//}

// GetTransactionByHash returns the information about a transaction requested by transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash
func (e *Ethereum) GetTransactionByHash(hashTx string) (*ResponseGetTransaction, error) {

	var resMap map[string]string
	err := e.rpcClient.CallContext(e.ctx, &resMap, "eth_getTransactionByHash", hashTx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getTransactionByHash)")
	}
	if len(resMap) == 0 {
		return nil, errors.New("response of eth_getTransactionByHash is empty")
	}

	blockNumber, err := hexutil.DecodeBig(setZeroHex(resMap["blockNumber"])) //blockNumber string = ""
	if err != nil {
		return nil, errors.New("response[blockNumber] is invalid")
	}
	gas, err := hexutil.DecodeBig(setZeroHex(resMap["gas"])) //gas string = "0x5208"
	if err != nil {
		return nil, errors.New("response[gas] is invalid")
	}
	gasPrice, err := hexutil.DecodeBig(setZeroHex(resMap["gasPrice"])) //gasPrice string = "0x0"
	if err != nil {
		return nil, errors.New("response[gasPrice] is invalid")
	}
	nonce, err := hexutil.DecodeBig(setZeroHex(resMap["nonce"])) //nonce string = "0x0"
	if err != nil {
		return nil, errors.New("response[nonce] is invalid")
	}
	transactionIndex, err := hexutil.DecodeBig(setZeroHex(resMap["transactionIndex"])) //transactionIndex string = ""
	if err != nil {
		return nil, errors.New("response[transactionIndex] is invalid")
	}
	value, err := hexutil.DecodeBig(setZeroHex(resMap["value"])) //value string = "0xde0b6b3a7640000"
	if err != nil {
		return nil, errors.New("response[value] is invalid")
	}
	v, err := hexutil.DecodeBig(setZeroHex(resMap["v"])) //v string = "0x2a"
	if err != nil {
		return nil, errors.New("response[v] is invalid")
	}

	return &ResponseGetTransaction{
		BlockHash:        resMap["blockHash"],
		BlockNumber:      blockNumber.Int64(),
		From:             resMap["from"],
		Gas:              gas.Int64(),
		GasPrice:         gasPrice.Int64(),
		Hash:             resMap["hash"],
		Input:            resMap["input"],
		Nonce:            nonce.Int64(),
		To:               resMap["to"],
		TransactionIndex: transactionIndex.Int64(),
		Value:            value.Int64(),
		V:                v.Int64(),
		R:                resMap["r"],
		S:                resMap["s"],
	}, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
// Note, tis is not available for pending transactions
func (e *Ethereum) GetTransactionReceipt(hashTx string) (*ResponseGetTransactionReceipt, error) {

	// timeout
	ch := make(chan error, 1)
	// FIXME: timeout configuration
	ctx, cancel := context.WithTimeout(e.ctx, 10*time.Second)
	defer func() {
		cancel()
	}()

	//call
	var resMap map[string]interface{}
	go func() {
		err := e.rpcClient.CallContext(ctx, &resMap, "eth_getTransactionReceipt", hashTx)
		if err != nil {
			ch <- errors.Wrap(err, "fail to call rpc.CallContext(eth_getTransactionReceipt)")
		}
		ch <- nil
	}()

	//wait by timeout
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err == context.Canceled {
			e.logger.Debug("context.Canceled for calling eth_getTransactionReceipt")
		} else if err == context.DeadlineExceeded {
			e.logger.Debug("context.DeadlineExceeded for calling eth_getTransactionReceipt")
		} else if err != nil {
			e.logger.Debug(err.Error())
			return nil, err
		}
	case retErr := <-ch:
		if retErr != nil {
			return nil, retErr
		}
	}

	if len(resMap) == 0 {
		return nil, errors.New("response is empty")
	}

	transactionHash, err := castToString(resMap["transactionHash"])
	if err != nil {
		return nil, errors.New("response[transactionHash] is invalid")
	}
	transactionIndex, err := castToInt64(resMap["transactionIndex"])
	if err != nil {
		return nil, errors.New("response[transactionIndex] is invalid")
	}
	blockHash, err := castToString(resMap["blockHash"])
	if err != nil {
		return nil, errors.New("response[blockHash] is invalid")
	}
	blockNumber, err := castToInt64(resMap["blockNumber"])
	if err != nil {
		return nil, errors.New("response[blockNumber] is invalid")
	}
	from, err := castToString(resMap["from"])
	if err != nil {
		return nil, errors.New("response[from] is invalid")
	}
	to, err := castToString(resMap["to"])
	if err != nil {
		return nil, errors.New("response[to] is invalid")
	}
	cumulativeGasUsed, err := castToInt64(resMap["cumulativeGasUsed"])
	if err != nil {
		return nil, errors.New("response[cumulativeGasUsed] is invalid")
	}
	gasUsed, err := castToInt64(resMap["gasUsed"])
	if err != nil {
		return nil, errors.New("response[gasUsed] is invalid")
	}
	//contractAddress would be nil sometimes
	var contractAddress string
	if resMap["contractAddress"] == nil {
		contractAddress = ""
	} else {
		contractAddress, err = castToString(resMap["contractAddress"])
		if err != nil {
			return nil, errors.New("response[contractAddress] is invalid")
		}
	}
	// logs would be empty interface{} sometimes
	logs, err := castToSliceString(resMap["logs"])
	if err != nil {
		return nil, errors.New("response[logs] is invalid")
	}

	logsBloom, err := castToString(resMap["logsBloom"])
	if err != nil {
		return nil, errors.New("response[logsBloom] is invalid")
	}
	status, err := castToInt64(resMap["status"])
	if err != nil {
		return nil, errors.New("response[status] is invalid")
	}

	return &ResponseGetTransactionReceipt{
		TransactionHash:   transactionHash,
		TransactionIndex:  transactionIndex,
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              from,
		To:                to,
		CumulativeGasUsed: cumulativeGasUsed,
		GasUsed:           gasUsed,
		ContractAddress:   contractAddress,
		Logs:              logs,
		LogsBloom:         logsBloom,
		Status:            status,
	}, nil
}
