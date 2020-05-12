package eth

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// RawTx is raw transaction
type RawTx struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Value big.Int `json:"value"`
	Nonce uint64  `json:"nonce"`
	TxHex string  `json:"txhex"`
	Hash  string  `json:"hash"`
}

// TODO: WIP: logic is not fixed
func (e *Ethereum) getNonce(fromAddr string) (uint64, error) {
	// by calling GetTransactionCount()
	nonce1, err := e.GetTransactionCount(fromAddr, QuantityTagPending)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call eth.GetTransactionCount()")
	}

	// Or by calling
	nonce2, err := e.ethClient.PendingNonceAt(e.ctx, common.HexToAddress(fromAddr))
	if err != nil {
		return 0, errors.Wrap(err, "fail to call ethClient.PendingNonceAt()")
	}

	e.logger.Debug("nonce",
		zap.Uint64("GetTransactionCount(fromAddr, QuantityTagLatest)", nonce1.Uint64()),
		zap.Uint64("ethClient.PendingNonceAt(e.ctx, common.HexToAddress(fromAddr))", nonce2),
	)

	return nonce2, nil
}

// How to calculate transaction fee?
// https://ethereum.stackexchange.com/questions/19665/how-to-calculate-transaction-fee
func (e *Ethereum) calculateFee(fromAddr, toAddr common.Address, balance, gasPrice, value *big.Int) (*big.Int, *big.Int, *big.Int, error) {

	msg := &ethereum.CallMsg{
		From:     fromAddr,
		To:       &toAddr,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    nil,
		Data:     nil,
	}
	// gasLimit
	estimatedGas, err := e.EstimateGas(msg)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "fail to call EstimateGas()")
	}
	//txFee := gasPrice * estimatedGas
	txFee := new(big.Int).Mul(gasPrice, estimatedGas)
	//newValue := value - txFee
	newValue := new(big.Int)
	if value.Uint64() == 0 {
		newValue = newValue.Sub(balance, txFee)
	} else {
		newValue = newValue.Sub(value, txFee)
		if balance.Cmp(value) == -1 {
			//   -1 if x <  y
			//    0 if x == y
			//   +1 if x >  y
			return nil, nil, nil, errors.Errorf("balance`%d` is insufficient to send `%d`", balance.Uint64(), newValue.Uint64())
		}
	}

	return newValue, txFee, estimatedGas, nil
}

// CreateRawTransaction creates raw transaction for watch only wallet
// TODO: which QuantityTag should be used?
// - Creating offline/raw transactions with Go-Ethereum
//   https://medium.com/@akshay_111meher/creating-offline-raw-transactions-with-go-ethereum-8d6cc8174c5d
func (e *Ethereum) CreateRawTransaction(fromAddr, toAddr string, amount uint64) (string, []byte, error) {
	// validation check
	if e.ValidationAddr(fromAddr) != nil || e.ValidationAddr(toAddr) != nil {
		return "", nil, errors.New("address validation error")
	}

	// TODO: pending status should be included in target balance??
	// TODO: if block is still syncing, proper balance is not returned
	balance, err := e.GetBalance(fromAddr, QuantityTagPending)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.GetBalance()")
	}
	e.logger.Info("balance", zap.Int64("balance", balance.Int64()))
	if balance.Uint64() == 0 {
		return "", nil, errors.New("balance is needed to send eth")
	}

	// nonce
	nonce, err := e.getNonce(fromAddr)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.GetTransactionCount()")
	}

	// gasPrice
	//e.ethClient.SuggestGasPrice()
	gasPrice, err := e.GasPrice()
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.GasPrice()")
	}
	e.logger.Info("gas_price", zap.Int64("gas_price", gasPrice.Int64()))

	//fromAddr, toAddr common.Address, gasPrice, value *big.Int
	newValue, txFee, estimatedGas, err := e.calculateFee(
		common.HexToAddress(fromAddr),
		common.HexToAddress(toAddr),
		balance,
		gasPrice,
		new(big.Int).SetUint64(amount),
	)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call eth.calculateFee()")
	}

	//TODO: which value should be used for args of types.NewTransaction()
	e.logger.Debug("comparison",
		zap.Uint64("GasLimit", GasLimit),
		zap.Uint64("estimatedGas", estimatedGas.Uint64()),
		zap.Uint64("txFee", txFee.Uint64()))

	// create transaction
	// data is required when contract transaction
	// NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddr), newValue, GasLimit, gasPrice, nil)
	txHash := tx.Hash().Hex()
	rawTxHex, err := encodeTx(tx)
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	//RawTx
	rawtx := &RawTx{
		From:  fromAddr,
		To:    toAddr,
		Value: *newValue,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
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
		Value: *msg.Value(),
		Nonce: msg.Nonce(),
		TxHex: *encodedTx,
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
