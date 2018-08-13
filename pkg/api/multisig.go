package api

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// AddMultisigAddressResult addmultisigaddressをcallしたresponseの型
type AddMultisigAddressResult struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// CreateMultiSig マルチシグを Rawトランザクション用に作成する
// requiredSigs: 取引成立に必要なサイン数
// addresses: 自分のアドレス+承認者のアドレスxN をいれていく
// TODO:各ユーザーに割り振るアドレスは、Multisigのアドレスになる？？
// TODO:ここで作ったアドレスに送金してみてtest
//func (b *Bitcoin) CreateMultiSig(requiredSigs int, addresses []string, accountName string) (btcutil.Address, error) {
func (b *Bitcoin) CreateMultiSig(requiredSigs int, addresses []string, accountName string) (*AddMultisigAddressResult, error) {

	if requiredSigs > len(addresses) {
		return nil, errors.New("number of given address should be at least same to requiredSigs or more")
	}

	//TODO:NativeのJson RPCを使うのであれば不要かと
	addrs := make([]btcutil.Address, len(addresses))
	for idx, ad := range addresses {
		add, err := b.DecodeAddress(ad)
		if err != nil {
			return nil, err
		}
		//
		//addrs = append(addrs, add)
		addrs[idx] = add
	}
	// deprecateされているので、こちらは使用しない。
	//res, err := b.client.CreateMultisig(requiredSigs, addrs)
	//error: -5: Invalid public key
	// Note that from v0.16, createmultisig no longer accepts addresses.
	// Clients must transition to using addmultisigaddress to create multisig addresses with addresses known to the wallet before upgrading to v0.17.
	// To use the deprecated functionality, start bitcoind with -deprecatedrpc=createmultisig

	// こちらのfuncはjsonのI/Fが実際のBitcoin coreのAPIから乖離してしまっている。。。
	// btcsuite/btcd/rpcclient/wallet.goの (r FutureAddMultisigAddressResult) Receive() のI/Fが古くてjsonとしてParseできん。。。
	//resAddr, err := b.client.AddMultisigAddress(requiredSigs, addrs, accountName)
	//if err != nil {
	//	//error: json: cannot unmarshal object into Go value of type string
	//	return nil, errors.Errorf("client.CreateMultisig(): error: %v", err)
	//}

	//requiredSigs
	bRequiredSigs, err := json.Marshal(requiredSigs)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(requiredSigs): error: %v", err)
	}

	//addresses
	bAddresses, err := json.Marshal(struct {
		Addresses []string `json:"keys"`
	}{
		Addresses: addresses,
	})
	if err != nil {
		return nil, errors.Errorf("json.Marchal(addresses): error: %v", err)
	}

	jsonRawMsg := []json.RawMessage{}

	//accountName
	if accountName != "" {
		bAccountName, err := json.Marshal(accountName)
		if err != nil {
			return nil, errors.Errorf("json.Marchal(accountName): error: %v", err)
		}
		jsonRawMsg = []json.RawMessage{bRequiredSigs, bAddresses, bAccountName}
	} else {
		jsonRawMsg = []json.RawMessage{bRequiredSigs, bAddresses}
	}

	//call addmultisigaddress
	rawResult, err := b.client.RawRequest("addmultisigaddress", jsonRawMsg)
	if err != nil {
		return nil, errors.Errorf("client.RawRequest(addmultisigaddress): error: %v", err)
	}

	multisigAddrResult := AddMultisigAddressResult{}
	err = json.Unmarshal([]byte(rawResult), &multisigAddrResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(multisigAddrResult): error: %v", err)
	}

	return &multisigAddrResult, nil
}
