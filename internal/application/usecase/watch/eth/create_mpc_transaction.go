package eth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

type createMPCTransactionUseCase struct {
	txCreator   apieth.TxCreator
	mpcFileRepo portfile.MPCFileRepositorier
	uuidHandler pkguuid.UUIDHandler
}

// NewCreateMPCTransactionUseCase constructs the use case for creating an unsigned MPC
// transaction file.
func NewCreateMPCTransactionUseCase(
	txCreator apieth.TxCreator,
	mpcFileRepo portfile.MPCFileRepositorier,
	uuidHandler pkguuid.UUIDHandler,
) watchusecase.CreateMPCTransactionUseCase {
	return &createMPCTransactionUseCase{
		txCreator:   txCreator,
		mpcFileRepo: mpcFileRepo,
		uuidHandler: uuidHandler,
	}
}

// Execute validates inputs, builds an unsigned EIP-1559 transaction, computes its hash, and
// writes an ETHMPCTransactionFile so that MPC nodes can cooperatively sign it.
func (u *createMPCTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.CreateMPCTransactionInput,
) (watchusecase.CreateMPCTransactionOutput, error) {
	// 1. Validate EIP-55 checksummed addresses.
	if !common.IsHexAddress(input.FromAddress) {
		return watchusecase.CreateMPCTransactionOutput{},
			fmt.Errorf("invalid from address %q: %w", input.FromAddress, dtoeth.ErrInvalidMPCAddress)
	}
	if !common.IsHexAddress(input.ToAddress) {
		return watchusecase.CreateMPCTransactionOutput{},
			fmt.Errorf("invalid to address %q: %w", input.ToAddress, dtoeth.ErrInvalidMPCAddress)
	}

	// 2. Generate a UUID for the MPC transaction file.
	id, err := u.uuidHandler.GenerateV4()
	if err != nil {
		return watchusecase.CreateMPCTransactionOutput{}, fmt.Errorf("failed to generate UUID: %w", err)
	}
	uuidStr := id.String()

	// 3. Convert Ether to Wei and create the unsigned EIP-1559 transaction.
	weiAmount := pkgeth.FromFloatEther(input.AmountEther)
	rawTx, txParams, err := u.txCreator.CreateRawTransactionEIP1559(
		ctx, input.FromAddress, input.ToAddress, weiAmount.Uint64(), 0)
	if err != nil {
		return watchusecase.CreateMPCTransactionOutput{},
			fmt.Errorf("failed to create EIP-1559 transaction: %w", err)
	}

	// 4. Populate the ETHMPCTransactionFile.
	file := &dtoeth.ETHMPCTransactionFile{
		Version:              1,
		TxType:               string(domainTx.TxTypeUnsigned),
		UUID:                 uuidStr,
		ActionType:           input.ActionType,
		From:                 input.FromAddress,
		To:                   input.ToAddress,
		Value:                weiAmount.String(),
		Nonce:                txParams.Nonce,
		GasLimit:             uint64(txParams.GasLimit),
		ChainID:              txParams.ChainID,
		MaxFeePerGas:         strconv.FormatUint(txParams.MaxFeePerGas, 10),
		MaxPriorityFeePerGas: strconv.FormatUint(txParams.MaxPriorityFeePerGas, 10),
		TxHash:               rawTx.Hash,
		RawTxHex:             rawTx.TxHex,
		Threshold:            input.Threshold,
		PartyIDs:             input.PartyIDs,
	}

	// 5. Write the unsigned file.
	writtenPath, err := u.mpcFileRepo.WriteETHMPCJSONFile(file, false)
	if err != nil {
		return watchusecase.CreateMPCTransactionOutput{},
			fmt.Errorf("failed to write MPC transaction file: %w", err)
	}

	return watchusecase.CreateMPCTransactionOutput{
		FilePath: writtenPath,
		UUID:     uuidStr,
		TxHash:   rawTx.Hash,
	}, nil
}
