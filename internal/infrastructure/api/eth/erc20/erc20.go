package erc20

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/ethtx"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/contract"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// Compile-time check to ensure ERC20 implements the ERC20er interface
var _ apieth.ERC20er = (*ERC20)(nil)

// ERC20 struct holds both a named Ethereum field (for EIP-1559 support) and the
// raw ethclient.Client (retained for balance, gas estimation, and nonce calls that
// use it directly until Tasks 3.3–3.4 delegate them through the eth field).
type ERC20 struct {
	eth             apieth.ERC20NodeAPI
	client          *ethclient.Client
	tokenClient     *contract.Token
	token           domainCoin.ERC20Token
	uuidHandler     uuid.UUIDHandler
	name            string
	contractAddress string
	masterAddress   string
	decimals        int
}

func NewERC20(
	eth apieth.ERC20NodeAPI,
	client *ethclient.Client,
	tokenClient *contract.Token,
	token domainCoin.ERC20Token,
	uuidHandler uuid.UUIDHandler,
	name string,
	contractAddress string,
	masterAddress string,
	decimals int,
) *ERC20 {
	return &ERC20{
		eth:             eth,
		client:          client,
		tokenClient:     tokenClient,
		token:           token,
		uuidHandler:     uuidHandler,
		name:            name,
		contractAddress: contractAddress,
		masterAddress:   masterAddress,
		decimals:        decimals,
	}
}

// func (e *ERC20) getOption(
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

func (*ERC20) ValidateAddr(addr string) error {
	// validation check
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("address:%s is invalid", addr)
	}
	return nil
}

// FloatToBigInt converts float64 to *big.Int
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

func (e *ERC20) GetBalance(ctx context.Context, hexAddr string, _ domainETH.QuantityTag) (*big.Int, error) {
	balance, err := e.tokenClient.BalanceOf(nil, common.HexToAddress(hexAddr))
	if err != nil {
		return nil, fmt.Errorf("fail to call e.contract.BalanceOf(%s): %w", hexAddr, err)
	}
	return balance, nil
}

// CreateRawTransaction creates raw transaction for watch only wallet
//   - Transferring Tokens (ERC-20)
//     https://goethereumbook.org/en/transfer-tokens/
//   - Transfer ERC20 Tokens Using Golang
//     https://www.youtube.com/watch?v=-Epg5Ub-fA0
//     https://github.com/what-the-func/golang-ethereum-transfer-tokens/blob/master/main.go
//
// Note:
// - master address takes fee
// - sender account delegates transfer to master address
// - 1. call approve(address spender, uint256 amount) by fromA, spender is masterAddr
// -  this task may be separated from normal flow `create tx`
// - => approve requires gas to call ... this pattern is impossible
// - 1.b. Or after approve is called, this transaction may be sent
func (e *ERC20) CreateRawTransaction(
	ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
) (*domainETH.RawTx, *apieth.TxCreateParams, error) {
	// validation check
	if e.ValidateAddr(fromAddr) != nil || e.ValidateAddr(toAddr) != nil {
		return nil, nil, errors.New("address validation error")
	}
	logger.Debug("eth.CreateRawTransaction()",
		"fromAddr", fromAddr,
		"toAddr", toAddr,
		"amount", amount,
	)

	balance, err := e.GetBalance(ctx, fromAddr, "")
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.GetBalance(): %w", err)
	}
	logger.Info("balance", "balance", balance.Int64())
	if balance.Uint64() < amount {
		return nil, nil, errors.New("balance is short to send token")
	}
	tokenAmount := big.NewInt(int64(amount))
	if amount == 0 {
		tokenAmount = balance
	}

	data := e.createTransferData(toAddr, tokenAmount)
	gasLimit, err := e.estimateGas(data)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call estimateGas(data): %w", err)
	}

	gasPrice, err := e.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call client.SuggestGasPrice(): %w", err)
	}

	// nonce
	nonce, err := e.getNonce(ctx, fromAddr, additionalNonce)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call e.getNonce(): %w", err)
	}

	logger.Debug("comparison",
		"Nonce", nonce,
		"TokenAmount", tokenAmount.Uint64(),
		"GasLimit", gasLimit,
		"GasPrice", gasPrice.Uint64(),
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
		return nil, nil, fmt.Errorf("fail to call encodeTx(): %w", err)
	}

	// generate UUID to trace transaction because unsignedTx is not unique
	uid, err := e.uuidHandler.GenerateV7()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
	}

	// create domain RawTx (infrastructure type)
	infraRawTx := &ethtx.RawTx{
		UUID:  uid.String(),
		From:  fromAddr,
		To:    toAddr,
		Value: *tokenAmount,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}

	// Calculate transaction fee (gasPrice * gasLimit)
	txFee := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasLimit))

	// create TxCreateParams DTO for use case layer
	txParams := &apieth.TxCreateParams{
		UUID:        uid.String(),
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      tokenAmount.Uint64(),
		Fee:         txFee.Uint64(),
		GasLimit:    uint32(gasLimit),
		Nonce:       nonce,
		EthTxType:   0, // legacy transaction
		GasPrice:    gasPrice.Uint64(),
	}

	// Convert infrastructure RawTx to domain RawTx
	return ethtx.ToDomainRawTx(infraRawTx), txParams, nil
}

