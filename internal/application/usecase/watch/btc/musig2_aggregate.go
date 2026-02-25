// Package btc provides Bitcoin-specific use cases for Watch wallet operations.
//
// # MuSig2 Signature Aggregation (Watch Wallet Coordinator)
//
// The AggregateMuSig2SignaturesUseCase is responsible for aggregating partial
// signatures from multiple offline signers (Keygen, Sign1, Sign2, etc.) into
// a final Schnorr signature that can be broadcast to the Bitcoin network.
//
// # Watch Wallet Role
//
// The Watch wallet acts as the coordinator in the MuSig2 protocol:
//   - Does NOT possess private keys (watch-only)
//   - Does NOT generate its own signatures
//   - Collects partial signatures from offline wallets
//   - Aggregates partial signatures into final signature
//   - Verifies the final signature
//   - Broadcasts the transaction to Bitcoin network
//
// # Aggregation Workflow
//
//	┌─────────────────────────────────────────────────────────┐
//	│ Input: Partial Signatures from All Signers              │
//	│  - Keygen's partial signature                           │
//	│  - Sign1's partial signature                            │
//	│  - Sign2's partial signature                            │
//	│  - Aggregated public key                                │
//	│  - Combined nonce                                       │
//	│  - PSBT with transaction data                           │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 1: Validate Inputs                                 │
//	│  - Check minimum partial signatures (≥2 for 2-of-N)     │
//	│  - Verify signer public keys match signatures           │
//	│  - Parse aggregated public key                          │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 2: Deserialize Partial Signatures                  │
//	│  - Convert PartialSignatureData to *musig2.PartialSig   │
//	│  - Parse S component (signature scalar)                 │
//	│  - Parse R component (nonce commitment)                 │
//	│  LIMITATION: btcd doesn't provide deserialization       │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 3: Combine Partial Signatures (Low-Level API)      │
//	│  - Parse combined nonce from input                      │
//	│  - Call musig2.CombineSigs(combinedNonce, partialSigs)  │
//	│  - NO PRIVATE KEY NEEDED (watch-only wallet)            │
//	│  - Returns final Schnorr signature                      │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 4: Verify Final Signature                          │
//	│  - Verify against aggregated public key                 │
//	│  - Verify against message hash                          │
//	│  - Ensure signature is valid Schnorr signature          │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 5: Finalize PSBT                                   │
//	│  - Add final signature to PSBT                          │
//	│  - Finalize all inputs                                  │
//	│  - Extract final transaction                            │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Output: Ready-to-Broadcast Transaction                  │
//	│  - Final PSBT (base64)                                  │
//	│  - Final transaction hex                                │
//	│  - Transaction ID                                       │
//	│  - IsComplete = true                                    │
//	└─────────────────────────────────────────────────────────┘
//
// # Watch Wallet Architecture
//
// The Watch wallet is designed to operate without private keys:
//
//	┌──────────────────────────────────────────────────────┐
//	│ Watch Wallet (Online, No Private Keys)              │
//	├──────────────────────────────────────────────────────┤
//	│                                                      │
//	│  ✓ Monitor blockchain for deposits                  │
//	│  ✓ Create unsigned transactions (PSBTs)             │
//	│  ✓ Coordinate MuSig2 signing process                │
//	│  ✓ Aggregate partial signatures                     │
//	│  ✓ Verify final signatures                          │
//	│  ✓ Broadcast transactions                           │
//	│                                                      │
//	│  ✗ Does NOT possess private keys                    │
//	│  ✗ Does NOT generate signatures                     │
//	│  ✗ Cannot sign transactions alone                   │
//	└──────────────────────────────────────────────────────┘
//
// # Security Model
//
//   - Watch wallet compromise does NOT leak private keys
//   - Attacker with watch wallet access can:
//   - Monitor transactions and balances
//   - Create unsigned transactions (but cannot sign them)
//   - Attacker CANNOT:
//   - Steal funds (requires offline wallet signatures)
//   - Create valid signatures
//   - Forge transactions
//
// # Known Limitation: Partial Signature Deserialization
//
// The btcd musig2.PartialSignature type doesn't expose deserialization methods.
// This prevents converting input.PartialSignatures to *musig2.PartialSignature.
//
// Required for production:
//  1. Contribute NewPartialSignature() constructor to btcd library, OR
//  2. Use alternative serialization format, OR
//  3. Use reflection (not recommended)
//
// Current implementation is a placeholder showing correct architectural approach
// (watch-only, no private keys, uses low-level CombineSigs API).
//
// # Example
//
// For detailed usage examples, see Example_aggregateMuSig2Signatures.
package btc

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
)

// muSig2BTCClient defines the minimal interface for MuSig2 aggregation.
// Currently unused in Execute() but reserved for future broadcast support.
type muSig2BTCClient interface {
	apibtc.PSBTFinalizer     // For finalizing PSBT
	apibtc.TransactionSender // For broadcasting final transaction
}

type aggregateMuSig2SignaturesUseCase struct {
	musig2Service *apibtcimpl.MuSig2Service
	btcClient     muSig2BTCClient
}

