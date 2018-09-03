package procedure

import "github.com/bookerzzz/grok"

//WalletType wallet種別
type WalletType string

// wallet_type
const (
	WalletTypeWatchOnly WalletType = "wallet"      //生成したアドレスのみを保持し、Bitcoin core NWに接続可能なwallet
	WalletTypeCold1     WalletType = "coldwallet1" //通常利用のkeyの生成から管理まで行う非ネットワーク環境下で利用するwallet
	WalletTypeCold2     WalletType = "coldwallet2" //承認用のアカウント及び、Multisigアドレスの生成を行うwallet
)

// Procedure 手順に伴う情報グループ
type Procedure struct {
	WalletType WalletType
	Indication string
	Command    string
}

//proceduresForPreparation セットアップに必要な手順
var proceduresForPreparation = []Procedure{
	{
		WalletTypeCold1,
		"for all",
		"coldwallet1 -w 1 -d",
	},
	{
		WalletTypeCold1,
		"generate seed",
		"coldwallet1 -w 1 -k -m 1",
	},
	{
		WalletTypeCold1,
		"generate [client, receipt, payment] key",
		"coldwallet1 -w 1 -k -m 10,11,12",
	},
	{
		WalletTypeCold1,
		"run `importprivkey` to register generated [client, receipt, payment] private key",
		"coldwallet1 -w 1 -k -m 20,21,22",
	},
	{
		WalletTypeCold1,
		"export [client, receipt, payment] public address as csv",
		"coldwallet1 -w 1 -k -m 30,31,32",
	},
	{
		WalletTypeCold2,
		"generate seed",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 1",
	},
	{
		WalletTypeCold2,
		"generate [authorization] key",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 13",
	},
	{
		WalletTypeCold2,
		"run `importprivkey` to register generated [authorization] private key",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 23",
	},
	{
		WalletTypeCold2,
		"import [receipt] public address from csv to DB",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 33 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold2,
		"import [payment] public address from csv to DB",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 34 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold2,
		"run `addmultisigaddress` with [receipt] address as param1,authorization address as param2",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 50",
	},
	{
		WalletTypeCold2,
		"run `addmultisigaddress` with [payment] address as param1,authorization address as param2",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 51",
	},
	{
		WalletTypeCold2,
		"export [receipt, payment] multisig address, public address and redeemScript as csv",
		"coldwallet1 -c ./data/toml/cold2_config.toml -w 2 -k -m 60,61",
	},
	{
		WalletTypeCold1,
		"import [receipt] multisig address and redeemScript from csv to DB",
		"coldwallet1 -w 1 -k -m 40 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold1,
		"import [payment] multisig address and redeemScript from csv to DB",
		"coldwallet1 -w 1 -k -m 41 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold1,
		"export [receipt, payment] multisig address as csv",
		"coldwallet1 -w 1 -k -m 50, 51",
	},
	{
		WalletTypeWatchOnly,
		"import [client] address for watch only wallet",
		"wallet -k -m 1 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"import [receipt] address for watch only wallet",
		"wallet -k -m 2 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"import [payment] address for watch only wallet",
		"wallet -k -m 3 -i ./data/pubkey/xxx.csv",
	},
}

// Show Procedureを表示する
func Show() {
	grok.Value(proceduresForPreparation)
}
