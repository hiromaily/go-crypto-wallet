package erc20

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"

	"github.com/hiromaily/go-crypto-wallet/pkg/contract"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/ethtx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

type ERC20 struct {
	client          *ethclient.Client
	tokenClient     *contract.Token
	token           coin.ERC20Token
	name            string
	contractAddress string
	masterAddress   string
	decimals        int
	logger          *zap.Logger
}

func NewERC20(
	client *ethclient.Client,
	tokenClient *contract.Token,
	token coin.ERC20Token,
	name string,
	contractAddress string,
	masterAddress string,
	decimals int,
	logger *zap.Logger,
) *ERC20 {
	return &ERC20{
		client:          client,
		tokenClient:     tokenClient,
		token:           token,
		name:            name,
		contractAddress: contractAddress,
		masterAddress:   masterAddress,
		decimals:        decimals,
		logger:          logger,
	}
}

//func (e *ERC20) getOption(
//	ctx context.Context,
//	isPending bool,
//	fromAddr common.Address,
//	blockNumber *big.Int) *bind.CallOpts {
//
//	opts := bind.CallOpts{}
//	if ctx != nil {
//		opts.Context = ctx
//	}
//	opts.Pending = isPending
//	opts.From = fromAddr
//	if blockNumber != nil {
//		opts.BlockNumber = blockNumber
//	}
//	return &opts
//}

func (e *ERC20) ValidateAddr(addr string) error {
	// validation check
	if !common.IsHexAddress(addr) {
		return errors.Errorf("address:%s is invalid", addr)
	}
	return nil
}

// FIXME: Is it correct to handle decimal??
func (e *ERC20) FloatToBigInt(v float64) *big.Int {
	if e.decimals == 18 {
		return big.NewInt(int64(v * 1e18))
	}
	// v * math.Pow(10, float64(e.decimals))
	for i := 0; i < e.decimals; i++ {
		v *= 10
	}
	return big.NewInt(int64(v))
}

func (e *ERC20) GetBalance(hexAddr string, _ eth.QuantityTag) (*big.Int, error) {
	balance, err := e.tokenClient.BalanceOf(nil, common.HexToAddress(hexAddr))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call e.contract.BalanceOf(%s)", hexAddr)
	}
	return balance, nil
}

// CreateRawTransaction creates raw transaction for watch only wallet
// - Transferring Tokens (ERC-20)
//   https://goethereumbook.org/en/transfer-tokens/
// - Transfer ERC20 Tokens Using Golang
//   https://www.youtube.com/watch?v=-Epg5Ub-fA0
//   https://github.com/what-the-func/golang-ethereum-transfer-tokens/blob/master/main.go
// Note: master address takes fee
// - if sender sends 5ETH, receiver receives 5ETH
// - master address has to pay 5ETH + fee
func (e *ERC20) CreateRawTransaction(fromAddr, toAddr string, amount uint64, additionalNonce int) (*ethtx.RawTx, *models.EthDetailTX, error) {
	// validation check
	if e.ValidateAddr(fromAddr) != nil || e.ValidateAddr(toAddr) != nil {
		return nil, nil, errors.New("address validation error")
	}
	e.logger.Debug("eth.CreateRawTransaction()",
		zap.String("fromAddr", fromAddr),
		zap.String("toAddr", toAddr),
		zap.Uint64("amount", amount),
	)

	balance, err := e.GetBalance(fromAddr, "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call eth.GetBalance()")
	}
	e.logger.Info("balance", zap.Int64("balance", balance.Int64()))
	if balance.Uint64() < amount {
		return nil, nil, errors.New("balance is short to send token")
	}
	amountToken := big.NewInt(int64(amount))

	data := e.createTransferData(toAddr, amountToken)
	gasLimit, err := e.estimateGas(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call estimateGas(data)")
	}
	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call client.SuggestGasPrice()")
	}

	// nonce
	nonce, err := e.getNonce(fromAddr, additionalNonce)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call e.getNonce()")
	}

	e.logger.Debug("comparison",
		zap.Uint64("Nonce", nonce),
		zap.Uint64("TokenAmount", amountToken.Uint64()),
		zap.Uint64("GasLimit", gasLimit),
		zap.Uint64("GasPrice", gasPrice.Uint64()),
	)

	// create transaction
	contractAddr := common.HexToAddress(e.contractAddress)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddr,
		Value:    new(big.Int), // value must be 0 for ERC-20
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	// From here, same as CreateRawTransaction() in ethgrop/eth/transaction.go
	txHash := tx.Hash().Hex()
	rawTxHex, err := ethtx.EncodeTx(tx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call encodeTx()")
	}

	// generate UUID to trace transaction because unsignedTx is not unique
	uid := uuid.NewV4().String()

	// create insert data for　eth_detail_tx
	txDetailItem := &models.EthDetailTX{
		UUID:            uid,
		SenderAccount:   "",
		SenderAddress:   fromAddr,
		ReceiverAccount: "",
		ReceiverAddress: toAddr,
		Amount:          amountToken.Uint64(),
		Fee:             0, // later update is required
		GasLimit:        uint32(gasLimit),
		Nonce:           nonce,
		UnsignedHexTX:   *rawTxHex,
	}

	// RawTx
	rawtx := &ethtx.RawTx{
		UUID:  uid,
		From:  fromAddr,
		To:    toAddr,
		Value: *amountToken,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}
	return rawtx, txDetailItem, nil
}

func (e *ERC20) createTransferData(toAddr string, amount *big.Int) []byte {
	// function signature as a byte slice
	transferFnSignature := []byte("transfer(address,uint256)")

	// methodID of function
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// account to address
	paddedToAddr := common.LeftPadBytes(common.HexToAddress(toAddr).Bytes(), 32)

	// set amount
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// create data
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedToAddr...)
	data = append(data, paddedAmount...)

	return data
}

func (e *ERC20) estimateGas(data []byte) (uint64, error) {
	contractAddr := common.HexToAddress(e.contractAddress)
	gasLimit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	})
	if err != nil {
		return 0, errors.Wrap(err, "fail to call client.EstimateGas()")
	}
	return gasLimit, nil
}

// FIXME: this logic is almost same to where getNonce() in ethgrp/eth/transaction.go
func (e *ERC20) getNonce(fromAddr string, additionalNonce int) (uint64, error) {
	nonce, err := e.client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddr))
	if err != nil {
		return 0, errors.Wrap(err, "fail to call ethClient.PendingNonceAt()")
	}
	nonce += uint64(additionalNonce)

	e.logger.Debug("nonce",
		zap.Uint64("client.PendingNonceAt(e.ctx, common.HexToAddress(fromAddr))", nonce),
	)
	return nonce, nil
}
