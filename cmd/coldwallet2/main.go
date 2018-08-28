package main

import (
	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/jessevdk/go-flags"
)

// HDウォレットとしてseed作成、keyを指定した数だけ生成し、出力する
// 対象アカウント: client, receipt, payment
// 1. create seed
// 2. create key
// 3. run `importprivkey`
// 4. export pubkey from DB

// 5. 未署名トランザクションへの署名

//TODO:encryptwalletコマンドによって、walletを暗号化した場合、秘密鍵を使用するタイミング(未署名トランザクションに署名する)
// でパスフレーズの入力が必要になる

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/cold1_config.toml" description:"Path for configuration toml file"`
	//実行される機能
	Mode uint8 `short:"m" long:"mode" description:"Mode i.e.Functionality"`
	//txファイルパス
	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
	//Debugモード
	Debug bool `short:"d" long:"debug" description:"for only development use"`
}

var (
	opts      Options
	chainConf *chaincfg.Params
)

func init() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}
}

func main() {
	// Config
	wallet, err := service.InitialSettings(opts.ConfPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer wallet.Done()

	if opts.Debug {
		//debug用 機能確認
		debugForCheck(wallet)
	} else {
		//switch mode
		switchFunction(wallet)
	}
}

// 実運用上利用するもののみ、こちらに定義する
func switchFunction(wallet *service.Wallet) {
	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Mode {
	case 1:
		//TODO:これはcoldwallet2(承認用)の機能
		//AuthorizationのKeyを作成する
		logger.Info("Run: AuthorizationのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeAuthorization, bSeed, 2)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)

	default:
		logger.Info("該当Mode無し")
	}
}

// 検証用
func debugForCheck(wallet *service.Wallet) {

}
