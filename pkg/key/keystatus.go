package key

//----------------------------------------------------
// KeyStatus
//----------------------------------------------------

//KeyStatus Key生成進捗ステータス
type KeyStatus string

// key_status
const (
	KeyStatusGenerated            KeyStatus = "generated"              //hd_walletによってkeyが生成された
	KeyStatusImportprivkey        KeyStatus = "importprivkey"          //importprivkeyが実行された
	KeyStatusPubkeyExported       KeyStatus = "pubkey_exported"        //pubkeyがexportされた(receipt/payment)
	KeyStatusMultiAddressImported KeyStatus = "multi_address_imported" //multiaddがimportされた(receipt/payment)
	KeyStatusAddressExported      KeyStatus = "address_exported"       //addressがexportされた
)

func (k KeyStatus) String() string {
	return string(k)
}

//KeyStatusValue key_statusの値
var KeyStatusValue = map[KeyStatus]uint8{
	KeyStatusGenerated:            0,
	KeyStatusImportprivkey:        1,
	KeyStatusPubkeyExported:       2,
	KeyStatusMultiAddressImported: 3,
	KeyStatusAddressExported:      4,
}

// ValidateKeyStatus KeyStatusのバリデーションを行う
func ValidateKeyStatus(val string) bool {
	if _, ok := KeyStatusValue[KeyStatus(val)]; ok {
		return true
	}
	return false
}
