package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

// Config root config
type Config struct {
	CoinType   string         `toml:"coin_type"`
	Bitcoin    BitcoinConf    `toml:"bitcoin"`
	Logger     LoggerConf     `toml:"bitcoin"`
	MySQL      MySQLConf      `toml:"mysql"`
	TxFile     TxFileConf     `toml:"tx_file"`
	PubkeyFile PubKeyFileConf `toml:"pubkey_file"`
}

// BitcoinConf Bitcoin information
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

// BitcoinBlockConf block information of Bitcoin
type BitcoinBlockConf struct {
	ConfirmationNum int `toml:"confirmation_num"`
}

// BitcoinFeeConf range of adjustment calculated fee when sending coin
type BitcoinFeeConf struct {
	AdjustmentMin float64 `toml:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max"`
}

// LoggerConf logger info
type LoggerConf struct {
	Service string `toml:"host"`
	Level   string `toml:"level"`
}

// MySQLConf MySQL info
type MySQLConf struct {
	Host string `toml:"host"`
	DB   string `toml:"dbname"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

// TxFileConf saved transaction file path which is used when import/export file
type TxFileConf struct {
	BasePath string `toml:"base_path"`
}

// PubKeyFileConf saved pubKey file path which is used when import/export file
type PubKeyFileConf struct {
	BasePath string `toml:"base_path"`
}

// New create config
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

// load config file
func loadConfig(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf(
			"toml file can't not be read. [path]:%s: [error]:%s", path, err)
	}

	var config Config
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.New("toml file can not be parsed")
	}

	return &config, nil
}

// validate config
func (c *Config) validate() error {
	// CoinType
	if !enum.ValidateBitcoinType(c.CoinType) {
		return errors.New("CoinType is invalid in toml file")
	}
	return nil
}
