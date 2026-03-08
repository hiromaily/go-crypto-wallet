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
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/contract"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// Compile-time check to ensure ERC20 implements the ERC20er interface
var _ apieth.ERC20er = (*erc20)(nil)

// erc20 holds both a named Ethereum field (for EIP-1559 support) and the
// raw ethclient.Client (retained for balance, gas estimation, and nonce calls).
type erc20 struct {
	eth             apieth.ERC20Operator
	ethClient       *ethclient.Client
	tokenClient     *contract.Token
	token           domainCoin.ERC20Token
	uuidHandler     uuid.UUIDHandler
	name            string
	contractAddress string
	masterAddress   string
	decimals        int
}

func NewERC20(
	eth apieth.ERC20Operator,
	ethClient *ethclient.Client,
	tokenClient *contract.Token,
	token domainCoin.ERC20Token,
	uuidHandler uuid.UUIDHandler,
	name string,
	contractAddress string,
	masterAddress string,
	decimals int,
) *erc20 {
	return &erc20{
		eth:             eth,
		ethClient:       ethClient,
		tokenClient:     tokenClient,
		token:           token,
		uuidHandler:     uuidHandler,
		name:            name,
		contractAddress: contractAddress,
		masterAddress:   masterAddress,
		decimals:        decimals,
	}
}

// FloatToBigInt converts float64 to *big.Int using token-specific decimal precision.
// Uses big.Float to avoid float64 precision loss and int64 overflow.
func (e *erc20) FloatToBigInt(v float64) *big.Int {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(e.decimals)), nil)
	bf := new(big.Float).SetPrec(256).SetFloat64(v)
	bm := new(big.Float).SetPrec(256).SetInt(multiplier)
	result, _ := new(big.Float).Mul(bf, bm).Int(nil)
	return result
}

func (e *erc20) GetBalance(ctx context.Context, hexAddr string, _ domainETH.QuantityTag) (*big.Int, error) {
	balance, err := e.tokenClient.BalanceOf(nil, common.HexToAddress(hexAddr))
	if err != nil {
		return nil, fmt.Errorf("fail to call e.contract.BalanceOf(%s): %w", hexAddr, err)
	}
	return balance, nil
}

// validateAddresses returns an error if either address fails validation.
func validateAddresses(fromAddr, toAddr string) error {
	if pkgeth.ValidateAddr(fromAddr) != nil || pkgeth.ValidateAddr(toAddr) != nil {
		return errors.New("address validation error")
	}
	return nil
}

// resolveTokenAmount fetches the token balance for fromAddr and returns the
// amount to transfer. When requested is 0, the full balance is returned.
func (e *erc20) resolveTokenAmount(ctx context.Context, fromAddr string, requested uint64) (*big.Int, error) {
	balance, err := e.GetBalance(ctx, fromAddr, "")
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.GetBalance(): %w", err)
	}
	logger.Info("token balance", "balance", balance.String())
	if requested == 0 {
		return new(big.Int).Set(balance), nil
	}
	tokenAmount := new(big.Int).SetUint64(requested)
	if balance.Cmp(tokenAmount) < 0 {
		return nil, errors.New("balance is short to send token")
	}
	return tokenAmount, nil
}

// checkETHBalanceForGas verifies that fromAddr has sufficient ETH to cover txFee.
// ERC-20 transfers send 0 ETH but still consume ETH for gas.
func (e *erc20) checkETHBalanceForGas(ctx context.Context, fromAddr string, txFee *big.Int) error {
	ethBalance, err := e.ethClient.BalanceAt(ctx, common.HexToAddress(fromAddr), nil)
	if err != nil {
		return fmt.Errorf("fail to call ethClient.BalanceAt(): %w", err)
	}
	if ethBalance.Cmp(txFee) < 0 {
		return fmt.Errorf("insufficient ETH balance for gas: have %s Wei, need %s Wei",
			ethBalance.String(), txFee.String())
	}
	return nil
}

// buildRawTx encodes the transaction and constructs the domain RawTx entity.
func buildRawTx(
	uid, fromAddr, toAddr string, value *big.Int, nonce uint64, tx *types.Transaction,
) (*domainETH.RawTx, error) {
	rawTxHex, err := pkgeth.EncodeTx(tx)
	if err != nil {
		return nil, fmt.Errorf("fail to call encodeTx(): %w", err)
	}
	return &domainETH.RawTx{
		UUID:  uid,
		From:  fromAddr,
		To:    toAddr,
		Value: *value,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  tx.Hash().Hex(),
	}, nil
}

// resolveEIP1559Fees fetches EIP-1559 fee parameters from the connected node.
// Returns maxPriorityFeePerGas and maxFeePerGas computed as (baseFee × 2) + tip.
func (e *erc20) resolveEIP1559Fees(ctx context.Context) (maxPriorityFeePerGas, maxFeePerGas *big.Int, err error) {
	maxPriorityFeePerGas, err = e.eth.SuggestGasTipCap(ctx)
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

	maxFeePerGas = new(big.Int).Add(
		new(big.Int).Mul(blockInfo.BaseFeePerGas, big.NewInt(2)),
		maxPriorityFeePerGas,
	)
	return maxPriorityFeePerGas, maxFeePerGas, nil
}

// CreateRawTransactionEIP1559 creates an EIP-1559 (Type 2) transaction for an
// ERC-20 token transfer.
//
// Fee formula:
//
//	maxPriorityFeePerGas = SuggestGasTipCap()
//	maxFeePerGas         = (baseFeePerGas × 2) + maxPriorityFeePerGas
func (e *erc20) CreateRawTransactionEIP1559(
	ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
) (*domainETH.RawTx, *apieth.TxCreateParams, error) {
	maxPriorityFeePerGas, maxFeePerGas, err := e.resolveEIP1559Fees(ctx)
	if err != nil {
		return nil, nil, err
	}

	chainID, err := e.ethClient.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call client.ChainID(): %w", err)
	}

	if err := validateAddresses(fromAddr, toAddr); err != nil {
		return nil, nil, err
	}

	tokenAmount, err := e.resolveTokenAmount(ctx, fromAddr, amount)
	if err != nil {
		return nil, nil, err
	}

	data := e.createTransferData(toAddr, tokenAmount)

	gasLimit, err := e.estimateGas(data)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call estimateGas(): %w", err)
	}

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

	uid, err := e.uuidHandler.GenerateV7()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
	}

	rawTx, err := buildRawTx(uid.String(), fromAddr, toAddr, tokenAmount, nonce, tx)
	if err != nil {
		return nil, nil, err
	}

	txFee := new(big.Int).Mul(maxFeePerGas, new(big.Int).SetUint64(gasLimit))
	if err := e.checkETHBalanceForGas(ctx, fromAddr, txFee); err != nil {
		return nil, nil, err
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

	return rawTx, txParams, nil
}

func (*erc20) createTransferData(toAddr string, amount *big.Int) []byte {
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

func (e *erc20) estimateGas(data []byte) (uint64, error) {
	contractAddr := common.HexToAddress(e.contractAddress)
	masterAddr := common.HexToAddress(e.masterAddress)
	gasLimit, err := e.ethClient.EstimateGas(context.Background(), ethereum.CallMsg{
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
func (e *erc20) getNonce(ctx context.Context, fromAddr string, additionalNonce int) (uint64, error) {
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
