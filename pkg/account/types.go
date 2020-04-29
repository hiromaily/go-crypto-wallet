package account

// AccountType utilization purpose of address
type AccountType string

// account_type
const (
	AccountTypeClient        AccountType = "client"    // users who created address
	AccountTypeDeposit       AccountType = "deposit"   // this address is used as receiver for deposit
	AccountTypePayment       AccountType = "payment"   // this address is used as sender for payment
	AccountTypeStored        AccountType = "stored"    // this address is used as receiver to store huge amount of coin
	AccountTypeAuthorization AccountType = "auth"      // authorization account for multisig address
	AccountTypeAnonymous     AccountType = "anonymous" // payment receiver account
	AccountTypeTest          AccountType = "test"      // unittest only
)

// String converter
func (a AccountType) String() string {
	return string(a)
}

// Is compare with params
func (a AccountType) Is(v string) bool {
	return a.String() == v
}

// Allow return true if acnt is in list
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

// NotAllow return true if acnt is not in list
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

// AccountTypeMap account_type mapper
var AccountTypeMap = map[string]AccountType{
	"client":    AccountTypeClient,
	"deposit":   AccountTypeDeposit,
	"payment":   AccountTypePayment,
	"stored":    AccountTypeStored,
	"auth":      AccountTypeAuthorization,
	"anonymous": AccountTypeAnonymous,
	"test":      AccountTypeTest,
}

// ValidateAccountType validate AccountType
func ValidateAccountType(v string) bool {
	if _, ok := AccountTypeMap[v]; ok {
		return true
	}
	return false
}

// AccountTypeValue account_type value
var AccountTypeValue = map[AccountType]uint32{
	AccountTypeClient:        0,
	AccountTypeDeposit:       1,
	AccountTypePayment:       2,
	AccountTypeStored:        3,
	AccountTypeAuthorization: 4,
	AccountTypeAnonymous:     99,
	AccountTypeTest:          100,
}

// AccountTypeMultisig true: account type is for multisig address
var AccountTypeMultisig = map[AccountType]bool{
	AccountTypeClient:        false,
	AccountTypeDeposit:       true,
	AccountTypePayment:       true,
	AccountTypeStored:        true,
	AccountTypeAuthorization: false,
}

//----------------------------------------------------
// AuthType
//----------------------------------------------------

// AuthType is for authorization account details
//  this account is used for authorization of multisig address on sing wallet
type AuthType string

// auth_account_type, this type is used for operators with sign wallet
const (
	AuthType1  AuthType = "auth1" // operator1 would own sign wallet with this account
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

// String converter
func (a AuthType) String() string {
	return string(a)
}
