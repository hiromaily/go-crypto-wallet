package api

import (
	"github.com/btcsuite/btcd/rpcclient"
)

// Bitcoin includes Client to call Json-RPC
type Bitcoin struct {
	Client *rpcclient.Client
}

// Connection is to local bitcoin core RPC server using HTTP POST mode
func Connection(host, user, pass string, postMode, tls bool) (*Bitcoin, error) {
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

	return &Bitcoin{Client: client}, err
}

func (b *Bitcoin) Close(){
	b.Client.Shutdown()
}
