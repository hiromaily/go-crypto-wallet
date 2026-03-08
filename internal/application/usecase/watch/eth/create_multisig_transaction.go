package eth

import (
	"context"
	"fmt"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// zeroAddrHex is the Ethereum zero address used for gas fields.
const zeroAddrHex = "0x0000000000000000000000000000000000000000"

type createETHMultisigTransactionUseCase struct {
	nonceReader  apieth.SafeNonceReader
	hashComputer apieth.SafeTxHashComputer
	multisigRepo portfile.MultisigFileRepositorier
	uuidHandler  pkguuid.UUIDHandler
	chainID      uint64
}

// NewCreateETHMultisigTransactionUseCase constructs the use case for creating
// an unsigned Safe multisig transaction proposal file.
func NewCreateETHMultisigTransactionUseCase(
	nonceReader apieth.SafeNonceReader,
	hashComputer apieth.SafeTxHashComputer,
	multisigRepo portfile.MultisigFileRepositorier,
	uuidHandler pkguuid.UUIDHandler,
	chainID uint64,
) watchusecase.CreateETHMultisigTransactionUseCase {
	return &createETHMultisigTransactionUseCase{
		nonceReader:  nonceReader,
		hashComputer: hashComputer,
		multisigRepo: multisigRepo,
		uuidHandler:  uuidHandler,
		chainID:      chainID,
	}
}

// Execute creates an unsigned multisig transaction proposal and writes it to a JSON file.
// No database record is created. The file path is returned in the output.
func (u *createETHMultisigTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.CreateETHMultisigTransactionInput,
) (watchusecase.CreateETHMultisigTransactionOutput, error) {
	// 1. Generate a unique proposal UUID.
	id, err := u.uuidHandler.GenerateV4()
	if err != nil {
		return watchusecase.CreateETHMultisigTransactionOutput{}, fmt.Errorf("failed to generate UUID: %w", err)
	}
	uuidStr := id.String()

	// 2. Convert Ether amount to Wei.
	weiValue := pkgeth.FromFloatEther(input.AmountEther)

	// 3. Fetch the Safe contract's current nonce.
	nonce, err := u.nonceReader.GetSafeNonce(ctx, input.SafeAddress)
	if err != nil {
		return watchusecase.CreateETHMultisigTransactionOutput{}, fmt.Errorf("failed to get Safe nonce: %w", err)
	}

	// 4. Compute the EIP-712 safeTxHash via the Safe contract.
	// For simple ETH transfers: call data is empty, all gas fields are zero.
	safeTxHash, err := u.hashComputer.GetSafeTxHash(
		ctx,
		input.SafeAddress,
		input.To,
		weiValue,
		[]byte{}, // empty call data for plain ETH transfer
		0,        // operation: Call
		nonce,
	)
	if err != nil {
		return watchusecase.CreateETHMultisigTransactionOutput{},
			fmt.Errorf("failed to compute safeTxHash: %w", err)
	}

	// 5. Build the unsigned multisig transaction file.
	actionType := domainTx.ActionType(input.ActionType)
	file := &dtoeth.ETHMultisigTransactionFile{
		Version:        1,
		TxType:         "unsigned",
		UUID:           uuidStr,
		SafeAddress:    input.SafeAddress,
		To:             input.To,
		Value:          weiValue.String(),
		Data:           "0x",
		Operation:      0,
		SafeTxGas:      "0",
		BaseGas:        "0",
		GasPrice:       "0",
		GasToken:       zeroAddrHex,
		RefundReceiver: zeroAddrHex,
		Nonce:          nonce.String(),
		ChainID:        u.chainID,
		SafeTxHash:     safeTxHash,
		Threshold:      input.Threshold,
		Signatures:     []dtoeth.SignEntry{},
	}

	// 6. Derive the canonical file path and write the file.
	filePath := u.multisigRepo.CreateMultisigFilePath(actionType, uuidStr, 0)
	writtenPath, err := u.multisigRepo.WriteETHMultisigJSONFile(filePath, file)
	if err != nil {
		return watchusecase.CreateETHMultisigTransactionOutput{},
			fmt.Errorf("failed to write multisig JSON file: %w", err)
	}

	return watchusecase.CreateETHMultisigTransactionOutput{
		FilePath: writtenPath,
		UUID:     uuidStr,
	}, nil
}
