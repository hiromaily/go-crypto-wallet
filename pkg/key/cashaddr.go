package key

//for Bitcoin cash

// BCHPrefix BCH用address prefix
type BCHPrefix string

// bch_address_prefix
const (
	BCHPrefixMain    BCHPrefix = "bitcoincash:"
	BCHPrefixTestNet BCHPrefix = "bchtest:"
	BCHPrefixRegTest BCHPrefix = "bchreg:"
)

//P2PKH:pubKeyHash
//bchutil.NewCashAddressPubKeyHash()
// NewAddressPubKeyHash returns a new AddressPubKeyHash.  pkHash mustbe 20
// bytes.
//func NewCashAddressPubKeyHash(pkHash []byte, net *chaincfg.Params) (*CashAddressPubKeyHash, error) {
//	return newCashAddressPubKeyHash(pkHash, net)
//}
//func newCashAddressPubKeyHash(pkHash []byte, net *chaincfg.Params) (*CashAddressPubKeyHash, error) {
//	// Check for a valid pubkey hash length.
//	if len(pkHash) != ripemd160.Size {
//		return nil, errors.New("pkHash must be 20 bytes")
//	}
//
//	prefix, ok := Prefixes[net.Name]
//	if !ok {
//		return nil, errors.New("unknown network parameters")
//	}
//
//	addr := &CashAddressPubKeyHash{prefix: prefix}
//	copy(addr.hash[:], pkHash)
//	return addr, nil

//btcutil.NewAddressScriptHash

//P2SH Hash (pay to script hash)
//bchutil.NewCashAddressScriptHashFromHash()
// NewAddressScriptHashFromHash returns a new AddressScriptHash.  scriptHash
// must be 20 bytes.
//func NewCashAddressScriptHashFromHash(scriptHash []byte, net *chaincfg.Params) (*CashAddressScriptHash, error) {
//	return newCashAddressScriptHashFromHash(scriptHash, net)
//}
//func newCashAddressScriptHashFromHash(scriptHash []byte, net *chaincfg.Params) (*CashAddressScriptHash, error) {
//	// Check for a valid script hash length.
//	if len(scriptHash) != ripemd160.Size {
//		return nil, errors.New("scriptHash must be 20 bytes")
//	}
//
//	pre, ok := Prefixes[net.Name]
//	if !ok {
//		return nil, errors.New("unknown network parameters")
//	}
//
//	addr := &CashAddressScriptHash{prefix: pre}
//	copy(addr.hash[:], scriptHash)
//	return addr, nil
//}

// Address converts the extended key to a standard bitcoin pay-to-pubkey-hash
// address for the passed network.
//func Address(k *hdkeychain.ExtendedKey, net *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
//	pkHash := btcutil.Hash160(k.pubKeyBytes())
//	return btcutil.NewAddressPubKeyHash(pkHash, net)
//}
