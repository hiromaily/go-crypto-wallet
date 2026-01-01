package cryptocurrency

import (
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewBitcoinRPCClient try to connect bitcoin core RPCserver to create client instance
// using HTTP POST mode
func NewBitcoinRPCClient(conf *config.Bitcoin) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Pass,
		HTTPPostMode: conf.PostMode,   // Bitcoin core only supports HTTP POST mode
		DisableTLS:   conf.DisableTLS, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("rpcclient.New() error: %s", err)
	}
	return client, err
}
