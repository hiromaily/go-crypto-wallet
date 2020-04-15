package coldwallet

// sign wallet

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// AddMultisigAddress add multisig address by address for auth and given account
// - allowed account is only who has multisig addresses and auth addresses
// - if 3:5 proportion is required, at least 4 different auth accounts should be parepared in advance
// - when sending coin from multisig address, 、related priv key is required which is related to addresses in parameters
// - actually address is overridden by multisig addresses in multisig acccount
// - 4th parameter must be`p2sh-segwit` addressType in Bitcoin
func (w *ColdWallet) AddMultisigAddress(accountType account.AccountType, addressType address.AddrType) error {
	//TODO:remove it
	if w.wtype != types.WalletTypeSignature {
		return errors.New("it's available on sign wallet")
	}
	// validate
	if !account.AccountTypeMultisig[accountType] {
		w.logger.Info("only multisig account is allowed")
		return nil
	}
	if !account.NotAllow(accountType.String(), []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
		return errors.Errorf("account: %s/%s is not allowed", account.AccountTypeAuthorization, account.AccountTypeClient)
	}

	// get one wallet_address for Authorization account from account_key_authorization table
	authKeyTable, err := w.storager.GetOneByMaxIDOnAccountKeyTable(account.AccountTypeAuthorization)
	if err != nil {
		return errors.Wrap(err, "fail to call storager.GetOneByMaxIDOnAccountKeyTable(AccountTypeAuthorization)")
	}
	// get full-pub-key for given account from added_pubkey_history_table
	addedPubkeyHistoryTable, err := w.storager.GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType)
	if err != nil {
		return errors.Wrapf(err, "fail to call storager.GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(%s)", accountType.String())
	}

	// call bitcoinAPI `addmultisigaddress`
	//FIXME: for now only 2:2 proportion is available
	// - however N:M should be adjustable
	for _, val := range addedPubkeyHistoryTable {
		resAddr, err := w.btc.AddMultisigAddress(
			2,
			[]string{
				val.FullPublicKey, // receipt, payment, stored ...
				authKeyTable.P2shSegwitAddress,
			},
			fmt.Sprintf("multi_%s", accountType), //TODO:what account name is understandable?
			addressType,
		)
		if err != nil {
			//[Error] -5: no full public key for address mkPmdpo59gpU7ZioGYwwoMTQJjh7MiqUvd
			w.logger.Error(
				"fail to call btc.CreateMultiSig(2,,) ",
				zap.String("full public key", val.FullPublicKey),
				zap.String("p2sh segwit address", authKeyTable.P2shSegwitAddress),
				zap.Error(err))
			continue
		}

		// store generated address into added_pubkey_history_table
		err = w.storager.UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType, resAddr.Address,
			resAddr.RedeemScript, authKeyTable.P2shSegwitAddress, val.FullPublicKey, nil, true)
		if err != nil {
			w.logger.Error(
				"fail to call db.UpdateMultisigAddrOnAddedPubkeyHistoryTable()",
				zap.String("accountType", accountType.String()),
				zap.Error(err))
		}
	}

	return nil
}
