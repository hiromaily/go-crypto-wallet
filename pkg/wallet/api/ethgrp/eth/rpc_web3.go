package eth

import (
	"github.com/pkg/errors"
)

// ClientVersion returns client version
// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion
// - returns like `Geth/v1.9.14-stable/darwin-amd64/go1.14.2`
//                `Parity-Ethereum//v2.7.2-stable-d961010f63-20200205/x86_64-apple-darwin/rustc1.41.0`
func (e *Ethereum) ClientVersion() (string, error) {
	var resClientVersion string
	err := e.rpcClient.CallContext(e.ctx, &resClientVersion, "web3_clientVersion")
	if err != nil {
		return "", errors.Wrap(err, "fail to call client.CallContext(web3_clientVersion)")
	}
	return resClientVersion, nil
}

// SHA3 returns Keccak-256 (not the standardized SHA3-256) of the given data
// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_sha3
func (e *Ethereum) SHA3(data string) (string, error) {
	var res string
	err := e.rpcClient.CallContext(e.ctx, &res, "web3_sha3", data)
	if err != nil {
		return "", errors.Wrap(err, "fail to call client.CallContext(web3_sha3)")
	}
	return res, nil
}