// SupportsEIP1559 delegates to the underlying Ethereum node to detect EIP-1559 support.
// Returns true when the connected node supports EIP-1559 (e.g., Anvil, post-London geth).
func (e *ERC20) SupportsEIP1559(ctx context.Context) bool {
	return e.eth.SupportsEIP1559(ctx)
}

// CreateRawTransactionEIP1559 creates an EIP-1559 (Type 2) transaction for an
// ERC-20 token transfer.
//
// When the connected node does not support EIP-1559 (e.g., pre-London private
// chains), the method falls back to a legacy Type 0 transaction via
// CreateRawTransaction.
//
// Fee formula (identical to Ethereum.CreateRawTransactionEIP1559):
//
//	maxPriorityFeePerGas = SuggestGasTipCap()
//	maxFeePerGas         = (baseFeePerGas × 2) + maxPriorityFeePerGas
//
// The calldata is ABI-encoded transfer(address,uint256) with method selector
// 0xa9059cbb, same as the legacy path.
func (e *ERC20) CreateRawTransactionEIP1559(
	ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
) (*domainETH.RawTx, *apieth.TxCreateParams, error) {
	// Fall back to legacy when the node does not support EIP-1559.
	if !e.SupportsEIP1559(ctx) {
		return e.CreateRawTransaction(ctx, fromAddr, toAddr, amount, additionalNonce)
	}

	// ── EIP-1559 fee parameters ──────────────────────────────────────────────
	maxPriorityFeePerGas, err := e.eth.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.SuggestGasTipCap(): %w", err)
	}

	currentBlockNum, err := e.eth.BlockNumber(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.BlockNumber(): %w", err)
	}

	blockInfo, err := e.eth.GetBlockByNumber(ctx, currentBlockNum.Uint64())
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.GetBlockByNumber(): %w", err)
	}

	if blockInfo.BaseFeePerGas == nil {
		return nil, nil, errors.New("baseFeePerGas not found in block (EIP-1559 not activated)")
	}

	// maxFeePerGas = (baseFee × 2) + tip — same formula as Ethereum.CreateRawTransactionEIP1559
	baseFeeTimesTwo := new(big.Int).Mul(blockInfo.BaseFeePerGas, big.NewInt(2))
	maxFeePerGas := new(big.Int).Add(baseFeeTimesTwo, maxPriorityFeePerGas)

	// ── Chain ID ─────────────────────────────────────────────────────────────
	chainID, err := e.client.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call client.ChainID(): %w", err)
	}

	// ── Address validation ───────────────────────────────────────────────────
	if err := e.ValidateAddr(fromAddr); err != nil {
		return nil, nil, fmt.Errorf("invalid fromAddr: %w", err)
	}
	if err := e.ValidateAddr(toAddr); err != nil {
		return nil, nil, fmt.Errorf("invalid toAddr: %w", err)
	}

	// ── Token balance ────────────────────────────────────────────────────────
	balance, err := e.GetBalance(ctx, fromAddr, "")
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.GetBalance(): %w", err)
	}
	logger.Info("token balance", "balance", balance.String())
	tokenAmount := new(big.Int).SetUint64(amount)
	if amount == 0 {
		tokenAmount = new(big.Int).Set(balance)
	} else if balance.Cmp(tokenAmount) < 0 {
		return nil, nil, errors.New("balance is short to send token")
	}

	// ── Calldata & gas ───────────────────────────────────────────────────────
	data := e.createTransferData(toAddr, tokenAmount)
	gasLimit, err := e.estimateGas(data)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call estimateGas(): %w", err)
	}

	// ── Nonce ────────────────────────────────────────────────────────────────
	nonce, err := e.getNonce(ctx, fromAddr, additionalNonce)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call e.getNonce(): %w", err)
	}

	logger.Debug("EIP-1559 ERC-20 tx parameters",
		"nonce", nonce,
		"tokenAmount", tokenAmount.Uint64(),
		"gasLimit", gasLimit,
		"maxPriorityFeePerGas", maxPriorityFeePerGas.Uint64(),
		"maxFeePerGas", maxFeePerGas.Uint64(),
		"chainID", chainID.Uint64(),
	)

	// ── Build DynamicFeeTx ────────────────────────────────────────────────────
	contractAddr := common.HexToAddress(e.contractAddress)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       gasLimit,
		To:        &contractAddr,
		Value:     new(big.Int), // ERC-20 transfers send 0 ETH
		Data:      data,
	})

	txHash := tx.Hash().Hex()
	rawTxHex, err := ethtx.EncodeTx(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call encodeTx(): %w", err)
	}

	uid, err := e.uuidHandler.GenerateV7()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
	}

	txFee := new(big.Int).Mul(maxFeePerGas, new(big.Int).SetUint64(gasLimit))

	infraRawTx := &ethtx.RawTx{
		UUID:  uid.String(),
		From:  fromAddr,
		To:    toAddr,
		Value: *tokenAmount,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}

	txParams := &apieth.TxCreateParams{
		UUID:                 uid.String(),
		FromAddress:          fromAddr,
		ToAddress:            toAddr,
		Amount:               tokenAmount.Uint64(),
		Fee:                  txFee.Uint64(),
		GasLimit:             uint32(gasLimit),
		Nonce:                nonce,
		EthTxType:            2, // EIP-1559 Type 2
		ChainID:              chainID.Uint64(),
		MaxFeePerGas:         maxFeePerGas.Uint64(),
		MaxPriorityFeePerGas: maxPriorityFeePerGas.Uint64(),
	}

	return ethtx.ToDomainRawTx(infraRawTx), txParams, nil
}

