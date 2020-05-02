package keygensrv

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// Multisiger is Multisiger service
type Multisiger interface {
	AddMultisigAddress(accountType account.AccountType, addressType address.AddrType) error
}

// Multisig type
type Multisig struct {
	btc                api.Bitcoiner
	logger             *zap.Logger
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier
	accountKeyRepo     coldrepo.AccountKeyRepositorier
	wtype              wallet.WalletType
}

// NewMultisig returns multisig
func NewMultisig(
	btc api.Bitcoiner,
	logger *zap.Logger,
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	wtype wallet.WalletType) *Multisig {

	return &Multisig{
		btc:                btc,
		logger:             logger,
		authFullPubKeyRepo: authFullPubKeyRepo,
		accountKeyRepo:     accountKeyRepo,
		wtype:              wtype,
	}
}

// AddMultisigAddress add multisig address by auth account address and given account address
// https://bitcoincore.org/en/doc/0.19.0/rpc/wallet/addmultisigaddress/
// - if 3:5 proportion is required, at least 4 different auth accounts should be prepared in advance
// - when sending coin from multisig address, 、related priv key is required which is related to addresses in parameters
// - 4th parameter must be`p2sh-segwit` addressType in Bitcoin
//  What is the difference between createmultisig and addmultisigaddress?
// - https://bitcointalk.org/index.php?topic=3402541.0
func (m *Multisig) AddMultisigAddress(accountType account.AccountType, addressType address.AddrType) error {
	//for sign wallet
	m.logger.Debug("addmultisigaddress",
		zap.String("account_type", accountType.String()),
	)

	// validate accountType
	if !account.IsMultisigAccount(accountType) {
		m.logger.Info("only multisig account is allowed")
		return nil
	}

	var requiredSig int
	var authFullPubKeys []string
	// get fullPubKey from auth_fullpubkey table

	//AccountTypeDeposit: { //2:5+1
	//	2: {AuthType1, AuthType2, AuthType3, AuthType4, AuthType5},
	//},
	for sigCount, authTypes := range account.MultisigAccounts[accountType] {
		requiredSig = sigCount
		for _, authType := range authTypes {
			// get record from
			fullPubKeyItem, err := m.authFullPubKeyRepo.GetOne(authType)
			if err != nil {
				return errors.Wrapf(err, "fail to call authFullPubKeyRepo.GetOne() %s", authType.String())
			}
			authFullPubKeys = append(authFullPubKeys, fullPubKeyItem.FullPublicKey)
		}
		m.logger.Debug("don't repeat again")
	}

	// get target addresses from account_key table, addr_status=AddrStatusPrivKeyImported
	accountKeyItems, err := m.accountKeyRepo.GetAllAddrStatus(accountType, address.AddrStatusPrivKeyImported)
	if err != nil {
		return errors.Wrapf(err, "fail to call accountKeyRepo.GetAllAddrStatus(%s)", accountType.String())
	}

	// call bitcoinAPI `addmultisigaddress`
	for _, item := range accountKeyItems {
		addrs := append(authFullPubKeys, item.FullPublicKey)
		resAddr, err := m.btc.AddMultisigAddress(
			requiredSig,
			addrs,
			fmt.Sprintf("multi_%s", accountType), //this is not important
			addressType,
		)
		if err != nil {
			//[Error] -5: no full public key for address mkPmdpo59gpU7ZioGYwwoMTQJjh7MiqUvd
			m.logger.Error(
				"fail to call btc.CreateMultiSig()",
				zap.Int("signature_count", requiredSig),
				zap.String("full public key for accountType", item.FullPublicKey),
				zap.Strings("full public key for authType", authFullPubKeys),
				zap.Error(err),
			)
			continue
		}
		// update generated multisig address, redeemScript, addrStatus
		item.MultisigAddress = resAddr.Address
		item.RedeemScript = resAddr.RedeemScript
		item.AddrStatus = address.AddrStatusMultisigAddressGenerated.Int8()

		_, err = m.accountKeyRepo.UpdateMultisigAddr(accountType, item)
		if err != nil {
			return errors.Wrapf(err, "fail to call accountKeyRepo.UpdateMultisigAddr(%s)", accountType.String())
		}
	}

	return nil
}