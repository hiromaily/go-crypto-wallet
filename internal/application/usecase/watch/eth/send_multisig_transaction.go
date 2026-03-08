package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type sendETHMultisigTransactionUseCase struct {
	safeExec     apieth.SafeExecuter
	multisigRepo portfile.MultisigFileRepositorier
}

// NewSendETHMultisigTransactionUseCase constructs the use case for submitting a
// fully signed Safe multisig transaction to Ethereum.
func NewSendETHMultisigTransactionUseCase(
	safeExec apieth.SafeExecuter,
	multisigRepo portfile.MultisigFileRepositorier,
) watchusecase.SendETHMultisigTransactionUseCase {
	return &sendETHMultisigTransactionUseCase{
		safeExec:     safeExec,
		multisigRepo: multisigRepo,
	}
}

// Execute reads the signed multisig file, assembles SafeExecParams, and submits
// the transaction via the Safe contract. Returns the submitted Ethereum tx hash.
func (u *sendETHMultisigTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.SendETHMultisigTransactionInput,
) (watchusecase.SendETHMultisigTransactionOutput, error) {
	// 1. Read and validate the multisig transaction file.
	f, err := u.multisigRepo.ReadETHMultisigJSONFile(input.FilePath)
	if err != nil {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("failed to read multisig JSON file: %w", err)
	}

	// 2. Guard: reject unsigned proposals.
	if f.TxType != "signed" {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("%w: got tx_type %q", dtoeth.ErrNotFullySigned, f.TxType)
	}

	// 3. Build concatenated signatures sorted by signer address ascending (required by Safe).
	sigBytes, err := buildSortedSignatures(f.Signatures)
	if err != nil {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("failed to build signatures: %w", err)
	}

	// 4. Parse numeric fields from decimal strings.
	value, ok := new(big.Int).SetString(f.Value, 10)
	if !ok {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("invalid value field %q", f.Value)
	}
	safeNonce, ok := new(big.Int).SetString(f.Nonce, 10)
	if !ok {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("invalid nonce field %q", f.Nonce)
	}
	safeTxGas, ok := new(big.Int).SetString(f.SafeTxGas, 10)
	if !ok {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("invalid safe_tx_gas field %q", f.SafeTxGas)
	}
	baseGas, ok := new(big.Int).SetString(f.BaseGas, 10)
	if !ok {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("invalid base_gas field %q", f.BaseGas)
	}
	gasPrice, ok := new(big.Int).SetString(f.GasPrice, 10)
	if !ok {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("invalid gas_price field %q", f.GasPrice)
	}

	// 5. Decode call data hex ("0x" → empty; other values → decoded bytes).
	callData, err := decodeHexData(f.Data)
	if err != nil {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("invalid data field: %w", err)
	}

	// 6. Assemble SafeExecParams and submit.
	params := apieth.SafeExecParams{
		SafeAddress:    f.SafeAddress,
		To:             f.To,
		Value:          value,
		Data:           callData,
		Operation:      f.Operation,
		SafeTxGas:      safeTxGas,
		BaseGas:        baseGas,
		GasPrice:       gasPrice,
		GasToken:       f.GasToken,
		RefundReceiver: f.RefundReceiver,
		Nonce:          safeNonce,
		Signatures:     sigBytes,
	}

	txHash, err := u.safeExec.ExecuteSafeTransaction(ctx, params)
	if err != nil {
		return watchusecase.SendETHMultisigTransactionOutput{},
			fmt.Errorf("failed to execute Safe transaction: %w", err)
	}

	logger.Info("Safe multisig transaction submitted",
		"uuid", f.UUID,
		"safe", f.SafeAddress,
		"to", f.To,
		"tx_hash", txHash,
	)

	return watchusecase.SendETHMultisigTransactionOutput{TxHash: txHash}, nil
}

// buildSortedSignatures sorts signature entries by signer address ascending
// (case-insensitive) and concatenates the decoded 65-byte signatures.
// The Safe contract requires this ordering for signature verification.
func buildSortedSignatures(entries []dtoeth.SignEntry) ([]byte, error) {
	sorted := make([]dtoeth.SignEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].SignerAddress) < strings.ToLower(sorted[j].SignerAddress)
	})

	var result []byte
	for _, entry := range sorted {
		sigHex := strings.TrimPrefix(entry.SignatureHex, "0x")
		sigBytes, err := hex.DecodeString(sigHex)
		if err != nil {
			return nil, fmt.Errorf("invalid signature hex for %s: %w", entry.SignerAddress, err)
		}
		result = append(result, sigBytes...)
	}
	return result, nil
}

// decodeHexData decodes a 0x-prefixed hex string to bytes.
// "0x" with no data returns an empty byte slice.
func decodeHexData(hexStr string) ([]byte, error) {
	stripped := strings.TrimPrefix(hexStr, "0x")
	if stripped == "" {
		return []byte{}, nil
	}
	return hex.DecodeString(stripped)
}
