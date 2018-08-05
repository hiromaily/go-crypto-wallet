package api

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

// Bitcoin includes Client to call Json-RPC
type Bitcoin struct {
	Client    *rpcclient.Client
	chainConf *chaincfg.Params
}

// Connection is to local bitcoin core RPC server using HTTP POST mode
func Connection(host, user, pass string, postMode, tls, isMain bool) (*Bitcoin, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: postMode, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   tls,      // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	bit := Bitcoin{Client: client}
	if isMain {
		bit.chainConf = &chaincfg.MainNetParams
	} else {
		bit.chainConf = &chaincfg.TestNet3Params
	}

	return &bit, err
}

func (b *Bitcoin) Close() {
	b.Client.Shutdown()
}

func (b *Bitcoin) GetChainConf() *chaincfg.Params {
	return b.chainConf
}
