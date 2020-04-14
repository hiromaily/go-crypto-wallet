package storage

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
)

type ImportExporter interface {
	ImportPubKey(fileName string, accountType account.AccountType) error
}

type Importer interface {
}

type Exporter interface {
}
