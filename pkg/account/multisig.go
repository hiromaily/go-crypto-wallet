package account

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
)

// Multisig address involved accounts

// MultisigAccounter is AddressRepository interface
type MultisigAccounter interface {
	IsMultisigAccount(accountType domainAccount.AccountType) bool
	MultiAccounts() map[domainAccount.AccountType]map[int][]domainAccount.AuthType
}

type multisigAccount struct {
	accountMap map[domainAccount.AccountType]map[int][]domainAccount.AuthType
}

// NewMultisigAccounts returns multisigAccount instance
func NewMultisigAccounts(confMultisig []AccountMultisig) MultisigAccounter {
	ma := multisigAccount{
		accountMap: make(map[domainAccount.AccountType]map[int][]domainAccount.AuthType, len(confMultisig)),
	}
	for _, val := range confMultisig {
		ma.accountMap[val.Type] = map[int][]domainAccount.AuthType{
			val.Required: val.AuthUsers,
		}
	}
	return &ma
}

// IsMultisigAccount validates Multisig account or not
func (m *multisigAccount) IsMultisigAccount(v domainAccount.AccountType) bool {
	if _, ok := m.accountMap[v]; ok {
		return true
	}
	return false
}

// MultiAccounts returns accountMap
func (m *multisigAccount) MultiAccounts() map[domainAccount.AccountType]map[int][]domainAccount.AuthType {
	return m.accountMap
}

// MultisigAccounts proportion of multisig address M:N
// var MultisigAccounts = map[AccountType]map[int][]AuthType{
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
// var MultisigAccounts = map[AccountType]map[int][]AuthType{}
