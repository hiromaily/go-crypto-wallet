package coldwallet

// sign wallet

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
)

// AddMultisigAddress add multisig address by address for auth and given account
// - allowed account is only who has multisig addresses and auth addresses
// - if 3:5 proportion is required, at least 4 different auth accounts should be parepared in advance
// - when sending coin from multisig address, 、related priv key is required which is related to addresses in parameters
// - actually address is overridden by multisig addresses in multisig acccount
// - 4th parameter must be`p2sh-segwit` addressType in Bitcoin
func (w *ColdWallet) AddMultisigAddress(accountType account.AccountType, addressType address.AddrType) error {
	// validate
	if !account.AccountTypeMultisig[accountType] {
		w.logger.Info("only multisig account is allowed")
		return nil
	}
	if !account.NotAllow(accountType.String(), []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
		return errors.Errorf("account: %s/%s is not allowed", account.AccountTypeAuthorization, account.AccountTypeClient)
	}

	// get one wallet_address for Authorization account from account_key_authorization table
	authKeyTable, err := w.repo.AccountKey().GetOneMaxID(account.AccountTypeAuthorization)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.GetOneByMaxIDOnAccountKeyTable(AccountTypeAuthorization)")
	}
	// get full-pub-key for given account from added_pubkey_history_table
	multisigHistoryTable, err := w.repo.MultisigHistory().GetAllNoMultisig(accountType)
	if err != nil {
		return errors.Wrapf(err, "fail to call repo.GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(%s)", accountType.String())
	}

	// call bitcoinAPI `addmultisigaddress`
	//FIXME: for now only 2:2 proportion is available
	// - however N:M should be adjustable
	for _, val := range multisigHistoryTable {
		resAddr, err := w.btc.AddMultisigAddress(
			2,
			[]string{
				val.FullPublicKey, // receipt, payment, stored ...
				authKeyTable.P2SHSegwitAddress,
			},
			fmt.Sprintf("multi_%s", accountType), //TODO:what account name is understandable?
			addressType,
		)
		if err != nil {
			//[Error] -5: no full public key for address mkPmdpo59gpU7ZioGYwwoMTQJjh7MiqUvd
			w.logger.Error(
				"fail to call btc.CreateMultiSig(2,,) ",
				zap.String("full public key", val.FullPublicKey),
				zap.String("p2sh segwit address", authKeyTable.P2SHSegwitAddress),
				zap.Error(err))
			continue
		}

		// store generated address into added_pubkey_history_table
		_, err = w.repo.MultisigHistory().UpdateMultisigAddr(
			accountType,
			resAddr.Address,
			resAddr.RedeemScript,
			authKeyTable.P2SHSegwitAddress,
			val.FullPublicKey)
		if err != nil {
			w.logger.Error(
				"fail to call repo.MultisigHistory().UpdateMultisigAddr()",
				zap.String("accountType", accountType.String()),
				zap.Error(err))
		}
	}

	return nil
}
