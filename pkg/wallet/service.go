package wallet

import (
	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/gcp"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// Wallet 基底オブジェクト
type Wallet struct {
	//BTC  *btc.Bitcoin
	BTC  api.Bitcoiner
	DB   *model.DB
	GCS  map[enum.ActionType]*gcp.Storage
	Env  enum.EnvironmentType
	Type enum.WalletType
	Seed string
}

//InitialSettings 実行前に必要なすべての設定をこちらで行う
//TODO:hotwalletとColdwalletで設定が異なるので要調整
func InitialSettings(confPath string) (*Wallet, error) {
	// Config
	conf, err := toml.New(confPath)
	if err != nil {
		return nil, errors.Errorf("toml.New() error: %s", err)
	}
	grok.Value(conf)

	// CoinType
	if !enum.ValidateBitcoinType(conf.CoinType) {
		return nil, errors.New("CoinType is invalid in toml file")
	}

	// Log
	logger.Initialize(enum.EnvironmentType(conf.Environment))

	// MySQL
	rds, err := rdb.Connection(&conf.MySQL)
	if err != nil {
		return nil, errors.Errorf("rds.Connection() error: %s", err)
	}

	// TxFile
	if conf.TxFile.BasePath != "" {
		txfile.SetFilePath(conf.TxFile.BasePath)
	}

	// PubkeyCSV
	if conf.PubkeyFile.BasePath != "" {
		key.SetFilePath(conf.PubkeyFile.BasePath)
	}

	// GCS (only watch only wallete)
	gcs := make(map[enum.ActionType]*gcp.Storage)
	if conf.GCS.ReceiptBucketName != "" {
		gcs[enum.ActionTypeReceipt] = gcp.NewStorage(conf.GCS.ReceiptBucketName, conf.GCS.StorageKeyPath)
	}
	if conf.GCS.PaymentBucketName != "" {
		gcs[enum.ActionTypePayment] = gcp.NewStorage(conf.GCS.PaymentBucketName, conf.GCS.StorageKeyPath)
	}
	if conf.GCS.TransferBucketName != "" {
		gcs[enum.ActionTypeTransfer] = gcp.NewStorage(conf.GCS.TransferBucketName, conf.GCS.StorageKeyPath)
	}

	// Connection to Bitcoin core
	bit, err := api.Connection(&conf.Bitcoin, enum.CoinType(conf.CoinType))
	if err != nil {
		return nil, errors.Errorf("btc.Connection error: %s", err)
	}

	//seed (only dev mode)
	var seed string
	if conf.Key.Seed != "" && enum.EnvironmentType(conf.Environment) == enum.EnvDev {
		seed = conf.Key.Seed
	}

	//Wallet Object
	wallet := Wallet{BTC: bit, DB: model.NewDB(rds), GCS: gcs, Env: enum.EnvironmentType(conf.Environment), Seed: seed}
	return &wallet, nil
}

// Done 終了時に必要な処理
func (w *Wallet) Done() {
	w.DB.RDB.Close()
	w.BTC.Close()
}
