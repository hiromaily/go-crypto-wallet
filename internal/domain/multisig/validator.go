package multisig

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// ValidateMultisigConfig validates the multisig configuration (M-of-N).
// requiredSigs (M) is the number of signatures required.
// totalSigs (N) is the total number of possible signers.
func ValidateMultisigConfig(requiredSigs, totalSigs int) error {
	if requiredSigs < 1 {
		return fmt.Errorf("required signatures must be at least 1: got %d", requiredSigs)
	}

	if totalSigs < 2 {
		return fmt.Errorf("multisig requires at least 2 total signers: got %d", totalSigs)
	}

	if requiredSigs > totalSigs {
		return fmt.Errorf("required signatures (%d) cannot exceed total signatures (%d)", requiredSigs, totalSigs)
	}

	// Practical limit check
	if totalSigs > 15 {
		return fmt.Errorf("total signatures cannot exceed 15: got %d", totalSigs)
	}

	return nil
}

// ValidateRedeemScript validates that a redeem script is not empty.
// Note: Full cryptographic validation would require infrastructure dependencies,
// so this is a basic domain-level check.
func ValidateRedeemScript(redeemScript string) error {
	if redeemScript == "" {
		return errors.New("redeem script cannot be empty")
	}
	// Length check: redeem scripts are typically hex-encoded
	if len(redeemScript) < 20 {
		return fmt.Errorf("redeem script appears too short: %d characters", len(redeemScript))
	}
	return nil
}

// ValidatePublicKeys validates that the required number of public keys are provided.
func ValidatePublicKeys(publicKeys []string, requiredCount int) error {
	if len(publicKeys) < requiredCount {
		return fmt.Errorf("insufficient public keys: got %d, need %d", len(publicKeys), requiredCount)
	}

	for i, pk := range publicKeys {
		if pk == "" {
			return fmt.Errorf("public key at index %d is empty", i)
		}
		// Basic length check for hex-encoded public keys
		if len(pk) < 66 {
			return fmt.Errorf("public key at index %d appears too short: %d characters", i, len(pk))
		}
	}

	return nil
}

// ValidateAccountForMultisig validates that an account type can use multisig.
func ValidateAccountForMultisig(accountType account.AccountType) error {
	if !account.IsMultisigEligible(accountType) {
		return fmt.Errorf("account type %s is not eligible for multisig addresses", accountType)
	}
	return nil
}

// ValidateAuthAccounts validates that all auth accounts are proper authorization accounts.
func ValidateAuthAccounts(authAccounts []account.AccountType) error {
	if len(authAccounts) == 0 {
		return errors.New("at least one auth account is required")
	}

	for i, acct := range authAccounts {
		if !account.IsAuthorizationAccount(acct) {
			return fmt.Errorf("account at index %d (%s) is not an authorization account", i, acct)
		}
	}

	return nil
}
