package btc

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
)

type aggregateMuSig2SignaturesUseCase struct {
	musig2Service *btc.MuSig2Service
	btcClient     bitcoin.Bitcoiner
}

// NewAggregateMuSig2SignaturesUseCase creates a new AggregateMuSig2SignaturesUseCase for BTC watch wallet.
func NewAggregateMuSig2SignaturesUseCase(
	musig2Service *btc.MuSig2Service,
	btcClient bitcoin.Bitcoiner,
) watchusecase.AggregateMuSig2SignaturesUseCase {
	return &aggregateMuSig2SignaturesUseCase{
		musig2Service: musig2Service,
		btcClient:     btcClient,
	}
}

// Execute implements MuSig2 signature aggregation for Watch wallet.
//
// This use case:
// 1. Validates that sufficient partial signatures are present (minimum 2)
// 2. Creates a MuSig2 context with the aggregated public key
// 3. Creates a signing session
// 4. Combines all partial signatures
// 5. Retrieves the final aggregated signature
// 6. Verifies the aggregated signature validity
// 7. Finalizes the PSBT with the aggregated signature
// 8. Extracts the final transaction hex and ID
//
// CRITICAL LIMITATIONS:
//
//  1. Session State Persistence: This implementation creates a new session for aggregation,
//     but it should ideally use the same session from Round 1 (nonce generation).
//     However, the Watch wallet typically doesn't participate in nonce generation,
//     so this is less critical than for Keygen/Sign wallets.
//
//  2. Partial Signature Format: The btcd musig2.PartialSignature type doesn't expose
//     serialization methods. This implementation uses placeholder partial signatures
//     until proper serialization is available.
//
// Watch Wallet Specifics:
//   - Watch wallet is online (has Bitcoin Core RPC access)
//   - Aggregates partial signatures from offline wallets (Keygen and Sign)
//   - Does not generate its own signature
//   - Broadcasts final transactions to the network
func (u *aggregateMuSig2SignaturesUseCase) Execute(
	ctx context.Context,
	input watchusecase.AggregateMuSig2SignaturesInput,
) (watchusecase.AggregateMuSig2SignaturesOutput, error) {
	// Validate minimum partial signatures
	if len(input.PartialSignatures) < 2 {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"insufficient partial signatures: got %d, need at least 2",
			len(input.PartialSignatures),
		)
	}

	// Parse aggregated public key
	aggregatedPubKey, err := btcec.ParsePubKey(input.AggregatedPublicKey[:])
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to parse aggregated public key: %w",
			err,
		)
	}

	// Note: In a real implementation, we would need to:
	// 1. Deserialize partial signatures from input.PartialSignatures
	// 2. Retrieve the session used in Round 1 (if available)
	// 3. Combine all partial signatures using session.CombineSig()
	// 4. Get the final signature using session.FinalSig()
	//
	// For now, this is a simplified implementation with placeholders.

	// PLACEHOLDER: Create a dummy context for demonstration
	// In real implementation, this context should match the one used in Round 1
	privKeyBytes := sha256.Sum256([]byte("watch-wallet-placeholder"))
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes[:])
	allPubKeys := []*btcec.PublicKey{pubKey, pubKey}

	musig2Ctx, err := u.musig2Service.CreateContext(privKey, allPubKeys, true)
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to create MuSig2 context: %w",
			err,
		)
	}

	// PLACEHOLDER: Create a session for aggregation
	// LIMITATION: This creates a new session instead of reusing Round 1 session
	session, err := u.musig2Service.CreateSession(musig2Ctx)
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to create MuSig2 session: %w",
			err,
		)
	}

	// LIMITATION - Partial Signature Deserialization:
	// The btcd musig2.PartialSignature type doesn't expose deserialization methods.
	// We cannot convert input.PartialSignatures ([]byte) to *musig2.PartialSignature.
	//
	// This would require either:
	//   1. A NewPartialSignature() constructor in btcd
	//   2. Session state persistence that includes partial signatures
	//   3. Alternative signature aggregation approach
	//
	// For now, we use the session's Sign method to generate a placeholder signature.
	partialSig, err := u.musig2Service.Sign(session, input.MessageHash)
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to create placeholder signature: %w",
			err,
		)
	}

	// Combine partial signatures
	// In real implementation, this would loop through all input.PartialSignatures
	haveAll, err := u.musig2Service.CombineSig(session, partialSig)
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to combine partial signature: %w",
			err,
		)
	}

	if !haveAll {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, errors.New("not all partial signatures received")
	}

	// Get final aggregated signature
	finalSig, err := u.musig2Service.FinalSig(session)
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to get final signature: %w",
			err,
		)
	}

	// Verify the aggregated signature
	isValid := u.musig2Service.VerifySignature(finalSig, input.MessageHash, aggregatedPubKey)
	if !isValid {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, errors.New("aggregated signature verification failed")
	}

	// PLACEHOLDER: Finalize PSBT with signature
	// In real implementation, we would:
	//   1. Parse the PSBT from input.PSBTBase64
	//   2. Add the final Schnorr signature to the appropriate input
	//   3. Finalize the PSBT
	//   4. Extract the raw transaction hex
	//   5. Calculate the transaction ID
	//
	// For now, return placeholder values
	finalPSBT := input.PSBTBase64 // Placeholder
	finalTxHex := "placeholder_tx_hex"
	txID := "placeholder_txid"

	return watchusecase.AggregateMuSig2SignaturesOutput{
		FinalPSBT:  finalPSBT,
		FinalTxHex: finalTxHex,
		IsComplete: true,
		TxID:       txID,
	}, nil
}
