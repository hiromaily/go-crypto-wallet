package xrppublic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

const signerListSetBaseFee = 12 // minimum fee in drops per XRPL protocol

// signerListSetTxInput is the infrastructure wire type for a SignerListSet transaction.
// JSON tags match the XRPL protocol wire format.
type signerListSetTxInput struct {
	TransactionType    string            `json:"TransactionType"`
	Account            string            `json:"Account"`
	SignerQuorum       uint32            `json:"SignerQuorum"`
	SignerEntries      []signerEntryWire `json:"SignerEntries"`
	Fee                string            `json:"Fee"`
	Flags              uint64            `json:"Flags"`
	LastLedgerSequence uint64            `json:"LastLedgerSequence"`
	Sequence           uint64            `json:"Sequence"`
}

// signerEntryWire wraps the inner signer entry to produce the nested JSON structure
// required by XRPL: {"SignerEntry": {"Account": "...", "SignerWeight": N}}.
type signerEntryWire struct {
	SignerEntry signerEntryInner `json:"SignerEntry"`
}

// signerEntryInner holds the actual signer account and weight.
type signerEntryInner struct {
	Account      string `json:"Account"`
	SignerWeight uint32 `json:"SignerWeight"`
}

// toDTOSignerListSetTxInput converts the infrastructure wire type to the application DTO.
func toDTOSignerListSetTxInput(infra *signerListSetTxInput) *dtoxrp.SignerListSetTxInput {
	if infra == nil {
		return nil
	}
	entries := make([]dtoxrp.SignerListEntry, len(infra.SignerEntries))
	for i, e := range infra.SignerEntries {
		entries[i] = dtoxrp.SignerListEntry{
			Account:      e.SignerEntry.Account,
			SignerWeight: e.SignerEntry.SignerWeight,
		}
	}
	return &dtoxrp.SignerListSetTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		SignerQuorum:       infra.SignerQuorum,
		SignerEntries:      entries,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
	}
}

// PrepareSignerListSetTransaction builds an unsigned SignerListSet transaction using
// WebSocket account_info to fetch the current sequence and ledger index.
//
// Fee is calculated as (N+1) * 12 drops where N is the number of signer entries,
// matching the XRPL protocol minimum for SignerListSet transactions.
// Instructions can override fee, sequence, and last ledger sequence.
func (p *PublicXRP) PrepareSignerListSetTransaction(
	ctx context.Context,
	senderAccount string,
	signerQuorum uint32,
	signerEntries []apixrp.SignerEntryInput,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.SignerListSetTxInput, string, error) {
	accInfo, err := p.PublicRPC.AccountInfo(ctx, senderAccount)
	if err != nil {
		return nil, "", fmt.Errorf("failed to call PrepareSignerListSetTransaction(): %w", err)
	}
	if accInfo.Error != "" {
		return nil, "", fmt.Errorf("failed to call PrepareSignerListSetTransaction(): %s", accInfo.Error)
	}

	sequence := uint64(accInfo.Result.AccountData.Sequence)
	lastLedgerSequence := uint64(accInfo.Result.LedgerCurrentIndex) + xrpkg.MaxLedgerVersionOffset

	// Fee: (N+1) * baseFee drops, where N is the number of signer entries.
	fee := strconv.Itoa((len(signerEntries) + 1) * signerListSetBaseFee)

	if instructions != nil {
		if instructions.Fee != "" {
			fee = instructions.Fee
		}
		if instructions.Sequence != 0 {
			sequence = instructions.Sequence
		}
		if instructions.MaxLedgerVersion != 0 {
			lastLedgerSequence = instructions.MaxLedgerVersion
		} else if instructions.MaxLedgerVersionOffset != 0 {
			lastLedgerSequence = uint64(accInfo.Result.LedgerCurrentIndex) + instructions.MaxLedgerVersionOffset
		}
	}

	wireEntries := make([]signerEntryWire, len(signerEntries))
	for i, e := range signerEntries {
		wireEntries[i] = signerEntryWire{
			SignerEntry: signerEntryInner{
				Account:      e.Account,
				SignerWeight: e.Weight,
			},
		}
	}

	txInput := &signerListSetTxInput{
		TransactionType:    "SignerListSet",
		Account:            senderAccount,
		SignerQuorum:       signerQuorum,
		SignerEntries:      wireEntries,
		Fee:                fee,
		Flags:              0,
		LastLedgerSequence: lastLedgerSequence,
		Sequence:           sequence,
	}

	jsonBytes, err := json.Marshal(txInput)
	if err != nil {
		return nil, "", fmt.Errorf("failed to call json.Marshal(txInput): %w", err)
	}

	logger.Debug("PrepareSignerListSetTransaction", "TxJSON", string(jsonBytes))

	return toDTOSignerListSetTxInput(txInput), string(jsonBytes), nil
}
