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

//WalletType Walletタイプ
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "watch_only"
	WalletTypeCold1     WalletType = "cold1"
	WalletTypeCold2     WalletType = "cold2"
)

//WalletTypeValue env_typeの値
var WalletTypeValue = map[WalletType]uint8{
	WalletTypeWatchOnly: 1,
	WalletTypeCold1:     2,
	WalletTypeCold2:     3,
}

//CoinType Bitcoin種別(CayenneWalletで取引するcoinの種別)
type CoinType string

// coin_type
const (
	BTC CoinType = "btc"
	BCH CoinType = "bch"
	//ETH CoinType = "eth"
)

//CoinTypeValue coin_typeの値
var CoinTypeValue = map[CoinType]uint8{
	BTC: 1,
	BCH: 2,
	//ETH: 3,
}

// ValidateBitcoinType BitcoinTypeのバリデーションを行う
func ValidateBitcoinType(val string) bool {
	if _, ok := CoinTypeValue[CoinType(val)]; ok {
		return true
	}
	return false
}

//BTCVersion 実行環境
type BTCVersion int

// environment
const (
	BTCVer16 BTCVersion = 160000
	BTCVer17 BTCVersion = 170000
	BTCVer18 BTCVersion = 180000
)

//NetworkType ネットワーク種別
type NetworkType string

// network type
const (
	NetworkTypeMainNet    NetworkType = "mainnet"
	NetworkTypeTestNet3   NetworkType = "testnet3"
	NetworkTypeRegTestNet NetworkType = "regtest"
)

//AddressType address種別
type AddressType string

// address type
const (
	AddressTypeLegacy     AddressType = "legacy"
	AddressTypeP2shSegwit AddressType = "p2sh-segwit"
	AddressTypeBech32     AddressType = "bech32"
)

//AddressTypeValue address_typeの値
var AddressTypeValue = map[AddressType]uint8{
	AddressTypeLegacy:     0,
	AddressTypeP2shSegwit: 1,
	AddressTypeBech32:     2,
}

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

//AccountTypeValue account_typeの値
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
	TxTypeUnsigned    TxType = "unsigned"
	TxTypeUnsigned2nd TxType = "unsigned2nd"
	TxTypeSigned      TxType = "signed"
	TxTypeSent        TxType = "sent"
	TxTypeDone        TxType = "done"
	TxTypeNotified    TxType = "notified"
	TxTypeCancel      TxType = "canceled"
)

//TxTypeValue tx_typeの値
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned:    1,
	TxTypeUnsigned2nd: 2,
	TxTypeSigned:      3,
	TxTypeSent:        4,
	TxTypeDone:        5,
	TxTypeNotified:    6,
	TxTypeCancel:      7,
}

// ValidateTxType TxTypeのバリデーションを行う
func ValidateTxType(val string) bool {
	if _, ok := TxTypeValue[TxType(val)]; ok {
		return true
	}
	return false
}

// Search SliceのtxTypes内にtが含まれているかチェックする
func (t TxType) Search(txTypes []TxType) bool {
	for _, v := range txTypes {
		if v == t {
			return true
		}
	}
	return false
}

//ActionType 入金/出金
type ActionType string

// action
const (
	ActionTypeReceipt  ActionType = "receipt"
	ActionTypePayment  ActionType = "payment"
	ActionTypeTransfer ActionType = "transfer"
)

//ActionTypeValue action_typeの値
var ActionTypeValue = map[ActionType]uint8{
	ActionTypeReceipt:  1,
	ActionTypePayment:  2,
	ActionTypeTransfer: 3,
}

// ValidateActionType ActionTypeのバリデーションを行う
func ValidateActionType(val string) bool {
	if _, ok := ActionTypeValue[ActionType(val)]; ok {
		return true
	}
	return false
}

//ActionToAccountMap ActionTypeをAccountTypeでマッピングする
//var ActionToAccountMap = map[ActionType]AccountType{
//	ActionTypeReceipt: AccountTypeClient,
//	ActionTypePayment: AccountTypePayment,
//}
