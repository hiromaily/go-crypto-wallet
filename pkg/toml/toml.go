package toml

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config ルート
type Config struct {
	Bitcoin BitcoinConf `toml:"bitcoin"`
	LevelDB LevelDBConf `toml:"leveldb"` //TODO:おそらく不要
	MySQL   MySQLConf   `toml:"mysql"`
	File    FileConf    `toml:"file"`
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
}

// BitcoinBlockConf Bitcoinブロック情報
type BitcoinBlockConf struct {
	ConfirmationNum int64 `toml:"confirmation_num"`
}

// BitcoinAddrConf 内部利用のためのBitcoin公開アドレス, アカウント情報
type BitcoinAddrConf struct {
	Address     string `toml:"address"`
	AccountName string `toml:"account"`
}

// LevelDBConf LevelDB情報
type LevelDBConf struct {
	Path string `toml:"path"`
}

// MySQLConf MySQL情報
type MySQLConf struct {
	Host string `toml:"host"`
	DB   string `toml:"dbname"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

// FileConf 保存されるtransactionファイル情報
// import/export共にこのパスが使われる
type FileConf struct {
	Path string `toml:"path"`
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
