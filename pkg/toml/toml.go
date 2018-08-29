package toml

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config ルート
type Config struct {
	Environment string        `toml:"environment"`
	Bitcoin     BitcoinConf   `toml:"bitcoin"`
	MySQL       MySQLConf     `toml:"mysql"`
	TxFile      TxFileConf     `toml:"tx_file"`
	PubkeyFile  PubKeyFileConf `toml:"pubkey_file"`
	GCS         GCSConf        `toml:"gcs"`
	//LevelDB LevelDBConf `toml:"leveldb"` //TODO:おそらく不要
}

// BitcoinConf Bitcoin情報
type BitcoinConf struct {
	Host       string `toml:"host"`
	User       string `toml:"user"`
	Pass       string `toml:"pass"`
	PostMode   bool   `toml:"http_post_mode"`
	DisableTLS bool   `toml:"disable_tls"`
	IsMain     bool   `toml:"is_main"`

	Block   BitcoinBlockConf `toml:"block"`
	Stored  BitcoinAddrConf  `toml:"stored"`
	Payment BitcoinAddrConf  `toml:"payment"`
	Fee     BitcoinFeeConf   `toml:"fee"`
}

// BitcoinBlockConf Bitcoinブロック情報
type BitcoinBlockConf struct {
	ConfirmationNum int `toml:"confirmation_num"`
}

// BitcoinAddrConf 内部利用のためのBitcoin公開アドレス, アカウント情報
type BitcoinAddrConf struct {
	Address     string `toml:"address"`
	AccountName string `toml:"account"`
}

// BitcoinFeeConf fee調整Range
type BitcoinFeeConf struct {
	AdjustmentMin float64 `toml:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max"`
}

// MySQLConf MySQL情報
type MySQLConf struct {
	Host string `toml:"host"`
	DB   string `toml:"dbname"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

// LevelDBConf LevelDB情報
//type LevelDBConf struct {
//	Path string `toml:"path"`
//}

// TxFileConf 保存されるtransactionファイル情報
// import/export共にこのパスが使われる
type TxFileConf struct {
	BasePath string `toml:"base_path"`
}

// PubKeyCSVConf 保存されるtransactionファイル情報
// import/export共にこのパスが使われる
type PubKeyFileConf struct {
	BasePath string `toml:"base_path"`
}

// GCSConf Google Cloud Storage
type GCSConf struct {
	StorageKeyPath    string `toml:"storage_key_path"`
	ReceiptBucketName string `toml:"receipt_bucket_name"`
	PaymentBucketName string `toml:"payment_bucket_name"`
}

// load configfile
func loadConfig(path string) (*Config, error) {
	//読み込み
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf(
			"toml file can't not be read. [path]:%s: [error]:%v", path, err)
	}

	//解析
	var config Config
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.New("toml file can't not be parsed")
	}

	return &config, nil
}

// New configオブジェクトを生成する
func New(file string) (*Config, error) {
	if file == "" {
		return nil, errors.New("file should be passed")
	}

	var err error
	conf, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
