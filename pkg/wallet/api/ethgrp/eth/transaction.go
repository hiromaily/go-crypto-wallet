package eth

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// Signing a raw transaction in Go
//https://ethereum.stackexchange.com/questions/16472/signing-a-raw-transaction-in-go

// How to create Raw Transactions in Ethereum ? — Part 1
// https://medium.com/blockchain-musings/how-to-create-raw-transactions-in-ethereum-part-1-1df91abdba7c

// Create and sign OFFLINE raw transactions?
// https://ethereum.stackexchange.com/questions/3386/create-and-sign-offline-raw-transactions

// maybe useful
// https://medium.com/@akshay_111meher/creating-offline-raw-transactions-with-go-ethereum-8d6cc8174c5d

// RawTx is raw transaction
type RawTx struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	TxHex string  `json:"txhex"`
	Value big.Int `json:"value"`
	Nonce uint64  `json:"nonce"`
	Hash  string  `json:"hash"`
}

// CreateRawTransaction creates raw transaction for watch only wallet
// TODO: which QuantityTag should be used?
func (e *Ethereum) CreateRawTransaction(fromAddr, toAddr string, amount uint64) (string, []byte, error) {
	// validation check
	if e.ValidationAddr(fromAddr) != nil || e.ValidationAddr(toAddr) != nil {
		return "", nil, errors.New("address validation error")
	}

	// TODO: pending status should be included in target balance??
	// TODO: if block is still syncing, proper balance is not returned
	balance, err := e.GetBalance(fromAddr, QuantityTagLatest)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.GetBalance()")
	}
	e.logger.Info("balance", zap.Int64("balance", balance.Int64()))
	if balance.Uint64() == 0 {
		return "", nil, errors.New("balance is needed to send eth")
	}

	// nonce
	nonce, err := e.GetTransactionCount(fromAddr, QuantityTagLatest)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.GetTransactionCount()")
	}
	e.logger.Info("transactionCount", zap.Int64("nonce", nonce.Int64()))

	// gasPrice
	gasPrice, err := e.GasPrice()
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.GasPrice()")
	}
	e.logger.Info("gas_price", zap.Int64("gas_price", gasPrice.Int64()))

	var (
		txFee = new(big.Int)
		value = new(big.Int)
	)

	// calculate fees
	txFee = txFee.Mul(gasPrice, big.NewInt(int64(GasLimit)))
	// if amount==0, all amount is sent
	if amount == 0 {
		value = value.Sub(balance, txFee)
	} else {
		value = value.Sub(new(big.Int).SetUint64(amount), txFee)
		if balance.Cmp(value) == -1 {
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
			return "", nil, errors.Errorf("balance`%d` is insufficient to send `%d`", balance.Uint64(), amount)
		}
	}

	// create transaction
	// data is required when contract transaction
	tx := types.NewTransaction(nonce.Uint64(), common.HexToAddress(fromAddr), value, GasLimit, gasPrice, nil)
	txHashHex := tx.Hash().Hex()
	rawTxHex, err := encodeTx(tx)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	//RawTx
	rawtx := &RawTx{
		From:  fromAddr,
		To:    toAddr,
		TxHex: *rawTxHex,
		Value: *value,
		Nonce: nonce.Uint64(),
		Hash:  txHashHex,
	}
	bTx, err := json.Marshal(rawtx)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call json.Marshal(rawtx)")
	}

	return *rawTxHex, bTx, nil
}

// SignOnRawTransaction signs on raw transaction
// - https://ethereum.stackexchange.com/questions/16472/signing-a-raw-transaction-in-go
func (e *Ethereum) SignOnRawTransaction(rawTx *RawTx, passphrase string, accountType account.AccountType) (*RawTx, []byte, error) {
	txHex := rawTx.TxHex
	fromAddr := rawTx.From
	tx, err := decodeTx(txHex)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call decodeTx(txHex)")
	}

	// get private key
	key, err := e.GetPrivKey(fromAddr, passphrase, accountType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call e.GetPrivKey()")
	}

	// chain id
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
	chainID := big.NewInt(int64(e.netID))

	// sign
	signedTX, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key.PrivateKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call types.SignTx()")
	}
	msg, err := signedTX.AsMessage(types.NewEIP155Signer(chainID))
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to cll signedTX.AsMessage()")
	}

	encodedTx, err := encodeTx(signedTX)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	resTx := &RawTx{
		From:  msg.From().Hex(),
		To:    msg.To().Hex(),
		TxHex: *encodedTx,
		Value: *msg.Value(),
		Nonce: msg.Nonce(),
		Hash:  signedTX.Hash().Hex(),
	}

	bTx, err := json.Marshal(resTx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to all json.Marshal(resTx)")
	}

	return resTx, bTx, nil
}

// SendSignedRawTransaction sends signed raw transaction
// - SendRawTransaction in rpc_eth_tx.go
// - SendRawTx in client.go
func (e *Ethereum) SendSignedRawTransaction(signedTxHex string) (string, error) {
	decodedTx, err := decodeTx(signedTxHex)
	if err != nil {
		return "", errors.Wrap(err, "fail to call decodeTx(signedTxHex)")
	}

	txHash, err := e.SendRawTransactionWithTypesTx(decodedTx)
	if err != nil {
		return "", errors.Wrap(err, "fail to call SendRawTransactionWithTypesTx()")
	}

	return txHash, err
}

func encodeTx(tx *types.Transaction) (*string, error) {
	txb, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	txHex := hexutil.Encode(txb)
	return &txHex, nil
}

func decodeTx(txHex string) (*types.Transaction, error) {
	txc, err := hexutil.Decode(txHex)
	if err != nil {
		return nil, err
	}

	var txde types.Transaction
	err = rlp.Decode(bytes.NewReader(txc), &txde)
	if err != nil {
		return nil, err
	}

	return &txde, nil
}
