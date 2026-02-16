package bchutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// PrevTx represents previous transaction data for BCH signing.
// This struct is shared between keygen and sign wallets for multisig transactions.
//
// In BCH multisig workflows, transaction files include prevTx metadata to enable
// subsequent signers to add their signatures. This is especially important due to
// the BCH Node RPC bug where signrawtransactionwithkey incorrectly reports
// complete=true for partially signed multisig transactions.
//
// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
type PrevTx struct {
	TxID         string
	Vout         uint32
	ScriptPubKey string
	RedeemScript string
	Amount       int64
}

// ParseRawTxContent parses BCH raw transaction content.
//
// Format:
//
//	Line 1: Transaction hex
//	Line 2+: prevtx:txid:vout:scriptPubKey:redeemScript:amount
//
// The prevTx metadata is required for multisig transactions to enable subsequent
// signers to add their signatures. This metadata is preserved as a workaround for
// the BCH Node RPC bug.
//
// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
func ParseRawTxContent(content string) (string, []PrevTx, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return "", nil, errors.New("empty transaction content")
	}

	txHex := lines[0]
	prevTxs := make([]PrevTx, 0, len(lines)-1)

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || !strings.HasPrefix(line, "prevtx:") {
			continue
		}

		// Parse prevtx:txid:vout:scriptPubKey:redeemScript:amount
		parts := strings.Split(strings.TrimPrefix(line, "prevtx:"), ":")
		if len(parts) < 5 {
			logger.Warn("invalid prevtx format, skipping", "line", line)
			continue
		}

		vout64, err := strconv.ParseUint(parts[1], 10, 32)
		if err != nil {
			logger.Warn("invalid vout in prevtx", "vout", parts[1], "error", err)
			continue
		}
		vout := uint32(vout64)

		amount, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			logger.Warn("invalid amount in prevtx", "amount", parts[4], "error", err)
			continue
		}

		prevTxs = append(prevTxs, PrevTx{
			TxID:         parts[0],
			Vout:         vout,
			ScriptPubKey: parts[2],
			RedeemScript: parts[3],
			Amount:       amount,
		})
	}

	return txHex, prevTxs, nil
}

// ConvertPrevTxsToDTO converts PrevTx to DTO format for BCH signing operations.
//
// This conversion is needed when calling the BCH Node RPC signrawtransactionwithkey
// method, which requires prevTx information in a specific DTO format.
func ConvertPrevTxsToDTO(prevTxs []PrevTx) []dtobtc.PreviousTx {
	result := make([]dtobtc.PreviousTx, len(prevTxs))
	for i, prev := range prevTxs {
		result[i] = dtobtc.PreviousTx{
			TxID:         prev.TxID,
			Vout:         prev.Vout,
			ScriptPubKey: prev.ScriptPubKey,
			RedeemScript: prev.RedeemScript,
			Amount:       btcutil.Amount(prev.Amount),
		}
	}
	return result
}

// FormatSignedTxContent formats the signed transaction content for BCH.
//
// This function creates the transaction file format used in BCH multisig workflows:
//
//	Line 1: Transaction hex
//	Line 2+: prevtx:txid:vout:scriptPubKey:redeemScript:amount
//
// The prevTx metadata is included to enable subsequent signers (Sign1, Sign2) to
// add their signatures. This is critical for the BCH Node RPC bug workaround.
//
// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
func FormatSignedTxContent(txHex string, prevTxs []PrevTx) string {
	var builder strings.Builder
	builder.WriteString(txHex)
	builder.WriteString("\n")

	for _, prev := range prevTxs {
		fmt.Fprintf(&builder, "prevtx:%s:%d:%s:%s:%d\n",
			prev.TxID,
			prev.Vout,
			prev.ScriptPubKey,
			prev.RedeemScript,
			prev.Amount,
		)
	}
	return builder.String()
}

// ShouldIncludePrevTxMetadata determines whether prevTx metadata should be included
// in the transaction file based on transaction state and multisig configuration.
//
// Workaround for BCH Node RPC bug:
//
// Issue: BCH Node's signrawtransactionwithkey incorrectly reports complete=true
// after adding the 2nd signature to a 3-of-3 multisig transaction.
//
// Solution: For multisig transactions, always include prevTx metadata regardless
// of the isSigned flag. For single-sig transactions, only include metadata when
// the transaction is not yet fully signed.
//
// Parameters:
//   - isSigned: Whether BCH Node reports the transaction as complete
//   - isMultisig: Whether the transaction uses multisig (detected via non-nil multisigAccount)
//
// Returns: true if prevTx metadata should be included in the transaction file
//
// Reference: docs/chains/bch/README.md - "Known Issues and Workarounds"
// Related: Issue #485, Issue #433
func ShouldIncludePrevTxMetadata(isSigned bool, isMultisig bool) bool {
	// For multisig transactions, always include prevTx metadata for subsequent signers,
	// even if BCH incorrectly reports isSigned=true for partial signatures
	if isMultisig {
		return true
	}

	// For single-sig transactions, only include metadata when not yet fully signed
	return !isSigned
}
