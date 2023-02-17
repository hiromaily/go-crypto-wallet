package account

// AccountType utilization purpose of address
type AccountType string

// account_type
const (
	AccountTypeClient  AccountType = "client"  // users who created address
	AccountTypeDeposit AccountType = "deposit" // this address is used as receiver for deposit
	AccountTypePayment AccountType = "payment" // this address is used as sender for payment
	AccountTypeStored  AccountType = "stored"  // this address is used as receiver to store huge amount of coin

	AccountTypeAuthorization AccountType = "auth"   // authorization account for multisig address
	AccountTypeAuth1         AccountType = "auth1"  // it's used for key creation
	AccountTypeAuth2         AccountType = "auth2"  // it's used for key creation
	AccountTypeAuth3         AccountType = "auth3"  // it's used for key creation
	AccountTypeAuth4         AccountType = "auth4"  // it's used for key creation
	AccountTypeAuth5         AccountType = "auth5"  // it's used for key creation
	AccountTypeAuth6         AccountType = "auth6"  // it's used for key creation
	AccountTypeAuth7         AccountType = "auth7"  // it's used for key creation
	AccountTypeAuth8         AccountType = "auth8"  // it's used for key creation
	AccountTypeAuth9         AccountType = "auth9"  // it's used for key creation
	AccountTypeAuth10        AccountType = "auth10" // it's used for key creation
	AccountTypeAuth11        AccountType = "auth11" // it's used for key creation
	AccountTypeAuth12        AccountType = "auth12" // it's used for key creation
	AccountTypeAuth13        AccountType = "auth13" // it's used for key creation
	AccountTypeAuth14        AccountType = "auth14" // it's used for key creation
	AccountTypeAuth15        AccountType = "auth15" // it's used for key creation

	AccountTypeAnonymous AccountType = "anonymous" // payment receiver account
	AccountTypeTest      AccountType = "test"      // unittest only
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

// ValidateAccountType validate AccountType
func ValidateAccountType(v string) bool {
	if _, ok := AccountTypeMap[v]; ok {
		return true
	}
	return false
}

// AccountTypeValue account_type value
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

// Uint32 converter
func (a AccountType) Uint32() uint32 {
	return AccountTypeValue[a]
}

//----------------------------------------------------
// AuthType
//----------------------------------------------------

// AuthType is for authorization account details
//
//	this account is used for authorization of multisig address on sing wallet
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

// ValidateAuthType validate AuthType
func ValidateAuthType(v string) bool {
	if _, ok := AuthTypeMap[v]; ok {
		return true
	}
	return false
}

// AuthTypeMap auth_type mapper
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

// String converter
func (a AuthType) String() string {
	return string(a)
}

// AccountType converter
func (a AuthType) AccountType() AccountType {
	return AccountTypeMap[a.String()]
}
