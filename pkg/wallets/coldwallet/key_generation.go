package coldwallet

//cold wallet (keygen, sing)

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/coin"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
)

//1. generate seed and store it in database
//2. generate multisig key and store it in database
//    auth device should be separated into each auth accounts (auth1, auth2 ...)
//3. generate pubkey for client and store it in database
//4.Receipt Keyの生成 + Multisig対応 + DBに登録 (1日1Key消費するイメージ)
//5.Payment Keyの生成+ Multisig + DBに登録 (1日1Key消費するイメージ)

// GenerateSeed generate seed and store it in database
func (w *ColdWallet) GenerateSeed() ([]byte, error) {

	// retrieve seed from database
	bSeed, err := w.retrieveSeed()
	if err == nil {
		return bSeed, nil
	}

	// generate seed
	bSeed, err = key.GenerateSeed()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.GenerateSeed()")
	}
	strSeed := key.SeedToString(bSeed)

	// insert seed in database
	_, err = w.storager.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.InsertSeed()")
	}

	return bSeed, nil
}

// store given seed from command line args
//  development use
func (w *ColdWallet) StoreSeed(strSeed string) ([]byte, error) {
	bSeed, err := key.SeedToByte(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.SeedToByte() ")
	}

	// insert seed in database
	_, err = w.storager.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.InsertSeed()")
	}

	return bSeed, nil
}

// retrieve seed from database
func (w *ColdWallet) retrieveSeed() ([]byte, error) {
	// get seed from database, seed is expected only one record
	seed, err := w.storager.GetSeedOne()
	if err == nil && seed.Seed != "" {
		w.logger.Info("seed have already been generated")
		return key.SeedToByte(seed.Seed)
	}
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.GetSeedOne()")
	}
	// in this case, though err didn't happen, but seed is blank
	return nil, errors.New("somehow seed retrieved from database is blank ")
}

// GeneratePubKey generate pubkey for account
// TODO: if account is AccountTypeAuthorization and there is already record, it should stop creation
func (w *ColdWallet) GeneratePubKey(
	accountType account.AccountType,
	coinType coin.CoinType,
	seed []byte, count uint32) ([]key.WalletKey, error) {

	//get latest index
	idxFrom, err := w.storager.GetMaxIndexOnAccountKeyTable(accountType)
	if err != nil {
		idxFrom = 0
	} else {
		idxFrom++
	}

	// generate hd wallet key
	walletKeys, err := w.generateHDKey(accountType, coinType, seed, uint32(idxFrom), count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.generateAccountKeyData()")
	}

	// insert key information to account_key_table
	accountKeyClients := make([]coldrepo.AccountKeyTable, len(walletKeys))
	for idx, key := range walletKeys {
		accountKeyClients[idx] = coldrepo.AccountKeyTable{
			WalletAddress:         key.Address,
			P2shSegwitAddress:     key.P2shSegwit,
			FullPublicKey:         key.FullPubKey,
			WalletMultisigAddress: "",
			RedeemScript:          key.RedeemScript,
			WalletImportFormat:    key.WIF,
			Account:               accountType.String(),
			Idx:                   uint32(idxFrom),
		}
		idxFrom++
	}
	err = w.storager.InsertAccountKeyTable(accountType, accountKeyClients, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call storager.InsertAccountKeyTable()")
	}

	return walletKeys, err
}

func (w *ColdWallet) generateHDKey(
	accountType account.AccountType,
	coinType coin.CoinType,
	seed []byte,
	idxFrom,
	count uint32) ([]key.WalletKey, error) {

	// key object
	//TODO: interface and repository
	keyData := key.NewKey(key.PurposeTypeBIP44, coinType, w.btc.GetChainConf(), w.logger)

	// generate key
	walletKeys, err := keyData.CreateKey(seed, accountType, idxFrom, count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call keyData.CreateKey()")
	}
	return walletKeys, nil
}
