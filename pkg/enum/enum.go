package enum

//EnvironmentType 実行環境
type EnvironmentType string

// environment
const (
	EnvDev  EnvironmentType = "dev"
	EnvProd EnvironmentType = "prod"
)

//EnvironmentTypeValue env_typeの値
var EnvironmentTypeValue = map[EnvironmentType]uint8{
	EnvDev:  1,
	EnvProd: 2,
}

// ValidateEnvironmentType EnvironmentTypeのバリデーションを行う
func ValidateEnvironmentType(val string) bool {
	if _, ok := EnvironmentTypeValue[EnvironmentType(val)]; ok {
		return true
	}
	return false
}

//NetworkType ネットワーク種別
type NetworkType string

// network type
const (
	NetworkTypeMainNet  NetworkType = "mainnet"
	NetworkTypeTestNet3 NetworkType = "testnet3"
)

//CoinType コインの種類
type CoinType uint32

// coin_type
const (
	CoinTypeBitcoin CoinType = 0 //Bitcoin
	CoinTypeTestnet CoinType = 1 //Testnet
)

//AccountType 利用目的
type AccountType string

// account_type
const (
	AccountTypeClient        AccountType = "client"        //ユーザーの入金受付用アドレス
	AccountTypeReceipt       AccountType = "receipt"       //入金を受け付けるアドレス用
	AccountTypePayment       AccountType = "payment"       //出金時に支払いをするアドレス
	AccountTypeAuthorization AccountType = "authorization" //マルチシグアドレスのための承認アドレス
	AccountTypeQuoine        AccountType = "quoine"        //Quoineから購入したcoinが入金されるであろうアドレス
	AccountTypeFee           AccountType = "fee"           //手数料保管用アドレス
)

//AccountTypeValue account_typeの値
var AccountTypeValue = map[AccountType]uint8{
	AccountTypeClient:        0,
	AccountTypeReceipt:       1,
	AccountTypePayment:       2,
	AccountTypeAuthorization: 3,
	AccountTypeQuoine:        4,
	AccountTypeFee:           5,
}

//AccountTypeMultisig account_type毎のmultisig対応アカウントかどうか
var AccountTypeMultisig = map[AccountType]bool{
	AccountTypeClient:        false,
	AccountTypeReceipt:       true,
	AccountTypePayment:       true,
	AccountTypeAuthorization: false,
	AccountTypeQuoine:        false,
	AccountTypeFee:           false,
}

// ValidateAccountType AccountTypeのバリデーションを行う
func ValidateAccountType(val string) bool {
	if _, ok := AccountTypeValue[AccountType(val)]; ok {
		return true
	}
	return false
}

//KeyStatus Key生成進捗ステータス
type KeyStatus string

// key_status
const (
	KeyStatusGenerated            KeyStatus = "generated"              //hd_walletによってkeyが生成された
	KeyStatusImportprivkey        KeyStatus = "importprivkey"          //importprivkeyが実行された
	KeyStatusPubkeyExported       KeyStatus = "pubkey_exported"        //pubkeyがexportされた(receipt/payment)
	KeyStatusMultiAddressImported KeyStatus = "multi_address_imported" //multiaddがimportされた(receipt/payment)
	KeyStatusAddressExported      KeyStatus = "address_exported"       //addressがexportされた
)

//KeyStatusValue key_statusの値
var KeyStatusValue = map[KeyStatus]uint8{
	KeyStatusGenerated:            0,
	KeyStatusImportprivkey:        1,
	KeyStatusPubkeyExported:       2,
	KeyStatusMultiAddressImported: 3,
	KeyStatusAddressExported:      4,
}

// ValidateKeyStatus KeyStatusのバリデーションを行う
func ValidateKeyStatus(val string) bool {
	if _, ok := KeyStatusValue[KeyStatus(val)]; ok {
		return true
	}
	return false
}

//TxType トランザクション種別
type TxType string

// tx_type
const (
	TxTypeUnsigned TxType = "unsigned"
	TxTypeSigned   TxType = "signed"
	TxTypeSent     TxType = "sent"
	TxTypeDone     TxType = "done"
	TxTypeNotified TxType = "notified"
	TxTypeCancel   TxType = "canceled"
)

//TxTypeValue tx_typeの値
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned: 1,
	TxTypeSigned:   2,
	TxTypeSent:     3,
	TxTypeDone:     4,
	TxTypeNotified: 5,
	TxTypeCancel:   6,
}

// ValidateTxType TxTypeのバリデーションを行う
func ValidateTxType(val string) bool {
	if _, ok := TxTypeValue[TxType(val)]; ok {
		return true
	}
	return false
}

//ActionType 入金/出金
type ActionType string

// action
const (
	ActionTypeReceipt ActionType = "receipt"
	ActionTypePayment ActionType = "payment"
)

//ActionTypeValue action_typeの値
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeReceipt: 1,
	ActionTypePayment: 2,
}

// ValidateActionType ActionTypeのバリデーションを行う
func ValidateActionType(val string) bool {
	if _, ok := ActionTypeValue[ActionType(val)]; ok {
		return true
	}
	return false
}
