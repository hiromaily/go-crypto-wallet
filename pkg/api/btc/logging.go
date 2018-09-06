package btc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// LoggingResult logging RPC のレスポンス
type LoggingResult struct {
	Net         int64 `json:"net"`
	Tor         int64 `json:"tor"`
	Mempool     int64 `json:"mempool"`
	HTTP        int64 `json:"http"`
	Bench       int64 `json:"bench"`
	Zmq         int64 `json:"zmq"`
	DB          int64 `json:"db"`
	RPC         int64 `json:"rpc"`
	EstimateFee int64 `json:"estimatefee"`
	Addrman     int64 `json:"addrman"`
	SelectCoins int64 `json:"selectcoins"`
	ReIndex     int64 `json:"reindex"`
	CmpctBlock  int64 `json:"cmpctblock"`
	Rand        int64 `json:"rand"`
	Prune       int64 `json:"prune"`
	Proxy       int64 `json:"proxy"`
	MempoolRej  int64 `json:"mempoolrej"`
	LibEvent    int64 `json:"libevent"`
	CoinDB      int64 `json:"coindb"`
	Qt          int64 `json:"qt"`
	LevelDB     int64 `json:"leveldb"`
}

// Logging logging RPC をcallする
func (b *Bitcoin) Logging() (*LoggingResult, error) {
	rawResult, err := b.client.RawRequest("logging", []json.RawMessage{})
	if err != nil {
		return nil, errors.Errorf("json.RawRequest(logging): error: %v", err)
	}

	loggingResult := LoggingResult{}
	err = json.Unmarshal([]byte(rawResult), &loggingResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %v", err)
	}

	return &loggingResult, nil
}
