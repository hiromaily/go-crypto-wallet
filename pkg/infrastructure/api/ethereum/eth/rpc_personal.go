package eth

import (
	"context"
	"fmt"
	"strings"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ImportRawKey Imports the given unencrypted private key (hex string) into the key store,
// encrypting it with the passphrase
//   - if address has 0x like 0x5d0a82e19564ae03ad3b834ac30b94c0ccce510e86d783d3e882efcb0e84b2af,
//     error would occur `invalid hex character 'x' in private key`
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_importrawkey
// - private key is stored and read from node server
func (e *Ethereum) ImportRawKey(ctx context.Context, hexKey, passPhrase string) (string, error) {
	if strings.HasPrefix(hexKey, "0x") {
		hexKey = strings.TrimLeft(hexKey, "0x")
	}

	var address string
	err := e.rpcClient.CallContext(ctx, &address, "personal_importRawKey", hexKey, passPhrase)
	if err != nil {
		if strings.Contains(err.Error(), "account already exists") {
			logger.Warn(err.Error())
		} else {
			return "", fmt.Errorf("fail to call client.CallContext(personal_importRawKey): %w", err)
		}
	}

	return address, nil
}

// ListAccounts returns all the Ethereum account addresses of all keys in the key store
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_listaccounts
// - private key is stored and read from node server
// - result is same as Accounts()
func (e *Ethereum) ListAccounts(ctx context.Context) ([]string, error) {
	var accounts []string
	err := e.rpcClient.CallContext(ctx, &accounts, "personal_listAccounts")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(personal_listAccounts): %w", err)
	}

	return accounts, nil
}

// NewAccount generates a new private key and stores it in the key store directory
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_newaccount
// - private key is stored and read from node server
func (e *Ethereum) NewAccount(
	ctx context.Context, passphrase string, accountType domainAccount.AccountType,
) (string, error) {
	var address string
	err := e.rpcClient.CallContext(ctx, &address, "personal_newAccount", passphrase)
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(personal_newAccount): %w", err)
	}
	return address, nil
}

// LockAccount removes the private key with given address from memory
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_lockaccount
func (e *Ethereum) LockAccount(ctx context.Context, hexAddr string) error {
	err := e.rpcClient.CallContext(ctx, nil, "personal_lockAccount", hexAddr)
	if err != nil {
		return fmt.Errorf("fail to call rpc.CallContext(personal_lockAccount): %w", err)
	}
	return err
}

// UnlockAccount decrypts the key with the given address from the key store.
//
//	duration: second
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_unlockaccount
// FIXME: how to fix error `account unlock with HTTP access is forbidden`
// - --allow-insecure-unlock option is required to run geth
func (e *Ethereum) UnlockAccount(ctx context.Context, hexAddr, passphrase string, duration uint64) (bool, error) {
	// eoncode duration
	// encodedDuration := hexutil.EncodeUint64(duration)

	var ret bool
	err := e.rpcClient.CallContext(ctx, &ret, "personal_unlockAccount", hexAddr, passphrase, duration)
	if err != nil {
		return false, fmt.Errorf("fail to call rpc.CallContext(personal_unlockAccount): %w", err)
	}
	return ret, nil
}
