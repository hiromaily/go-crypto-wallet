package coldsrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
)

// Seed type
type Seed struct {
	logger   *zap.Logger
	seedRepo coldrepo.SeedRepositorier
	wtype    wallet.WalletType
}

// NewSeed returns seed object
func NewSeed(
	logger *zap.Logger,
	seedRepo coldrepo.SeedRepositorier,
	wtype wallet.WalletType) *Seed {

	return &Seed{
		logger:   logger,
		seedRepo: seedRepo,
		wtype:    wtype,
	}
}

// Generate generate seed and store it in database
func (s *Seed) Generate() ([]byte, error) {

	// retrieve seed from database
	bSeed, err := s.retrieveSeed()
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
	err = s.seedRepo.Insert(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.Seed().Insert()")
	}

	return bSeed, nil
}

// Store stores given seed from command line args
//  development use
func (s *Seed) Store(strSeed string) ([]byte, error) {
	bSeed, err := key.SeedToByte(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.SeedToByte() ")
	}

	// insert seed in database
	err = s.seedRepo.Insert(strSeed)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.InsertSeed()")
	}

	return bSeed, nil
}

// retrieve seed from database
func (s *Seed) retrieveSeed() ([]byte, error) {
	// get seed from database, seed is expected only one record
	seed, err := s.seedRepo.GetOne()
	if err == nil && seed.Seed != "" {
		s.logger.Info("seed have already been generated")
		return key.SeedToByte(seed.Seed)
	}
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.GetSeedOne()")
	}
	// in this case, though err didn't happen, but seed is blank
	return nil, errors.New("somehow seed retrieved from database is blank ")
}
