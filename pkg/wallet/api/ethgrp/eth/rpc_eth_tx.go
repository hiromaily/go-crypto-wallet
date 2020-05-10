package eth

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
)

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
func (e *Ethereum) Call(msg ethereum.CallMsg) (string, error) {
	var txHash string
	err := e.rpcClient.CallContext(e.ctx, &txHash, "eth_call", msg)
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_call)")
	}

	return txHash, nil
}
