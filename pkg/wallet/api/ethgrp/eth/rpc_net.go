package eth

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// NetVersion returns the current network id
// "1": Ethereum Mainnet
// "2": Morden Testnet (deprecated)
// "3": Ropsten Testnet
// "4": Rinkeby Testnet
// "5": Goerli Testnet
// "42": Kovan Testnet
func (e *Ethereum) NetVersion() (uint16, error) {
	var resNetVersion string
	err := e.rpcClient.CallContext(e.ctx, &resNetVersion, "net_version")
	if err != nil {
		return 0, errors.Wrap(err, "fail to call client.CallContext(net_version)")
	}
	u, err := strconv.ParseUint(resNetVersion, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call strconv.ParseUint(%s)", resNetVersion)
	}

	return uint16(u), nil
}

// NetListening returns true if client is actively listening for network connections
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_listening
func (e *Ethereum) NetListening() (bool, error) {
	var isConnected bool
	err := e.rpcClient.CallContext(e.ctx, &isConnected, "net_listening")
	if err != nil {
		return false, errors.Wrap(err, "fail to call rpc.CallContext(net_listening)")
	}

	return isConnected, nil
}

// NetPeerCount returns number of peers currently connected to the client
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount
func (e *Ethereum) NetPeerCount() (*big.Int, error) {
	var resPeerNumber string
	err := e.rpcClient.CallContext(e.ctx, &resPeerNumber, "net_peerCount")
	if err != nil {
		return nil, errors.Wrap(err, "fail to call client.CallContext(net_peerCount)")
	}
	return hexutil.DecodeBig(resPeerNumber)
}
