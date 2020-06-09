package account

// Multisig address involved accounts

// MultisigAccounter is AddressRepository interface
type MultisigAccounter interface {
	IsMultisigAccount(accountType AccountType) bool
	MultiAccounts() map[AccountType]map[int][]AuthType
}

type multisigAccount struct {
	accountMap map[AccountType]map[int][]AuthType
}

// NewMultisigAccounts returns multisigAccount instance
func NewMultisigAccounts(confMultisig []AccountMultisig) MultisigAccounter {
	multisigAccount := multisigAccount{
		accountMap: make(map[AccountType]map[int][]AuthType, len(confMultisig)),
	}
	for _, val := range confMultisig {
		multisigAccount.accountMap[val.Type] = map[int][]AuthType{
			val.Required: val.AuthUsers,
		}
	}
	return &multisigAccount
}

// IsMultisigAccount validates Multisig account or not
func (m *multisigAccount) IsMultisigAccount(v AccountType) bool {
	if _, ok := m.accountMap[v]; ok {
		return true
	}
	return false
}

// MultiAccounts returns accountMap
func (m *multisigAccount) MultiAccounts() map[AccountType]map[int][]AuthType {
	return m.accountMap
}

// MultisigAccounts proportion of multisig address M:N
//var MultisigAccounts = map[AccountType]map[int][]AuthType{
//	AccountTypeDeposit: { //2:5+1
//		2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
//	},
//	AccountTypePayment: { //3:5+1
//		3: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
//	},
//	AccountTypeStored: { //4:5+1
//		4: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
//	},
//}
//var MultisigAccounts = map[AccountType]map[int][]AuthType{}
