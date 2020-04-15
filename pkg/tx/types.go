package tx

//----------------------------------------------------
// TxType
//----------------------------------------------------

//TxType トランザクション種別
type TxType string

// tx_type
const (
	TxTypeUnsigned    TxType = "unsigned"
	TxTypeUnsigned2nd TxType = "unsigned2nd"
	TxTypeSigned      TxType = "signed"
	TxTypeSent        TxType = "sent"
	TxTypeDone        TxType = "done"
	TxTypeNotified    TxType = "notified"
	TxTypeCancel      TxType = "canceled"
)

func (t TxType) String() string {
	return string(t)
}

//TxTypeValue tx_typeの値
var TxTypeValue = map[TxType]uint8{
	TxTypeUnsigned:    1,
	TxTypeUnsigned2nd: 2,
	TxTypeSigned:      3,
	TxTypeSent:        4,
	TxTypeDone:        5,
	TxTypeNotified:    6,
	TxTypeCancel:      7,
}

// ValidateTxType TxTypeのバリデーションを行う
func ValidateTxType(val string) bool {
	if _, ok := TxTypeValue[TxType(val)]; ok {
		return true
	}
	return false
}

// Search SliceのtxTypes内に`t`が含まれているかチェックする
func (t TxType) Search(txTypes []TxType) bool {
	for _, v := range txTypes {
		if v == t {
			return true
		}
	}
	return false
}
