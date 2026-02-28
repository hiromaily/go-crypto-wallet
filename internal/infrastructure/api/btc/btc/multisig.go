package btc

import (
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// AddMultisigAddress create multisig address
//   - requiredSigs: required number of signature for transaction
//   - addresses:    list of addresses(e.g. client, auth1, auth2, auth3
//   - [N:M] e.g. 2:5 => requiredSigs=2, addresses[addr1, addr2, addr3, addr4, addr5]
//
// Note: This is a legacy wallet method. For descriptor wallets (Bitcoin Core v23.0+),
// consider using descriptor-based multisig with importdescriptors for better functionality.
func (b *Bitcoin) AddMultisigAddress(
	requiredSigs int,
	addresses []string,
	accountName string,
	addressType domainBTC.AddressType,
) (*dtobtc.MultisigAddress, error) {
	if requiredSigs > len(addresses) {
		return nil, fmt.Errorf(
			"number of given address doesn't meet number of requiredSigs: requiredSigs:%d, len(addresses):%d",
			requiredSigs, len(addresses))
	}

	// addressType: BTC passes the type string; BCH passes "" to omit the parameter
	var addrTypeStr string
	switch b.coinTypeCode {
	case domainCoin.BTC:
		addrTypeStr = FromAddressType(addressType).String()
	case domainCoin.BCH:
		addrTypeStr = ""
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.HYT:
		return nil, fmt.Errorf("not implemented for %s in AddMultisigAddress()", b.coinTypeCode.String())
	default:
		return nil, fmt.Errorf("not implemented for %s in AddMultisigAddress()", b.coinTypeCode.String())
	}

	result, err := btcrpc.AddMultisigAddress(b.Client, requiredSigs, addresses, accountName, addrTypeStr)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.AddMultisigAddress(): %w", err)
	}

	return ToMultisigAddressFromPkg(result), nil
}
