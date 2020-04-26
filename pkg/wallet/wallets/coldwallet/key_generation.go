package coldwallet

//cold wallet (keygen, sing)

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
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
	err = w.repo.Seed().Insert(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.Seed().Insert()")
	}

	return bSeed, nil
}

// StoreSeed stores given seed from command line args
//  development use
func (w *ColdWallet) StoreSeed(strSeed string) ([]byte, error) {
	bSeed, err := key.SeedToByte(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.SeedToByte() ")
	}

	// insert seed in database
	err = w.repo.Seed().Insert(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.InsertSeed()")
	}

	return bSeed, nil
}

// retrieve seed from database
func (w *ColdWallet) retrieveSeed() ([]byte, error) {
	// get seed from database, seed is expected only one record
	seed, err := w.repo.Seed().GetOne()
	if err == nil && seed.Seed != "" {
		w.logger.Info("seed have already been generated")
		return key.SeedToByte(seed.Seed)
	}
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.GetSeedOne()")
	}
	// in this case, though err didn't happen, but seed is blank
	return nil, errors.New("somehow seed retrieved from database is blank ")
}

// GeneratePubKey generate pubkey for account
// TODO: if account is AccountTypeAuthorization and there is already record, it should stop creation
func (w *ColdWallet) GeneratePubKey(
	accountType account.AccountType,
	seed []byte, count uint32) ([]key.WalletKey, error) {

	//get latest index
	idxFrom, err := w.repo.AccountKey().GetMaxIndex(accountType)
	if err != nil {
		w.logger.Debug("fail by repo.AccountKey().GetMaxIndex()", zap.Error(err))
		idxFrom = 0
	} else {
		idxFrom++
	}
	w.logger.Debug("max_index",
		zap.String("account_type", accountType.String()),
		zap.Int64("current_index", idxFrom),
	)

	// generate hd wallet key
	walletKeys, err := w.generateHDKey(accountType, seed, uint32(idxFrom), count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.generateAccountKeyData()")
	}

	// insert key information to account_key_table
	accountKeyClients := make([]*models.AccountKey, len(walletKeys))
	for idx, key := range walletKeys {
		accountKeyClients[idx] = &models.AccountKey{
			Coin:                  w.GetBTC().CoinTypeCode().String(),
			Account:               accountType.String(),
			WalletAddress:         key.Address,
			P2SHSegwitAddress:     key.P2shSegwit,
			FullPublicKey:         key.FullPubKey,
			WalletMultisigAddress: "",
			RedeemScript:          key.RedeemScript,
			WalletImportFormat:    key.WIF,
			Idx:                   idxFrom,
		}
		idxFrom++
	}
	err = w.repo.AccountKey().InsertBulk(accountKeyClients)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.AccountKey().InsertBulk()")
	}

	return walletKeys, err
}

func (w *ColdWallet) generateHDKey(
	accountType account.AccountType,
	seed []byte,
	idxFrom,
	count uint32) ([]key.WalletKey, error) {

	// generate key
	walletKeys, err := w.keyGenerator.CreateKey(seed, accountType, idxFrom, count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call keyData.CreateKey()")
	}
	return walletKeys, nil
}