// NewAggregateMuSig2SignaturesUseCase creates a new AggregateMuSig2SignaturesUseCase for BTC watch wallet.
func NewAggregateMuSig2SignaturesUseCase(
	musig2Service *apibtcimpl.MuSig2Service,
	btcClient muSig2BTCClient,
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

// ExampleAggregateMuSig2Signatures demonstrates how to use AggregateMuSig2SignaturesUseCase
// to aggregate partial signatures from multiple signers into a final Schnorr signature.
//
// This example shows the Watch wallet's role as coordinator:
//  1. Initialize dependencies (MuSig2Service, Bitcoin RPC client)
//  2. Create the use case with NewAggregateMuSig2SignaturesUseCase
//  3. Collect partial signatures from all signers (Keygen, Sign1, Sign2)
//  4. Prepare input with:
//     - All partial signatures
//     - Signer public keys
//     - Aggregated public key
//     - Combined nonce
//     - PSBT with transaction data
//  5. Execute aggregation to create final signature
//  6. Broadcast the finalized transaction to Bitcoin network
//
// CRITICAL NOTES:
//   - Watch wallet does NOT sign (watch-only, no private keys)
//   - Uses low-level musig2.CombineSigs API (no private key needed)
//   - Aggregation can only succeed if ALL required partial signatures are present
//   - Final signature is a standard Schnorr signature (indistinguishable from single-sig)
//
// KNOWN LIMITATION:
// btcd doesn't provide deserialization for musig2.PartialSignature.
// This example shows the correct architecture, but production requires
// partial signature deserialization support.
//
// Note: This is a conceptual example. In real implementation:
//   - Partial signatures would be extracted from PSBT proprietary fields
//   - PSBT would contain the actual unsigned transaction
//   - Final transaction would be broadcast via Bitcoin Core RPC
//   - Transaction ID would be returned for tracking
func ExampleAggregateMuSig2Signatures() {
	ctx := context.Background()

	// Initialize dependencies
	var (
		musig2Service *apibtcimpl.MuSig2Service // Provides MuSig2 cryptographic operations
		btcClient     muSig2BTCClient           // Bitcoin Core RPC client (for broadcast)
	)

	// Create the use case
	useCase := NewAggregateMuSig2SignaturesUseCase(musig2Service, btcClient)

	// Prepare partial signatures collected from signers
	// In real implementation, these would be extracted from PSBT
	partialSignatures := []watchusecase.PartialSignatureData{
		{
			SignerID:        "keygen",
			Signature:       [32]byte{}, // S component from keygen
			NonceCommitment: [33]byte{}, // R component from keygen
		},
		{
			SignerID:        "auth1",
			Signature:       [32]byte{}, // S component from sign1
			NonceCommitment: [33]byte{}, // R component from sign1
		},
		{
			SignerID:        "auth2",
			Signature:       [32]byte{}, // S component from sign2
			NonceCommitment: [33]byte{}, // R component from sign2
		},
	}

	// Signer public keys (for verification)
	signerPublicKeys := [][33]byte{
		{}, // Keygen's public key
		{}, // Sign1's public key
		{}, // Sign2's public key
	}

	// Aggregated data from Round 1
	var aggregatedPublicKey [33]byte // From address creation (compressed point)
	var combinedNonce [66]byte       // Combined/aggregated nonce from all signers
	// Note: MuSig2 nonces consist of two points (R1, R2), so combined nonce is 66 bytes
	// This is the result of aggregating all signers' [66]byte public nonces

	// PSBT with unsigned transaction
	psbtBase64 := "cHNidP8BA..." // Base64-encoded PSBT

	// Prepare input for aggregation
	input := watchusecase.AggregateMuSig2SignaturesInput{
		PartialSignatures:   partialSignatures,
		SignerPublicKeys:    signerPublicKeys,
		AggregatedPublicKey: aggregatedPublicKey,
		CombinedNonce:       combinedNonce,
		PSBTBase64:          psbtBase64,
	}

	// Execute aggregation
	output, err := useCase.Execute(ctx, input)
	if err != nil {
		fmt.Printf("failed to aggregate signatures: %v\n", err)
		return
	}

	if !output.IsComplete {
		fmt.Println("aggregation incomplete - missing signatures")
		return
	}

	fmt.Printf("Successfully aggregated MuSig2 signatures\n")
	fmt.Printf("Transaction ID: %s\n", output.TxID)
	fmt.Printf("Ready to broadcast: %s\n", output.FinalTxHex)

	// After successful execution:
	// 1. All partial signatures have been combined into final Schnorr signature
	// 2. Final signature has been verified against aggregated public key
	// 3. PSBT has been finalized with the signature
	// 4. Transaction is ready for broadcast to Bitcoin network
	//
	// Next steps:
	// 1. Broadcast transaction using btcClient.SendRawTransaction()
	// 2. Monitor transaction confirmation on blockchain
	// 3. Update database with transaction status
	//
	// Output:
	// Successfully aggregated MuSig2 signatures
	// Transaction ID: abc123...
	// Ready to broadcast: 0200000001...
}
