package account

// MultisigConfig represents multisig account configuration.
// This is a domain entity that holds multisig account settings.
type MultisigConfig struct {
	AccountMap map[AccountType]map[int][]AuthType
}

// IsMultisigAccount checks if the account type is a multisig account.
func (m *MultisigConfig) IsMultisigAccount(accountType AccountType) bool {
	if m == nil {
		return false
	}
	_, ok := m.AccountMap[accountType]
	return ok
}

// MultiAccounts returns the multisig configuration map.
func (m *MultisigConfig) MultiAccounts() map[AccountType]map[int][]AuthType {
	if m == nil {
		return nil
	}
	return m.AccountMap
}

// NewMultisigConfig creates a MultisigConfig domain entity from a map structure.
//
// This function is placed in the domain layer because:
//
// 1. Domain entity factory pattern:
//   - This is a factory function that creates a domain entity (MultisigConfig)
//   - Factory functions for domain entities belong in the domain layer
//   - It's conceptually part of the domain entity's lifecycle management
//
// 2. Co-location with domain entity:
//   - MultisigConfig and its factory function are closely related
//   - Having them in the same package improves code discoverability and understanding
//   - Developers looking for MultisigConfig will naturally find its factory function
//
// 3. Pure domain types:
//   - This function accepts only domain types (AccountType, AuthType)
//   - It has no dependencies on infrastructure or configuration packages
//   - This maintains the domain layer's independence from infrastructure
//
// 4. Separation of concerns:
//   - Domain layer: Domain entity creation using domain types (Domain concern)
//   - Infrastructure layer: Conversion from config structures to domain types (Infrastructure concern)
//   - The infrastructure layer calls this function after converting config.AccountMultisig to domain types
//
// The function accepts a map structure that represents multisig configuration:
//   - Key: AccountType (which account type is multisig)
//   - Value: map[int][]AuthType (required signatures -> list of auth types)
//
// This allows the domain layer to create MultisigConfig without depending on
// configuration structures, while infrastructure layer handles the conversion.
func NewMultisigConfig(accountMap map[AccountType]map[int][]AuthType) *MultisigConfig {
	if accountMap == nil {
		return &MultisigConfig{
			AccountMap: make(map[AccountType]map[int][]AuthType),
		}
	}
	return &MultisigConfig{
		AccountMap: accountMap,
	}
}
