package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

// Config ルート
type Config struct {
	CoinType   string         `toml:"coin_type"`
	Bitcoin    BitcoinConf    `toml:"bitcoin"`
	MySQL      MySQLConf      `toml:"mysql"`
	TxFile     TxFileConf     `toml:"tx_file"`
	PubkeyFile PubKeyFileConf `toml:"pubkey_file"`
	Key        KeyConf        `toml:"key"`
}

// BitcoinConf Bitcoin情報
type BitcoinConf struct {
	Host       string `toml:"host"`
	User       string `toml:"user"`
	Pass       string `toml:"pass"`
	PostMode   bool   `toml:"http_post_mode"`
	DisableTLS bool   `toml:"disable_tls"`
	IsMain     bool   `toml:"is_main"`

	Block BitcoinBlockConf `toml:"block"`
	Fee   BitcoinFeeConf   `toml:"fee"`
}

// BitcoinBlockConf Bitcoinブロック情報
type BitcoinBlockConf struct {
	ConfirmationNum int `toml:"confirmation_num"`
}

// BitcoinFeeConf fee調整Range
type BitcoinFeeConf struct {
	AdjustmentMin float64 `toml:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max"`
}

// KeyConf keyのデフォルト情報(devモード時にしか利用しない)
type KeyConf struct {
	Seed string `toml:"seed"`
}

// MySQLConf MySQL情報
type MySQLConf struct {
	Host string `toml:"host"`
	DB   string `toml:"dbname"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

// TxFileConf 保存されるtransactionファイル情報
// import/export共にこのパスが使われる
type TxFileConf struct {
	BasePath string `toml:"base_path"`
}

// PubKeyFileConf 保存されるtransactionファイル情報
// import/export共にこのパスが使われる
type PubKeyFileConf struct {
	BasePath string `toml:"base_path"`
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

	//validate
	if err = conf.validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

// load configfile
func loadConfig(path string) (*Config, error) {
	//読み込み
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf(
			"toml file can't not be read. [path]:%s: [error]:%s", path, err)
	}

	//解析
	var config Config
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.New("toml file can not be parsed")
	}

	return &config, nil
}

func (c *Config) validate() error {
	// CoinType
	if !enum.ValidateBitcoinType(c.CoinType) {
		return errors.New("CoinType is invalid in toml file")
	}
	return nil
}