package account

// AccountType represents the utilization purpose of an address in the wallet system.
//
// The system uses different account types for various operational purposes:
//   - Client: User-created addresses for receiving deposits
//   - Deposit: Aggregates coins from client accounts
//   - Payment: Sends payments to external addresses
//   - Stored: Long-term cold storage for large amounts
//   - Authorization (auth1-auth15): Multisig signers
//   - Anonymous: External payment receiver addresses
//   - Test: Testing purposes only
type AccountType string

// Account type constants
const (
	// AccountTypeClient represents addresses created for users who deposit funds
	AccountTypeClient AccountType = "client"

	// AccountTypeDeposit is used as receiver for aggregating deposits from client accounts
	AccountTypeDeposit AccountType = "deposit"

	// AccountTypePayment is used as sender for making payments to external addresses
	AccountTypePayment AccountType = "payment"

	// AccountTypeStored is used as receiver to store large amounts of cryptocurrency
	AccountTypeStored AccountType = "stored"

	// AccountTypeAuthorization is a generic authorization account for multisig addresses
	AccountTypeAuthorization AccountType = "auth"

	// AccountTypeAuth1-15 are specific authorization accounts used for key creation
	// and multisig signing (up to 15 different operators/signers)
	AccountTypeAuth1  AccountType = "auth1"
	AccountTypeAuth2  AccountType = "auth2"
	AccountTypeAuth3  AccountType = "auth3"
	AccountTypeAuth4  AccountType = "auth4"
	AccountTypeAuth5  AccountType = "auth5"
	AccountTypeAuth6  AccountType = "auth6"
	AccountTypeAuth7  AccountType = "auth7"
	AccountTypeAuth8  AccountType = "auth8"
	AccountTypeAuth9  AccountType = "auth9"
	AccountTypeAuth10 AccountType = "auth10"
	AccountTypeAuth11 AccountType = "auth11"
	AccountTypeAuth12 AccountType = "auth12"
	AccountTypeAuth13 AccountType = "auth13"
	AccountTypeAuth14 AccountType = "auth14"
	AccountTypeAuth15 AccountType = "auth15"

	// AccountTypeAnonymous represents external payment receiver addresses
	AccountTypeAnonymous AccountType = "anonymous"

	// AccountTypeTest is used for unit testing only
	AccountTypeTest AccountType = "test"
)

// String returns the string representation of the account type.
func (a AccountType) String() string {
	return string(a)
}

// Is compares the account type with the given string value.
func (a AccountType) Is(v string) bool {
	return a.String() == v
}

// Uint32 returns the numeric value of the account type.
func (a AccountType) Uint32() uint32 {
	return AccountTypeValue[a]
}

// Allow returns true if the given account string is in the allowed list.
func Allow(acnt string, accountList []AccountType) bool {
	if !ValidateAccountType(acnt) {
		return false
	}
	for _, v := range accountList {
		if acnt == v.String() {
			return true
		}
	}
	return false
}

// NotAllow returns true if the given account string is not in the disallowed list.
func NotAllow(acnt string, accountList []AccountType) bool {
	if !ValidateAccountType(acnt) {
		return false
	}
	for _, v := range accountList {
		if acnt == v.String() {
			return false
		}
	}
	return true
}

// AccountTypeMap provides string to AccountType mapping for validation and conversion.
var AccountTypeMap = map[string]AccountType{
	"client":  AccountTypeClient,
	"deposit": AccountTypeDeposit,
	"payment": AccountTypePayment,
	"stored":  AccountTypeStored,

	"auth":   AccountTypeAuthorization,
	"auth1":  AccountTypeAuth1,
	"auth2":  AccountTypeAuth2,
	"auth3":  AccountTypeAuth3,
	"auth4":  AccountTypeAuth4,
	"auth5":  AccountTypeAuth5,
	"auth6":  AccountTypeAuth6,
	"auth7":  AccountTypeAuth7,
	"auth8":  AccountTypeAuth8,
	"auth9":  AccountTypeAuth9,
	"auth10": AccountTypeAuth10,
	"auth11": AccountTypeAuth11,
	"auth12": AccountTypeAuth12,
	"auth13": AccountTypeAuth13,
	"auth14": AccountTypeAuth14,
	"auth15": AccountTypeAuth15,

	"anonymous": AccountTypeAnonymous,
	"test":      AccountTypeTest,
}

// ValidateAccountType validates that the given string is a valid account type.
func ValidateAccountType(v string) bool {
	_, ok := AccountTypeMap[v]
	return ok
}

// AccountTypeValue provides numeric values for account types.
// These values are used for database storage and ordering.
// Note: if key is not found, value is 0
var AccountTypeValue = map[AccountType]uint32{
	AccountTypeClient:  0,
	AccountTypeDeposit: 1,
	AccountTypePayment: 2,
	AccountTypeStored:  3,

	AccountTypeAuthorization: 10,
	AccountTypeAuth1:         11,
	AccountTypeAuth2:         12,
	AccountTypeAuth3:         13,
	AccountTypeAuth4:         14,
	AccountTypeAuth5:         15,
	AccountTypeAuth6:         16,
	AccountTypeAuth7:         17,
	AccountTypeAuth8:         18,
	AccountTypeAuth9:         19,
	AccountTypeAuth10:        20,
	AccountTypeAuth11:        21,
	AccountTypeAuth12:        22,
	AccountTypeAuth13:        23,
	AccountTypeAuth14:        24,
	AccountTypeAuth15:        25,

	AccountTypeAnonymous: 99,
	AccountTypeTest:      100,
}

// AuthType represents authorization account details for multisig addresses.
//
// This type is used for operators with sign wallet to authorize multisig transactions.
type AuthType string

// Authorization type constants for multisig signing
const (
	// AuthType1-15 represent different operators who own sign wallets
	// Each operator has their own auth account for signing
	AuthType1  AuthType = "auth1"
	AuthType2  AuthType = "auth2"
	AuthType3  AuthType = "auth3"
	AuthType4  AuthType = "auth4"
	AuthType5  AuthType = "auth5"
	AuthType6  AuthType = "auth6"
	AuthType7  AuthType = "auth7"
	AuthType8  AuthType = "auth8"
	AuthType9  AuthType = "auth9"
	AuthType10 AuthType = "auth10"
	AuthType11 AuthType = "auth11"
	AuthType12 AuthType = "auth12"
	AuthType13 AuthType = "auth13"
	AuthType14 AuthType = "auth14"
	AuthType15 AuthType = "auth15"
)

// String returns the string representation of the auth type.
func (a AuthType) String() string {
	return string(a)
}

// AccountType converts AuthType to AccountType.
func (a AuthType) AccountType() AccountType {
	return AccountTypeMap[a.String()]
}

// ValidateAuthType validates that the given string is a valid auth type.
func ValidateAuthType(v string) bool {
	_, ok := AuthTypeMap[v]
	return ok
}

// AuthTypeMap provides string to AuthType mapping for validation and conversion.
var AuthTypeMap = map[string]AuthType{
	"auth1":  AuthType1,
	"auth2":  AuthType2,
	"auth3":  AuthType3,
	"auth4":  AuthType4,
	"auth5":  AuthType5,
	"auth6":  AuthType6,
	"auth7":  AuthType7,
	"auth8":  AuthType8,
	"auth9":  AuthType9,
	"auth10": AuthType10,
	"auth11": AuthType11,
	"auth12": AuthType12,
	"auth13": AuthType13,
	"auth14": AuthType14,
	"auth15": AuthType15,
}
