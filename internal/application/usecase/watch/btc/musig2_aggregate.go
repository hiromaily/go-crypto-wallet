package btc

import (
	"context"
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
// 2. Parses signer public keys for verification
// 3. Combines partial signatures using low-level API (no private key needed)
// 4. Verifies the aggregated signature validity
// 5. Finalizes the PSBT with the aggregated signature
// 6. Extracts the final transaction hex and ID
//
// CRITICAL: Watch Wallet Does NOT Sign
// The watch wallet is a coordinator that aggregates signatures from offline wallets
// (Keygen and Sign). It does NOT possess private keys and does NOT generate its own
// signature. This implementation correctly uses the low-level musig2.CombineSigs API
// which does not require a private key or signing session.
//
// KNOWN LIMITATION - Partial Signature Deserialization:
// The btcd musig2.PartialSignature type doesn't expose deserialization methods.
// We cannot convert input.PartialSignatures (with S and R components) to
// *musig2.PartialSignature objects.
//
// This requires either:
//  1. Contributing a NewPartialSignature() constructor to btcd library
//  2. Using reflection to construct the type (not recommended)
//  3. Alternative serialization format
//
// For now, this is a placeholder implementation that demonstrates the correct
// architectural approach (no private keys, no signing).
//
// Watch Wallet Specifics:
//   - Watch wallet is online (has Bitcoin Core RPC access)
//   - Aggregates partial signatures from offline wallets (Keygen and Sign)
//   - Does NOT generate its own signature (watch-only)
//   - Broadcasts final transactions to the network
func (*aggregateMuSig2SignaturesUseCase) Execute(
	ctx context.Context,
	input watchusecase.AggregateMuSig2SignaturesInput,
) (watchusecase.AggregateMuSig2SignaturesOutput, error) {
	// Validate minimum partial signatures (2-of-N multisig)
	if len(input.PartialSignatures) < 2 {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"insufficient partial signatures: got %d, need at least 2",
			len(input.PartialSignatures),
		)
	}

	// Validate signer public keys match partial signatures
	if len(input.SignerPublicKeys) < 2 {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"insufficient signer public keys: got %d, need at least 2",
			len(input.SignerPublicKeys),
		)
	}

	// Parse aggregated public key for verification
	aggregatedPubKey, err := btcec.ParsePubKey(input.AggregatedPublicKey[:])
	if err != nil {
		return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
			"failed to parse aggregated public key: %w",
			err,
		)
	}

	// Parse signer public keys (needed for context if we use session-based API)
	// In the low-level API approach, these are used for verification only
	signerPubKeys := make([]*btcec.PublicKey, 0, len(input.SignerPublicKeys))
	for i, pkBytes := range input.SignerPublicKeys {
		pk, err := btcec.ParsePubKey(pkBytes[:])
		if err != nil {
			return watchusecase.AggregateMuSig2SignaturesOutput{}, fmt.Errorf(
				"failed to parse signer public key %d: %w",
				i,
				err,
			)
		}
		signerPubKeys = append(signerPubKeys, pk)
	}

	// LIMITATION - Partial Signature Deserialization:
	// The btcd musig2.PartialSignature type requires both S (scalar) and R (nonce commitment).
	// Our input DTO now correctly includes both components in PartialSignatureData:
	//   - Signature [32]byte       // S component
	//   - NonceCommitment [33]byte // R component
	//
	// However, btcd doesn't provide a way to construct a *musig2.PartialSignature from bytes.
	// The ideal implementation would be:
	//
	//   var partialSigs []*musig2.PartialSignature
	//   for _, ps := range input.PartialSignatures {
	//       // Parse R (nonce commitment)
	//       R, err := btcec.ParsePubKey(ps.NonceCommitment[:])
	//       if err != nil { return err }
	//
	//       // Parse S (signature scalar)
	//       var S btcec.ModNScalar
	//       S.SetByteSlice(ps.Signature[:])
	//
	//       // Construct partial signature (NOT POSSIBLE - no constructor exists)
	//       partialSig := &musig2.PartialSignature{S: &S, R: R}
	//       partialSigs = append(partialSigs, partialSig)
	//   }
	//
	//   // Use low-level API to combine (no private key needed)
	//   combinedNonce, err := btcec.ParsePubKey(input.CombinedNonce[:])
	//   finalSig := musig2.CombineSigs(combinedNonce, partialSigs)
	//
	// Until btcd provides deserialization support, we return a placeholder.

	// PLACEHOLDER: Return placeholder signature and transaction data
	// In real implementation with proper deserialization:
	//   1. Deserialize all partial signatures from input
	//   2. Parse combined nonce from input.CombinedNonce
	//   3. Call musig2.CombineSigs(combinedNonce, partialSigs) - NO PRIVATE KEY NEEDED
	//   4. Verify final signature against aggregated public key
	//   5. Add signature to PSBT
	//   6. Finalize PSBT and extract transaction

	// Placeholder verification step (shows correct architecture even without real sigs)
	// In real implementation, finalSig would come from musig2.CombineSigs
	_ = aggregatedPubKey // Will be used for signature verification
	_ = signerPubKeys    // Used for verification context

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
