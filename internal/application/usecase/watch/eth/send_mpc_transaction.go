package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ErrSenderMismatch is returned when the address recovered from the TSS signature does not
// match the from address in the MPC transaction file. Broadcasting is aborted.
var ErrSenderMismatch = errors.New("recovered sender does not match from address in transaction file")

type sendMPCTransactionUseCase struct {
	mpcSigner   apieth.MPCTransactionSigner
	txSender    apieth.TxSender
	mpcFileRepo portfile.MPCFileRepositorier
}

// NewSendMPCTransactionUseCase constructs the use case for initiating a TSS signing session
// and broadcasting the resulting signed transaction.
func NewSendMPCTransactionUseCase(
	mpcSigner apieth.MPCTransactionSigner,
	txSender apieth.TxSender,
	mpcFileRepo portfile.MPCFileRepositorier,
) watchusecase.SendMPCTransactionUseCase {
	return &sendMPCTransactionUseCase{
		mpcSigner:   mpcSigner,
		txSender:    txSender,
		mpcFileRepo: mpcFileRepo,
	}
}

// Execute reads the unsigned MPC transaction file, initiates a TSS signing session with the
// participating nodes, applies the 65-byte ECDSA signature, verifies the recovered sender,
// broadcasts the signed transaction, and writes the signed file.
func (u *sendMPCTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.SendMPCTransactionInput,
) (watchusecase.SendMPCTransactionOutput, error) {
	// 1. Read and validate the unsigned MPC transaction file.
	f, err := u.mpcFileRepo.ReadETHMPCJSONFile(input.FilePath)
	if err != nil {
		return watchusecase.SendMPCTransactionOutput{},
			fmt.Errorf("failed to read MPC transaction file: %w", err)
	}

	// 2. Guard: only unsigned files can initiate signing.
	if f.TxType != string(domainTx.TxTypeUnsigned) {
		return watchusecase.SendMPCTransactionOutput{},
			fmt.Errorf("%w: got tx_type %q", dtoeth.ErrNotUnsigned, f.TxType)
	}

	// 3. Build the MPCSigningRequest from file fields.
	sigReq := apieth.MPCSigningRequest{
		SessionID: f.UUID,
		Hash:      common.FromHex(f.TxHash),
		PartyIDs:  f.PartyIDs,
		Threshold: f.Threshold,
		PeerAddrs: input.PeerAddrs,
	}

	// 4. Initiate the TSS signing session and collect the 65-byte signature.
	sigResult, err := u.mpcSigner.SignTransaction(ctx, sigReq)
	if err != nil {
		return watchusecase.SendMPCTransactionOutput{},
			fmt.Errorf("TSS signing session failed (session=%s): %w", f.UUID, err)
	}

	// 5. Apply the signature to the raw transaction.
	signedTxHex, err := u.applySignature(f, sigResult.Signature)
	if err != nil {
		return watchusecase.SendMPCTransactionOutput{}, err
	}

	// 6. Broadcast the signed transaction.
	txHash, err := u.txSender.SendSignedRawTransaction(ctx, signedTxHex)
	if err != nil {
		return watchusecase.SendMPCTransactionOutput{},
			fmt.Errorf("failed to broadcast signed transaction: %w", err)
	}

	// 7. Write the signed file.
	f.TxType = string(domainTx.TxTypeSigned)
	f.SignedTxHex = signedTxHex
	if _, writeErr := u.mpcFileRepo.WriteETHMPCJSONFile(f, true); writeErr != nil {
		// Non-fatal: tx already broadcast. Log and continue.
		logger.Warn("failed to write signed MPC transaction file",
			"uuid", f.UUID, "error", writeErr)
	}

	logger.Info("MPC transaction broadcast",
		"uuid", f.UUID,
		"from", f.From,
		"to", f.To,
		"txHash", txHash,
	)

	return watchusecase.SendMPCTransactionOutput{TxHash: txHash}, nil
}

// applySignature decodes the raw unsigned transaction, applies the 65-byte TSS signature,
// verifies that the recovered sender matches file.From, and returns the signed hex.
func (*sendMPCTransactionUseCase) applySignature(
	f *dtoeth.ETHMPCTransactionFile,
	sig []byte,
) (string, error) {
	if len(sig) != 65 {
		return "", fmt.Errorf("expected 65-byte signature, got %d bytes", len(sig))
	}

	// Decode the unsigned transaction.
	tx, err := pkgeth.DecodeTx(f.RawTxHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode raw transaction: %w", err)
	}

	signer := types.LatestSignerForChainID(new(big.Int).SetUint64(f.ChainID))

	// Normalise v: WithSignature for EIP-1559 expects v ∈ {0, 1}, not 27/28.
	sigCopy := make([]byte, 65)
	copy(sigCopy, sig)
	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}

	signedTx, err := tx.WithSignature(signer, sigCopy)
	if err != nil {
		return "", fmt.Errorf("failed to apply TSS signature to transaction: %w", err)
	}

	// Verify that the recovered sender matches the expected from address.
	recovered, err := types.Sender(signer, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to recover sender from signed transaction: %w", err)
	}
	expected := common.HexToAddress(f.From)
	if recovered != expected {
		return "", fmt.Errorf("%w: recovered=%s expected=%s", ErrSenderMismatch,
			recovered.Hex(), expected.Hex())
	}

	// Encode the signed transaction.
	encoded, err := pkgeth.EncodeTx(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to encode signed transaction: %w", err)
	}

	return *encoded, nil
}
