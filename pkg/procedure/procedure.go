package procedure

import "github.com/bookerzzz/grok"

//WalletType wallet種別
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "watch_only" //生成したアドレスのみを保持し、Bitcoin core NWに接続可能なwallet
	WalletTypeCold1     WalletType = "cold1"      //通常利用のkeyの生成から管理まで行う非ネットワーク環境下で利用するwallet
	WalletTypeCold2     WalletType = "cold2"      //承認用のアカウント及び、Multisigアドレスの生成を行うwallet
)

// Procedure 手順に伴う情報グループ
type Procedure struct {
	WalletType WalletType
	Indication string
	Command    string
}

//Procedure env_typeの値
var procedures = []Procedure{
	{
		WalletTypeCold1,
		"generate seed",
		"make xxxx",
	},
	{
		WalletTypeCold1,
		"generate [client, receipt, payment] key",
		"make xxxx",
	},
	{
		WalletTypeCold1,
		"run `importprivkey` to register generated [client, receipt, payment] private key",
		"make xxxx",
	},
	{
		WalletTypeCold1,
		"export [client, receipt, payment] public address as csv",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"generate seed",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"generate [authorization] key",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"run `importprivkey` to register generated [authorization] private key",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"import [receipt, payment] public address from csv to DB",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"run `addmultisigaddress` with [receipt] address as param1,authorization address as param2",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"run `addmultisigaddress` with [payment] address as param1,authorization address as param2",
		"make xxxx",
	},
	{
		WalletTypeCold2,
		"export [receipt, payment] multisig address, public address and redeemScript as csv",
		"make xxxx",
	},
	{
		WalletTypeCold1,
		"import [receipt, payment] multisig address, public address and redeemScript from csv to DB",
		"make xxxx",
	},
	{
		WalletTypeCold1,
		"export [receipt, payment] multisig address as csv",
		"make xxxx",
	},
	{
		WalletTypeWatchOnly,
		"import [client] public address",
		"make xxxx",
	},
	{
		WalletTypeWatchOnly,
		"import [receipt, payment] multisig address as public address",
		"make xxxx",
	},
}

// Show Procedureを表示する
func Show() {
	grok.Value(procedures)
}
