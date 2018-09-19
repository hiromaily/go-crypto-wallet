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

//proceduresForWallet コマンド一覧
var proceduresForWallet = []Procedure{
	{
		WalletTypeWatchOnly,
		"[キー管理]coldwalletで生成した[client]アドレスをwalletにimportする",
		"wallet -k -m 1 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"[キー管理]coldwalletで生成した[receipt]アドレスをwalletにimportする",
		"wallet -k -m 2 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"[キー管理]coldwalletで生成した[payment]アドレスをwalletにimportする",
		"wallet -k -m 3 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"[入金管理]入金処理検知 + 未署名トランザクション作成",
		"wallet -r -m 1",
	},
	{
		WalletTypeWatchOnly,
		"[入金管理]入金処理検知 (確認のみ)",
		"wallet -r -m 2",
	},
	{
		WalletTypeWatchOnly,
		"[入金管理][Debug用]入金から送金までの一連の流れを確認",
		"wallet -r -m 10",
	},
	{
		WalletTypeWatchOnly,
		"[出金管理]出金のための未署名トランザクション作成",
		"wallet -p -m 1",
	},
	{
		WalletTypeWatchOnly,
		"[出金管理][Debug用]出金から送金までの一連の流れを確認",
		"wallet -p -m 10",
	},
	{
		WalletTypeWatchOnly,
		"[署名送信管理]ファイルから署名済みtxを送信する",
		"wallet -s -m 1 -i ./data/tx/xxxxx/xxxxx",
	},
	{
		WalletTypeWatchOnly,
		"[監視管理]送信済ステータスのトランザクションを監視する",
		"wallet -n -m 1",
	},
	{
		WalletTypeWatchOnly,
		"[bitcoin cli] listlockunspent lockされたトランザクションの解除",
		"wallet -b -m 1",
	},
	{
		WalletTypeWatchOnly,
		"[bitcoin cli] estimatesmartfee 手数料算出",
		"wallet -b -m 2",
	},
	{
		WalletTypeWatchOnly,
		"[bitcoin cli] logging ロギング",
		"wallet -b -m 3",
	},
	{
		WalletTypeWatchOnly,
		"[bitcoin cli] getnetworkinfo 情報取得",
		"wallet -b -m 4",
	},
	{
		WalletTypeWatchOnly,
		"[bitcoin cli] validateaddress AddressのValidationチェック",
		"wallet -b -m 5",
	},
	{
		WalletTypeWatchOnly,
		"[Debug]payment_requestテーブルを作成する",
		"wallet -d -m 1",
	},
	{
		WalletTypeWatchOnly,
		"[Debug]payment_requestテーブルの情報を初期化する",
		"wallet -d -m 2",
	},
}

//proceduresForColdWallet1 セットアップに必要な手順
var proceduresForColdWallet1 = []Procedure{
	{
		WalletTypeCold1,
		"for all",
		"coldwallet1 -d",
	},
	{
		WalletTypeCold1,
		"[キー管理] generate seed",
		"coldwallet1 -k -m 1",
	},
	{
		WalletTypeCold1,
		"[キー管理] generate [client, receipt, payment] key",
		"coldwallet1 -k -m 10,11,12",
	},
	{
		WalletTypeCold1,
		"[キー管理] run `importprivkey` to register generated [client, receipt, payment] private key",
		"coldwallet1 -k -m 20,21,22",
	},
	{
		WalletTypeCold1,
		"[キー管理] export [client, receipt, payment] public address as csv => next coldwallet2",
		"coldwallet1 -k -m 30,31,32",
	},
	{
		WalletTypeCold1,
		"[キー管理] import [receipt] multisig address and redeemScript from csv to DB",
		"coldwallet1 -k -m 40 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold1,
		"[キー管理] import [payment] multisig address and redeemScript from csv to DB",
		"coldwallet1 -k -m 41 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold1,
		"[キー管理] export [receipt, payment] multisig address as csv",
		"coldwallet1 -k -m 50, 51",
	},
	{
		WalletTypeWatchOnly,
		"[キー管理] import [client] address for watch only wallet",
		"wallet -k -m 1 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"[キー管理] import [receipt] address for watch only wallet",
		"wallet -k -m 2 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeWatchOnly,
		"[キー管理] import [payment] address for watch only wallet",
		"wallet -k -m 3 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold1,
		"[署名管理] sign [receipt] tx",
		"wallet -s -m 1 -i ./data/tx/receipt/xxx",
	},
	{
		WalletTypeCold1,
		"[署名管理] sign [payment] tx",
		"wallet -s -m 1 -i ./data/tx/payment/xxx",
	},
}

//proceduresForColdWallet2 セットアップに必要な手順
var proceduresForColdWallet2 = []Procedure{
	{
		WalletTypeCold2,
		"[キー管理] generate seed",
		"coldwallet2 -k -m 1",
	},
	{
		WalletTypeCold2,
		"[キー管理] generate [authorization] key",
		"coldwallet2 -k -m 13",
	},
	{
		WalletTypeCold2,
		"[キー管理] run `importprivkey` to register generated [authorization] private key",
		"coldwallet2 -k -m 23",
	},
	{
		WalletTypeCold2,
		"[キー管理] import [receipt] public address from csv to DB",
		"coldwallet2 -k -m 33 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold2,
		"[キー管理] import [payment] public address from csv to DB",
		"coldwallet2 -k -m 34 -i ./data/pubkey/xxx.csv",
	},
	{
		WalletTypeCold2,
		"[キー管理] run `addmultisigaddress` with [receipt] address as param1,authorization address as param2",
		"coldwallet2 -k -m 50",
	},
	{
		WalletTypeCold2,
		"[キー管理] run `addmultisigaddress` with [payment] address as param1,authorization address as param2",
		"coldwallet2 -k -m 51",
	},
	{
		WalletTypeCold2,
		"[キー管理] export [receipt, payment] multisig address, public address and redeemScript as csv, next coldwallet1",
		"coldwallet2 -k -m 60,61",
	},
	{
		WalletTypeCold2,
		"[署名管理] sign [payment] tx",
		"wallet -s -m 1 -i ./data/tx/payment/xxx",
	},
}

// ShowWallet Procedureを表示する
func ShowWallet() {
	grok.Value(proceduresForWallet)
}

// ShowColdWallet1 coldwallet1のProcedureを表示する
func ShowColdWallet1() {
	grok.Value(proceduresForColdWallet1)
}

// ShowColdWallet2 coldwallet2のProcedureを表示する
func ShowColdWallet2() {
	grok.Value(proceduresForColdWallet2)
}
