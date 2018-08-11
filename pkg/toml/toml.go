package toml

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var conf *Config

// Config is of root
type Config struct {
	HokanAddress         string `toml:"hokan_address"`
	ConfirmationBlockNum int64  `toml:"confirmation_blockNum"`
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
		return nil, errors.New("toml file can't not be parsed.")
	}

	return &config, nil
}

// New configオブジェクトを生成する
func New(file string) (*Config, error) {
	if file == "" {
		return nil, errors.New("file should be passed.")
	}

	var err error
	conf, err = loadConfig(file)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
