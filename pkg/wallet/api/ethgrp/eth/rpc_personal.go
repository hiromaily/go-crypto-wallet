package eth

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// ImportRawKey Imports the given unencrypted private key (hex string) into the key store, encrypting it with the passphrase
// - if address has 0x like 0x5d0a82e19564ae03ad3b834ac30b94c0ccce510e86d783d3e882efcb0e84b2af,
//    error would occur `invalid hex character 'x' in private key`
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_importrawkey
// - private key is stored and read from node server
func (e *Ethereum) ImportRawKey(hexKey, passPhrase string) (string, error) {
	if strings.HasPrefix(hexKey, "0x") {
		hexKey = strings.TrimLeft(hexKey, "0x")
	}

	var address string
	err := e.rpcClient.CallContext(e.ctx, &address, "personal_importRawKey", hexKey, passPhrase)
	if err != nil {
		if strings.Contains(err.Error(), "account already exists") {
			e.logger.Warn(err.Error())
		} else {
			return "", errors.Wrap(err, "fail to call client.CallContext(personal_importRawKey)")
		}
	}

	return address, nil
}

// ListAccounts returns all the Ethereum account addresses of all keys in the key store
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_listaccounts
// - private key is stored and read from node server
// - result is same as Accounts()
func (e *Ethereum) ListAccounts() ([]string, error) {
	var accounts []string
	err := e.rpcClient.CallContext(e.ctx, &accounts, "personal_listAccounts")
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(personal_listAccounts)")
	}

	return accounts, nil

}

// NewAccount generates a new private key and stores it in the key store directory
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_newaccount
// - private key is stored and read from node server
func (e *Ethereum) NewAccount(passphrase string, accountType account.AccountType) (string, error) {

	var address string
	err := e.rpcClient.CallContext(e.ctx, &address, "personal_newAccount", passphrase)
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(personal_newAccount)")
	}
	if e.isParity {
		//TODO:parity client generates file with UUID, so add address to filename if parity client
		err = e.RenameParityKeyFile(address, accountType)
		if err != nil {
			return "", errors.Wrap(err, "fail to call RenameParityKeyFile(address)")
		}
	}
	return address, nil
}

// LockAccount removes the private key with given address from memory
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_lockaccount
func (e *Ethereum) LockAccount(hexAddr string) error {

	err := e.rpcClient.CallContext(e.ctx, nil, "personal_lockAccount", hexAddr)
	if err != nil {
		return errors.Wrap(err, "fail to call rpc.CallContext(personal_lockAccount)")
	}
	return err
}

// UnlockAccount decrypts the key with the given address from the key store.
//  duration: second
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_unlockaccount
// FIXME: how to fix error `account unlock with HTTP access is forbidden`
// - --allow-insecure-unlock option is required to run geth
func (e *Ethereum) UnlockAccount(hexAddr, passphrase string, duration uint64) (bool, error) {

	//eoncode duration
	//encodedDuration := hexutil.EncodeUint64(duration)

	var ret bool
	err := e.rpcClient.CallContext(e.ctx, &ret, "personal_unlockAccount", hexAddr, passphrase, duration)
	if err != nil {
		return false, errors.Wrap(err, "fail to call rpc.CallContext(personal_unlockAccount)")
	}
	return ret, nil
}
