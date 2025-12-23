package account

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
)

// Deprecated: Use github.com/hiromaily/go-crypto-wallet/pkg/domain/account instead.
// This package provides backward compatibility aliases.

// AccountType utilization purpose of address
// Deprecated: Use domain/account.AccountType
type AccountType = domainAccount.AccountType

// account_type
// Deprecated: Use constants from domain/account package
const (
	AccountTypeClient        = domainAccount.AccountTypeClient
	AccountTypeDeposit       = domainAccount.AccountTypeDeposit
	AccountTypePayment       = domainAccount.AccountTypePayment
	AccountTypeStored        = domainAccount.AccountTypeStored
	AccountTypeAuthorization = domainAccount.AccountTypeAuthorization
	AccountTypeAuth1         = domainAccount.AccountTypeAuth1
	AccountTypeAuth2         = domainAccount.AccountTypeAuth2
	AccountTypeAuth3         = domainAccount.AccountTypeAuth3
	AccountTypeAuth4         = domainAccount.AccountTypeAuth4
	AccountTypeAuth5         = domainAccount.AccountTypeAuth5
	AccountTypeAuth6         = domainAccount.AccountTypeAuth6
	AccountTypeAuth7         = domainAccount.AccountTypeAuth7
	AccountTypeAuth8         = domainAccount.AccountTypeAuth8
	AccountTypeAuth9         = domainAccount.AccountTypeAuth9
	AccountTypeAuth10        = domainAccount.AccountTypeAuth10
	AccountTypeAuth11        = domainAccount.AccountTypeAuth11
	AccountTypeAuth12        = domainAccount.AccountTypeAuth12
	AccountTypeAuth13        = domainAccount.AccountTypeAuth13
	AccountTypeAuth14        = domainAccount.AccountTypeAuth14
	AccountTypeAuth15        = domainAccount.AccountTypeAuth15
	AccountTypeAnonymous     = domainAccount.AccountTypeAnonymous
	AccountTypeTest          = domainAccount.AccountTypeTest
)

// Allow return true if acnt is in list
// Deprecated: Use domain/account.Allow
func Allow(acnt string, accountList []AccountType) bool {
	return domainAccount.Allow(acnt, accountList)
}

// NotAllow return true if acnt is not in list
// Deprecated: Use domain/account.NotAllow
func NotAllow(acnt string, accountList []AccountType) bool {
	return domainAccount.NotAllow(acnt, accountList)
}

// AccountTypeMap account_type mapper
// Deprecated: Use domain/account.AccountTypeMap
var AccountTypeMap = domainAccount.AccountTypeMap

// ValidateAccountType validate AccountType
// Deprecated: Use domain/account.ValidateAccountType
func ValidateAccountType(v string) bool {
	return domainAccount.ValidateAccountType(v)
}

// AccountTypeValue account_type value
// Deprecated: Use domain/account.AccountTypeValue
var AccountTypeValue = domainAccount.AccountTypeValue

//----------------------------------------------------
// AuthType
//----------------------------------------------------

// AuthType is for authorization account details
// Deprecated: Use domain/account.AuthType
type AuthType = domainAccount.AuthType

// auth_account_type, this type is used for operators with sign wallet
// Deprecated: Use constants from domain/account package
const (
	AuthType1  = domainAccount.AuthType1
	AuthType2  = domainAccount.AuthType2
	AuthType3  = domainAccount.AuthType3
	AuthType4  = domainAccount.AuthType4
	AuthType5  = domainAccount.AuthType5
	AuthType6  = domainAccount.AuthType6
	AuthType7  = domainAccount.AuthType7
	AuthType8  = domainAccount.AuthType8
	AuthType9  = domainAccount.AuthType9
	AuthType10 = domainAccount.AuthType10
	AuthType11 = domainAccount.AuthType11
	AuthType12 = domainAccount.AuthType12
	AuthType13 = domainAccount.AuthType13
	AuthType14 = domainAccount.AuthType14
	AuthType15 = domainAccount.AuthType15
)

// ValidateAuthType validate AuthType
// Deprecated: Use domain/account.ValidateAuthType
func ValidateAuthType(v string) bool {
	return domainAccount.ValidateAuthType(v)
}

// AuthTypeMap auth_type mapper
// Deprecated: Use domain/account.AuthTypeMap
var AuthTypeMap = domainAccount.AuthTypeMap
