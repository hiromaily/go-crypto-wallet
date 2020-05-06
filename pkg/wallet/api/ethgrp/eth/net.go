package eth

import (
	"strconv"

	"github.com/pkg/errors"
)

// NetVersion returns network id
// "1": Ethereum Mainnet
// "2": Morden Testnet (deprecated)
// "3": Ropsten Testnet
// "4": Rinkeby Testnet
// "42": Kovan Testnet
func (e *Ethereum) NetVersion() (uint16, error) {
	var resNetVersion string
	err := e.client.CallContext(e.ctx, &resNetVersion, "net_version")
	if err != nil {
		return 0, errors.Wrap(err, "fail to call client.CallContext(net_version)")
	}
	u, err := strconv.ParseUint(resNetVersion, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call strconv.ParseUint(%s)", resNetVersion)
	}

	return uint16(u), nil
}
