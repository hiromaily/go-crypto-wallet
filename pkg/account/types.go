package account

// AccountType utilization purpose of address
type AccountType string

// account_type
const (
	AccountTypeClient        AccountType = "client"    // users who created address
	AccountTypeDeposit       AccountType = "deposit"   // this address is used as receiver for deposit
	AccountTypePayment       AccountType = "payment"   // this address is used as sender for payment
	AccountTypeFee           AccountType = "fee"       // for transaction fee
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
	"fee":       AccountTypeFee,
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
	AccountTypeFee:           3,
	AccountTypeStored:        4,
	AccountTypeAuthorization: 5,
	AccountTypeAnonymous:     99,
	AccountTypeTest:          100,
}

// AccountTypeMultisig true: account type is for multisig address
var AccountTypeMultisig = map[AccountType]bool{
	AccountTypeClient:        false,
	AccountTypeDeposit:       true,
	AccountTypePayment:       true,
	AccountTypeFee:           true,
	AccountTypeStored:        true,
	AccountTypeAuthorization: false,
}