func (*ERC20) createTransferData(toAddr string, amount *big.Int) []byte {
	// function signature as a byte slice
	transferFnSignature := []byte("transfer(address,uint256)")

	// methodID of function
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	// set parameter for account: to address
	paddedToAddr := common.LeftPadBytes(common.HexToAddress(toAddr).Bytes(), 32)
	// set parameter for amount
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// create data
	data := make([]byte, 0, len(methodID)+len(paddedToAddr)+len(paddedAmount))
	data = append(data, methodID...)
	data = append(data, paddedToAddr...)
	data = append(data, paddedAmount...)

	return data
}

func (e *ERC20) estimateGas(data []byte) (uint64, error) {
	contractAddr := common.HexToAddress(e.contractAddress)
	masterAddr := common.HexToAddress(e.masterAddress)
	gasLimit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: masterAddr,
		To:   &contractAddr,
		Data: data,
	})
	if err != nil {
		return 0, fmt.Errorf("fail to call client.EstimateGas(): %w", err)
	}
	return gasLimit, nil
}

// getNonce retrieves the pending nonce for fromAddr by delegating to e.eth.GetTransactionCount.
// It mirrors the nonce retrieval in Ethereum.getNonce, removing the previous duplication
// that called e.client.PendingNonceAt directly (resolves task 3.4 FIXME).
func (e *ERC20) getNonce(ctx context.Context, fromAddr string, additionalNonce int) (uint64, error) {
	nonce, err := e.eth.GetTransactionCount(ctx, fromAddr, domainETH.QuantityTagPending)
	if err != nil {
		return 0, fmt.Errorf("fail to call eth.GetTransactionCount(): %w", err)
	}
	if additionalNonce != 0 {
		nonce = nonce.Add(nonce, new(big.Int).SetUint64(uint64(additionalNonce)))
	}
	logger.Debug("nonce",
		"eth.GetTransactionCount(fromAddr, QuantityTagPending)", nonce.Uint64(),
	)
	return nonce.Uint64(), nil
}
