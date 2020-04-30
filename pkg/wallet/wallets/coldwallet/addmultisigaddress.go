package coldwallet

// sign wallet

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
)

// AddMultisigAddress add multisig address by auth account address and given account address
// https://bitcoincore.org/en/doc/0.19.0/rpc/wallet/addmultisigaddress/
// - if 3:5 proportion is required, at least 4 different auth accounts should be prepared in advance
// - when sending coin from multisig address, 、related priv key is required which is related to addresses in parameters
// - 4th parameter must be`p2sh-segwit` addressType in Bitcoin
//  What is the difference between createmultisig and addmultisigaddress?
// - https://bitcointalk.org/index.php?topic=3402541.0
func (w *ColdWallet) AddMultisigAddress(accountType account.AccountType, authType account.AuthType, addressType address.AddrType) error {
	//for sign wallet
	w.logger.Debug("addmultisigaddress",
		zap.String("account_type", accountType.String()),
		zap.String("auth_type", authType.String()),
	)

	// validate
	if !account.AccountTypeMultisig[accountType] {
		w.logger.Info("only multisig account is allowed")
		return nil
	}

	// get one wallet_address for Authorization account from account_key_authorization table
	authKeyTable, err := w.repo.AccountKey().GetOneMaxID(account.AccountTypeAuthorization)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.GetOneByMaxIDOnAccountKeyTable(AccountTypeAuthorization)")
	}
	// get full-pub-key for given account from multisig_history_table
	multisigHistoryTable, err := w.repo.MultisigHistory().GetAllNoMultisig(accountType)
	if err != nil {
		return errors.Wrapf(err, "fail to call repo.MultisigHistory().GetAllNoMultisig(%s)", accountType.String())
	}

	// call bitcoinAPI `addmultisigaddress`
	//FIXME: for now only 2:2 proportion is available
	// - however N:M should be adjustable
	for _, val := range multisigHistoryTable {
		resAddr, err := w.btc.AddMultisigAddress(
			2,
			[]string{
				val.FullPublicKey,              // deposit, payment, stored ...
				authKeyTable.P2SHSegwitAddress, //TODO: what if address is changed to authKeyTable.FullPublicKey??
			},
			fmt.Sprintf("multi_%s", accountType), //this is not important
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
