package service

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/sign"
)

//-----------------------------------------------------------------------------
// Cold Wallet Services - Type aliases for backward compatibility
// These interfaces have been split and moved to:
// - Keygen interfaces: pkg/wallet/service/keygen/interfaces.go
// - Sign interfaces: pkg/wallet/service/sign/interfaces.go
//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------
// Keygen / Sign
//-----------------------------------------------------------------------------

// HDWalleter is HD wallet key generation service
type HDWalleter = keygen.HDWalleter

// Seeder is Seeder service
type Seeder = keygen.Seeder

// Signer is Signer service
type Signer = sign.Signer

//-----------------------------------------------------------------------------
// keygen
//-----------------------------------------------------------------------------

// AddressExporter is AddressExporter service
type AddressExporter = keygen.AddressExporter

// PrivKeyer is PrivKeyer service
type PrivKeyer = keygen.PrivKeyer

// FullPubKeyImporter is FullPubkeyImport service
type FullPubKeyImporter = keygen.FullPubKeyImporter

// Multisiger is Multisiger service
type Multisiger = keygen.Multisiger

//-----------------------------------------------------------------------------
// Sign
//-----------------------------------------------------------------------------

// FullPubkeyExporter is FullPubkeyExporter service
type FullPubkeyExporter = sign.FullPubkeyExporter
