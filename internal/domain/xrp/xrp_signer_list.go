package xrp

import (
	"errors"
	"time"
)

// XRPSignerList represents a multi-signature signer list configuration for an XRP account.
// A signer list enables multi-signature transactions where multiple parties must sign
// to authorize a transaction.
//
// Reference: https://xrpl.org/docs/concepts/accounts/multi-signing
type XRPSignerList struct {
	ID           int64
	AccountID    string     // XRP account address (r...) that owns this signer list
	SignerQuorum uint32     // Minimum total weight of signatures required
	IsActive     bool       // True if this signer list is currently active on the ledger
	SetTxHash    *string    // Transaction hash of SignerListSet that created/updated this list
	CreatedAt    *time.Time // When the signer list was created
	UpdatedAt    *time.Time // When the signer list was last updated
}

// NewXRPSignerList creates a new XRPSignerList entity.
func NewXRPSignerList(
	accountID string,
	signerQuorum uint32,
) (*XRPSignerList, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}
	if signerQuorum == 0 {
		return nil, errors.New("signer quorum must be greater than 0")
	}

	now := time.Now()
	return &XRPSignerList{
		AccountID:    accountID,
		SignerQuorum: signerQuorum,
		IsActive:     true,
		CreatedAt:    &now,
	}, nil
}

// SetTxHashValue sets the transaction hash that created this signer list.
func (s *XRPSignerList) SetTxHashValue(txHash string) {
	s.SetTxHash = &txHash
}

// Deactivate marks this signer list as inactive.
func (s *XRPSignerList) Deactivate() {
	s.IsActive = false
	now := time.Now()
	s.UpdatedAt = &now
}

// XRPSignerEntry represents a single signer in a multi-signature signer list.
type XRPSignerEntry struct {
	ID            int64
	SignerListID  int64      // Reference to XRPSignerList.ID
	SignerAccount string     // XRP address of the authorized signer (r...)
	SignerWeight  uint32     // Weight of this signer (contributes to quorum)
	CreatedAt     *time.Time // When this signer entry was created
}

// NewXRPSignerEntry creates a new XRPSignerEntry entity.
func NewXRPSignerEntry(
	signerListID int64,
	signerAccount string,
	signerWeight uint32,
) (*XRPSignerEntry, error) {
	if signerListID <= 0 {
		return nil, errors.New("signer list ID must be positive")
	}
	if signerAccount == "" {
		return nil, errors.New("signer account cannot be empty")
	}
	if signerWeight == 0 {
		return nil, errors.New("signer weight must be greater than 0")
	}

	now := time.Now()
	return &XRPSignerEntry{
		SignerListID:  signerListID,
		SignerAccount: signerAccount,
		SignerWeight:  signerWeight,
		CreatedAt:     &now,
	}, nil
}

// SignerEntryInput is a simplified input structure for creating signer entries.
// Used when setting up a new signer list.
type SignerEntryInput struct {
	Account string
	Weight  uint32
}

// ValidateSignerEntries validates a list of signer entries.
// - At least 1 signer required (max 32 per XRPL rules)
// - Total weight must meet or exceed quorum
// - No duplicate accounts
func ValidateSignerEntries(entries []SignerEntryInput, quorum uint32) error {
	if len(entries) == 0 {
		return errors.New("at least one signer entry is required")
	}
	if len(entries) > 32 {
		return errors.New("maximum 32 signer entries allowed")
	}

	var totalWeight uint32
	seen := make(map[string]bool)
	for _, entry := range entries {
		if entry.Account == "" {
			return errors.New("signer account cannot be empty")
		}
		if entry.Weight == 0 {
			return errors.New("signer weight must be greater than 0")
		}
		if seen[entry.Account] {
			return errors.New("duplicate signer account: " + entry.Account)
		}
		seen[entry.Account] = true
		totalWeight += entry.Weight
	}

	if totalWeight < quorum {
		return errors.New("total signer weight must be at least equal to quorum")
	}

	return nil
}
