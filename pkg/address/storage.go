package address

import "github.com/hiromaily/go-bitcoin/pkg/account"

type Storager interface {
	CreateFilePath(accountType account.AccountType, keyStatus uint8) string
	ValidateFilePath(fileName string, accountType account.AccountType) error
	ImportPubKey(fileName string) ([]string, error)
}

//type Importer interface {
//}
//
//type Exporter interface {
//}
