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
