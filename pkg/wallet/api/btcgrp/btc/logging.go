package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// LoggingResult is response type of PRC `logging`
type LoggingResult struct {
	Net         bool `json:"net"`
	Tor         bool `json:"tor"`
	Mempool     bool `json:"mempool"`
	HTTP        bool `json:"http"`
	Bench       bool `json:"bench"`
	Zmq         bool `json:"zmq"`
	Walletdb    bool `json:"walletdb"`
	RPC         bool `json:"rpc"`
	Estimatefee bool `json:"estimatefee"`
	Addrman     bool `json:"addrman"`
	Selectcoins bool `json:"selectcoins"`
	Reindex     bool `json:"reindex"`
	Cmpctblock  bool `json:"cmpctblock"`
	Rand        bool `json:"rand"`
	Prune       bool `json:"prune"`
	Proxy       bool `json:"proxy"`
	Mempoolrej  bool `json:"mempoolrej"`
	Libevent    bool `json:"libevent"`
	Coindb      bool `json:"coindb"`
	Qt          bool `json:"qt"`
	Leveldb     bool `json:"leveldb"`
	Validation  bool `json:"validation"`
}

// Logging calls RPC `logging`
func (b *Bitcoin) Logging() (*LoggingResult, error) {
	rawResult, err := b.Client.RawRequest("logging", []json.RawMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.RawRequest(logging)")
	}

	loggingResult := LoggingResult{}
	err = json.Unmarshal(rawResult, &loggingResult)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call json.Unmarshal(rawResult)")
	}

	return &loggingResult, nil
}
