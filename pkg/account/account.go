package account

//AccountType 利用目的
type AccountType string

// account_type
const (
	AccountTypeClient        AccountType = "client"        //ユーザーの入金受付用アドレス
	AccountTypeReceipt       AccountType = "receipt"       //入金を受け付けるアドレス用
	AccountTypePayment       AccountType = "payment"       //出金時に支払いをするアドレス
	AccountTypeQuoine        AccountType = "quoine"        //Quoineから購入したcoinが入金されるであろうアドレス
	AccountTypeFee           AccountType = "fee"           //手数料保管用アドレス
	AccountTypeStored        AccountType = "stored"        //保管用アドレス(多額のコインはこちらに保管しておく
	AccountTypeAuthorization AccountType = "authorization" //マルチシグアドレスのための承認アドレス
)

func (a AccountType) String() string {
	return string(a)
}

func (a AccountType) Is(v string) bool {
	return a.String() == v
}

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

//AccountTypeMap account_type mapper
var AccountTypeMap = map[string]AccountType{
	"client":        AccountTypeClient,
	"receipt":       AccountTypeReceipt,
	"payment":       AccountTypePayment,
	"quoine":        AccountTypeQuoine,
	"fee":           AccountTypeFee,
	"stored":        AccountTypeStored,
	"authorization": AccountTypeAuthorization,
}

// ValidateAccountType AccountTypeのバリデーションを行う
func ValidateAccountType(v string) bool {
	if _, ok := AccountTypeMap[v]; ok {
		return true
	}
	return false
}

//AccountTypeValue account_type value
var AccountTypeValue = map[AccountType]uint8{
	AccountTypeClient:        0,
	AccountTypeReceipt:       1,
	AccountTypePayment:       2,
	AccountTypeQuoine:        3,
	AccountTypeFee:           4,
	AccountTypeStored:        5,
	AccountTypeAuthorization: 6,
}

//AccountTypeMultisig account_type毎のmultisig対応アカウントかどうか
var AccountTypeMultisig = map[AccountType]bool{
	AccountTypeClient:        false,
	AccountTypeReceipt:       true,
	AccountTypePayment:       true,
	AccountTypeQuoine:        true,
	AccountTypeFee:           true,
	AccountTypeStored:        true,
	AccountTypeAuthorization: false,
}
