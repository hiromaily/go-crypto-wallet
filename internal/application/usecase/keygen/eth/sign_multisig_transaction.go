package eth

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type signMultisigTransactionUseCase struct {
	multisigRepo   portfile.MultisigFileRepositorier
	accountKeyRepo repocold.ETHAccountKeyRepositorier
}

// NewSignMultisigTransactionUseCase constructs the use case for offline EIP-712 signing
// of a Safe multisig proposal. No network calls are made — suitable for air-gapped wallets.
// This use case is shared by both ETHKeygen and ETHSign wallet adapters.
func NewSignMultisigTransactionUseCase(
	multisigRepo portfile.MultisigFileRepositorier,
	accountKeyRepo repocold.ETHAccountKeyRepositorier,
) keygenusecase.SignMultisigTransactionUseCase {
	return &signMultisigTransactionUseCase{
		multisigRepo:   multisigRepo,
		accountKeyRepo: accountKeyRepo,
	}
}

// Sign reads the multisig file, verifies the safeTxHash offline via EIP-712 recomputation,
// signs the hash with the signer's private key, appends the signature, and writes the updated
// file with an incremented counter in the path.
func (u *signMultisigTransactionUseCase) Sign(
	ctx context.Context,
	input keygenusecase.SignMultisigTransactionInput,
) (keygenusecase.SignMultisigTransactionOutput, error) {
	// 1. Read the multisig transaction file.
	f, err := u.multisigRepo.ReadETHMultisigJSONFile(input.FilePath)
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("failed to read multisig file: %w", err)
	}

	// 2. Recompute safeTxHash offline and verify against the file.
	computedHash, err := ComputeSafeTxHash(f)
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("failed to compute safeTxHash: %w", err)
	}

	fileHashHex := strings.TrimPrefix(f.SafeTxHash, "0x")
	fileHashBytes, err := hex.DecodeString(fileHashHex)
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("invalid safeTxHash in file: %w", err)
	}
	if !bytes.Equal(computedHash[:], fileHashBytes) {
		return keygenusecase.SignMultisigTransactionOutput{}, dtoeth.ErrSafeTxHashMismatch
	}

	// 3. Guard against duplicate signers.
	signerAddrLower := strings.ToLower(input.SignerAddress)
	for _, entry := range f.Signatures {
		if strings.ToLower(entry.SignerAddress) == signerAddrLower {
			return keygenusecase.SignMultisigTransactionOutput{}, dtoeth.ErrDuplicateSigner
		}
	}

	// 4. Look up the signer's account key and derive the child private key.
	accountKey, err := u.accountKeyRepo.GetByAddress(input.SignerAddress)
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("key not found for address %s: %w", input.SignerAddress, err)
	}
	if accountKey.AccountExtendedPrivkey == "" {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("accountXpriv not found for address %s: run key generation first", input.SignerAddress)
	}

	privKey, err := deriveChildPrivKey(accountKey.AccountExtendedPrivkey, uint32(accountKey.Idx))
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("failed to derive private key for %s: %w", input.SignerAddress, err)
	}

	// 5. Sign the hash. crypto.Sign returns [R || S || V] with V ∈ {0, 1}.
	//    Safe contract requires V = 27 or 28 (Ethereum legacy convention), so V += 27.
	sig, err := crypto.Sign(computedHash[:], privKey)
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("failed to sign safeTxHash: %w", err)
	}
	sig[64] += 27

	sigHex := "0x" + hex.EncodeToString(sig)

	// 6. Append the new signature.
	f.Signatures = append(f.Signatures, dtoeth.SignEntry{
		SignerAddress: input.SignerAddress,
		SignatureHex:  sigHex,
	})

	// 7. Transition to "signed" when threshold is reached.
	signCount := len(f.Signatures)
	isComplete := signCount >= f.Threshold
	if isComplete {
		f.TxType = string(domainTx.TxTypeSigned)
	}

	// 8. Write the updated file with an incremented counter in the path.
	newPath := u.multisigRepo.CreateMultisigFilePath(domainTx.ActionType(f.ActionType), f.UUID, signCount)

	outPath, err := u.multisigRepo.WriteETHMultisigJSONFile(newPath, f)
	if err != nil {
		return keygenusecase.SignMultisigTransactionOutput{},
			fmt.Errorf("failed to write signed multisig file: %w", err)
	}

	logger.Info("Safe multisig signature appended",
		"uuid", f.UUID,
		"safe", f.SafeAddress,
		"signer", input.SignerAddress,
		"signCount", signCount,
		"threshold", f.Threshold,
		"isComplete", isComplete,
	)

	return keygenusecase.SignMultisigTransactionOutput{
		FilePath:   outPath,
		IsComplete: isComplete,
		SignCount:  signCount,
		Threshold:  f.Threshold,
	}, nil
}
