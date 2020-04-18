package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

//TODO:
// - use https://github.com/spf13/viper
// - use https://github.com/go-playground/validator [done]

// Config root config
type Config struct {
	CoinTypeCode coin.CoinTypeCode `toml:"coin_type" validate:"oneof=btc bch"`
	Bitcoin      Bitcoin           `toml:"bitcoin"`
	Logger       Logger            `toml:"logger"`
	Tracer       Tracer            `toml:"tracer"`
	MySQL        MySQL             `toml:"mysql"`
	TxFile       TxFile            `toml:"tx_file"`
	PubkeyFile   PubKeyFile        `toml:"pubkey_file"`
}

// Bitcoin Bitcoin information
type Bitcoin struct {
	Host        string           `toml:"host" validate:"required"`
	User        string           `toml:"user" validate:"required"`
	Pass        string           `toml:"pass" validate:"required"`
	PostMode    bool             `toml:"http_post_mode"`
	DisableTLS  bool             `toml:"disable_tls"`
	NetworkType coin.NetworkType `toml:"network_type" validate:"oneof=mainnet testnet3 regtest"`

	Block BitcoinBlock `toml:"block"`
	Fee   BitcoinFee   `toml:"fee"`
}

// BitcoinBlock block information of Bitcoin
// FIXME: keygen/signature wallet doesn't have this value
//  so validation can not be used
type BitcoinBlock struct {
	ConfirmationNum int `toml:"confirmation_num"`
}

// BitcoinFee range of adjustment calculated fee when sending coin
type BitcoinFee struct {
	AdjustmentMin float64 `toml:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max"`
}

// Logger logger info
type Logger struct {
	Service string `toml:"service" validate:"required"`
	Level   string `toml:"level" validate:"required"`
	Env     string `toml:"env" validate:"oneof=dev prod custom"`
}

// Tracer is open tracing
type Tracer struct {
	Type    string       `toml:"type" validate:"oneof=none jaeger datadog"`
	Jaeger  TracerDetail `toml:"jaeger"`
	Datadog TracerDetail `toml:"datadog"`
}

type TracerDetail struct {
	ServiceName         string  `toml:"service_name"`
	CollectorEndpoint   string  `toml:"collector_endpoint"`
	SamplingProbability float64 `toml:"sampling_probability"`
	IsDebug             bool    `toml:"is_debug"`
}

// MySQL MySQL info
type MySQL struct {
	Host string `toml:"host" validate:"required"`
	DB   string `toml:"dbname" validate:"required"`
	User string `toml:"user" validate:"required"`
	Pass string `toml:"pass" validate:"required"`
}

// TxFile saved transaction file path which is used when import/export file
type TxFile struct {
	BasePath string `toml:"base_path" validate:"required"`
}

// PubKeyFile saved pubKey file path which is used when import/export file
type PubKeyFile struct {
	BasePath string `toml:"base_path" validate:"required"`
}

// New create config
func New(file string, wtype types.WalletType) (*Config, error) {
	if file == "" {
		return nil, errors.New("file should be passed")
	}

	var err error
	conf, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	//debug
	//grok.Value(conf)

	//validate
	if err = conf.validate(wtype); err != nil {
		return nil, err
	}

	return conf, nil
}

// loadConfig load config file
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
func (c *Config) validate(wtype types.WalletType) error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}

	switch wtype {
	case types.WalletTypeWatchOnly:
		if c.Bitcoin.Block.ConfirmationNum == 0 {
			return errors.New("Block ConfirmationNum is required in toml file")
		}
	default:
	}
	//if !ctype.ValidateBitcoinType(c.CoinType) {
	//	return errors.New("CoinType is invalid in toml file")
	//}
	return nil
}
