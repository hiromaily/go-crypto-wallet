package eth

import (
	"github.com/pkg/errors"
)

// ClientVersion returns client version
func (e *Ethereum) ClientVersion() (string, error) {
	var resClientVersion string
	err := e.client.CallContext(e.ctx, &resClientVersion, "web3_clientVersion")
	if err != nil {
		return "", errors.Wrap(err, "fail to call client.CallContext(web3_clientVersion)")
	}
	return resClientVersion, nil
}
