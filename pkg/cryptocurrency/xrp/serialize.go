// Package xrp provides XRP (Ripple) cryptocurrency utilities including
// key generation, address encoding, and cryptographic operations.
package xrp

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// XRP Serialization Format
// Reference: https://xrpl.org/serialization.html

// Field type codes for XRP serialization
const (
	// Type codes (uint8 for standard field types)
	typeUInt16    uint8 = 1
	typeUInt32    uint8 = 2
	typeUInt64    uint8 = 3
	typeHash256   uint8 = 5
	typeAmount    uint8 = 6
	typeBlob      uint8 = 7
	typeAccountID uint8 = 8
	typeSTObject  uint8 = 14
	typeSTArray   uint8 = 15
	typeVector256 uint8 = 19
)

// HashPrefix values for XRP transaction hashing
// Reference: https://xrpl.org/transaction-basics.html
const (
	HashPrefixTransaction         = 0x54584E00 // 'TXN' + 0x00, for unsigned transactions
	HashPrefixTransactionSig      = 0x53545800 // 'STX' + 0x00, for signing
	HashPrefixTransactionID       = 0x54584E00 // 'TXN' + 0x00, for transaction ID
	HashPrefixTransactionMultiSig = 0x534D5400 // 'SMT' + 0x00, for multisig
)

// fieldInfo contains the type code and field code for serialization
type fieldInfo struct {
	TypeCode  uint8
	FieldCode uint8
	IsVL      bool // Variable length encoding
	IsSigning bool // Include in signing data
}

// Transaction field definitions
// Reference: https://xrpl.org/serialization.html#definitions-fields
var fieldDefinitions = map[string]fieldInfo{
	// Common transaction fields
	"TransactionType":    {typeUInt16, 2, false, true},
	"Flags":              {typeUInt32, 2, false, true},
	"Sequence":           {typeUInt32, 4, false, true},
	"LastLedgerSequence": {typeUInt32, 27, false, true},
	"Amount":             {typeAmount, 1, false, true},
	"Fee":                {typeAmount, 8, false, true},
	"SigningPubKey":      {typeBlob, 3, true, false},
	"TxnSignature":       {typeBlob, 4, true, false},
	"Account":            {typeAccountID, 1, false, true},
	"Destination":        {typeAccountID, 3, false, true},
	"SourceTag":          {typeUInt32, 3, false, true},
	"DestinationTag":     {typeUInt32, 14, false, true},

	// SignerListSet fields
	"SignerQuorum":  {typeUInt32, 32, false, true},
	"SignerWeight":  {typeUInt16, 3, false, true},
	"SignerEntry":   {typeSTObject, 16, false, true},
	"SignerEntries": {typeSTArray, 4, false, true},

	// SetRegularKey fields
	"RegularKey": {typeAccountID, 8, false, true},

	// AccountSet fields
	"SetFlag":   {typeUInt32, 33, false, true},
	"ClearFlag": {typeUInt32, 34, false, true},

	// Escrow fields
	"Owner":         {typeAccountID, 2, false, true},
	"OfferSequence": {typeUInt32, 25, false, true},
	"CancelAfter":   {typeUInt32, 36, false, true},
	"FinishAfter":   {typeUInt32, 37, false, true},
	"Condition":     {typeBlob, 17, true, true},
	"Fulfillment":   {typeBlob, 16, true, true},

	// Payment Channel fields
	"SettleDelay": {typeUInt32, 39, false, true},
	"PublicKey":   {typeBlob, 19, true, true},
	"Expiration":  {typeUInt32, 10, false, true},
	"Channel":     {typeHash256, 22, false, true},
	"Balance":     {typeAmount, 5, false, true},
	"Signature":   {typeBlob, 21, true, true},

	// NFToken fields
	"NFTokenTaxon":     {typeUInt32, 42, false, true},
	"TransferFee":      {typeUInt16, 41, false, true},
	"Issuer":           {typeAccountID, 4, false, true},
	"URI":              {typeBlob, 18, true, true},
	"NFTokenID":        {typeHash256, 10, false, true},
	"NFTokenSellOffer": {typeHash256, 29, false, true},
	"NFTokenBuyOffer":  {typeHash256, 30, false, true},
	"NFTokenBrokerFee": {typeAmount, 9, false, true},
	"NFTokenOffers":    {typeVector256, 4, false, true},

	// TrustSet fields
	"LimitAmount": {typeAmount, 3, false, true},
	"QualityIn":   {typeUInt32, 20, false, true},
	"QualityOut":  {typeUInt32, 21, false, true},
}

// Transaction type codes
// Reference: https://xrpl.org/transaction-types.html
var transactionTypes = map[string]uint16{
	"Payment":              0,
	"EscrowCreate":         1,
	"EscrowFinish":         2,
	"AccountSet":           3,
	"EscrowCancel":         4,
	"RegularKeySet":        5, // "SetRegularKey" in JSON
	"SetRegularKey":        5,
	"NickNameSet":          6,
	"OfferCreate":          7,
	"OfferCancel":          8,
	"Contract":             9,
	"TicketCreate":         10,
	"SignerListSet":        12,
	"PaymentChannelCreate": 13,
	"PaymentChannelFund":   14,
	"PaymentChannelClaim":  15,
	"CheckCreate":          16,
	"CheckCash":            17,
	"CheckCancel":          18,
	"DepositPreauth":       19,
	"TrustSet":             20,
	"AccountDelete":        21,
	"NFTokenMint":          25,
	"NFTokenBurn":          26,
	"NFTokenCreateOffer":   27,
	"NFTokenCancelOffer":   28,
	"NFTokenAcceptOffer":   29,
}

// SerializedField represents a field ready for binary serialization
type SerializedField struct {
	Name      string
	FieldID   []byte
	Value     []byte
	TypeCode  uint8
	FieldCode uint8
}

// Serializer handles XRP transaction binary serialization
type Serializer struct{}

// NewSerializer creates a new XRP transaction serializer
func NewSerializer() *Serializer {
	return &Serializer{}
}

// encodeFieldID encodes the field type and field code into bytes
// Reference: https://xrpl.org/serialization.html#field-ids
func encodeFieldID(typeCode, fieldCode uint8) []byte {
	// If both type and field fit in 4 bits each, use single byte
	if typeCode < 16 && fieldCode < 16 {
		return []byte{(typeCode << 4) | fieldCode}
	}
	// If type fits in 4 bits but field doesn't
	if typeCode < 16 {
		return []byte{typeCode << 4, fieldCode}
	}
	// If field fits in 4 bits but type doesn't
	if fieldCode < 16 {
		return []byte{fieldCode, typeCode}
	}
	// Neither fits in 4 bits
	return []byte{0, typeCode, fieldCode}
}

// encodeVL encodes a variable-length prefix
// Reference: https://xrpl.org/serialization.html#length-prefixing
func encodeVL(length int) []byte {
	if length <= 192 {
		return []byte{byte(length)}
	}
	if length <= 12480 {
		length -= 193
		return []byte{byte(193 + (length >> 8)), byte(length & 0xFF)}
	}
	if length <= 918744 {
		length -= 12481
		return []byte{
			byte(241 + (length >> 16)),
			byte((length >> 8) & 0xFF),
			byte(length & 0xFF),
		}
	}
	// Should not happen for valid XRP transactions
	return nil
}

// encodeUInt16 encodes a uint16 in big-endian format
func encodeUInt16(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}

// encodeUInt32 encodes a uint32 in big-endian format
func encodeUInt32(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf
}

// encodeUInt64 encodes a uint64 in big-endian format
func encodeUInt64(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf
}

// encodeXRPAmount encodes an XRP amount (in drops) for serialization
// Reference: https://xrpl.org/serialization.html#amount-fields
func encodeXRPAmount(drops string) ([]byte, error) {
	// Parse the drops value
	val, err := strconv.ParseUint(drops, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid XRP amount: %w", err)
	}

	// XRP amounts have bit 63 set to indicate native XRP (not IOU)
	// and bit 62 set to indicate positive
	encoded := val | 0x4000000000000000
	return encodeUInt64(encoded), nil
}

// decodeAccountAddress decodes an XRP address (r...) to 20-byte account ID
func decodeAccountAddress(address string) ([]byte, error) {
	if !strings.HasPrefix(address, "r") {
		return nil, errors.New("invalid XRP address format: must start with 'r'")
	}

	hash, err := NewRippleHashCheck(address, RippleAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode address %s: %w", address, err)
	}

	return hash.Payload(), nil
}

// encodeAccountID encodes an XRP address for serialization
func encodeAccountID(address string) ([]byte, error) {
	accountID, err := decodeAccountAddress(address)
	if err != nil {
		return nil, err
	}

	// Account ID is VL encoded (20 bytes)
	vl := encodeVL(len(accountID))
	return append(vl, accountID...), nil
}

// encodeBlob encodes binary data with VL prefix
func encodeBlob(data []byte) []byte {
	vl := encodeVL(len(data))
	return append(vl, data...)
}

// encodeHexBlob encodes hex-encoded data with VL prefix
func encodeHexBlob(hexStr string) ([]byte, error) {
	if hexStr == "" {
		return encodeBlob([]byte{}), nil
	}
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex data: %w", err)
	}
	return encodeBlob(data), nil
}

// encodeHash256 encodes a 256-bit hash from hex string
func encodeHash256(hexStr string) ([]byte, error) {
	if hexStr == "" {
		return make([]byte, 32), nil
	}
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hash256: %w", err)
	}
	if len(data) != 32 {
		return nil, fmt.Errorf("hash256 must be 32 bytes, got %d", len(data))
	}
	return data, nil
}

// Transaction represents a generic XRP transaction for serialization
type Transaction struct {
	Fields map[string]any
}

// NewTransaction creates a new transaction from a map of fields
func NewTransaction(fields map[string]any) *Transaction {
	return &Transaction{Fields: fields}
}

// serializeField serializes a single transaction field
func (*Serializer) serializeField(name string, value any, forSigning bool) (*SerializedField, error) {
	info, found := fieldDefinitions[name]
	if !found {
		return nil, fmt.Errorf("unknown field: %s", name)
	}

	// Skip non-signing fields when serializing for signing
	if forSigning && !info.IsSigning {
		return nil, nil
	}

	fieldID := encodeFieldID(info.TypeCode, info.FieldCode)

	encoded, err := encodeFieldValue(name, value, info.TypeCode)
	if err != nil {
		return nil, err
	}

	return &SerializedField{
		Name:      name,
		FieldID:   fieldID,
		Value:     encoded,
		TypeCode:  info.TypeCode,
		FieldCode: info.FieldCode,
	}, nil
}

// encodeFieldValue encodes a field value based on its type code
func encodeFieldValue(name string, value any, typeCode uint8) ([]byte, error) {
	switch typeCode {
	case typeUInt16:
		v, valid := toUint64(value)
		if !valid {
			return nil, fmt.Errorf("field %s: expected uint16 value", name)
		}
		return encodeUInt16(uint16(v)), nil

	case typeUInt32:
		v, valid := toUint64(value)
		if !valid {
			return nil, fmt.Errorf("field %s: expected uint32 value", name)
		}
		return encodeUInt32(uint32(v)), nil

	case typeUInt64:
		v, valid := toUint64(value)
		if !valid {
			return nil, fmt.Errorf("field %s: expected uint64 value", name)
		}
		return encodeUInt64(v), nil

	case typeAmount:
		str, valid := value.(string)
		if !valid {
			return nil, fmt.Errorf("field %s: expected string amount", name)
		}
		return encodeXRPAmount(str)

	case typeBlob:
		str, valid := value.(string)
		if !valid {
			return nil, fmt.Errorf("field %s: expected hex string", name)
		}
		return encodeHexBlob(str)

	case typeAccountID:
		str, valid := value.(string)
		if !valid {
			return nil, fmt.Errorf("field %s: expected address string", name)
		}
		return encodeAccountID(str)

	case typeHash256:
		str, valid := value.(string)
		if !valid {
			return nil, fmt.Errorf("field %s: expected hex string", name)
		}
		return encodeHash256(str)

	default:
		return nil, fmt.Errorf("unsupported field type for %s: %d", name, typeCode)
	}
}

// toUint64 converts various numeric types to uint64.
// Note: float64 is intentionally not supported to enforce type safety
// at the call site. Transaction fields are expected to be integers.
func toUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case uint64:
		return v, true
	case uint32:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case int:
		return uint64(v), true
	default:
		return 0, false
	}
}

// Serialize serializes a transaction to binary format
// If forSigning is true, excludes SigningPubKey and TxnSignature fields
func (s *Serializer) Serialize(tx *Transaction, forSigning bool) ([]byte, error) {
	if tx == nil || tx.Fields == nil {
		return nil, errors.New("transaction is nil")
	}

	// Collect and serialize all fields
	serializedFields, err := s.serializeAllFields(tx, forSigning)
	if err != nil {
		return nil, err
	}

	// Sort fields by type code, then by field code (canonical ordering)
	sort.Slice(serializedFields, func(i, j int) bool {
		if serializedFields[i].TypeCode != serializedFields[j].TypeCode {
			return serializedFields[i].TypeCode < serializedFields[j].TypeCode
		}
		return serializedFields[i].FieldCode < serializedFields[j].FieldCode
	})

	// Concatenate all serialized fields
	var result []byte
	for _, field := range serializedFields {
		result = append(result, field.FieldID...)
		result = append(result, field.Value...)
	}

	return result, nil
}

// serializeAllFields collects and serializes all transaction fields
func (s *Serializer) serializeAllFields(tx *Transaction, forSigning bool) ([]*SerializedField, error) {
	var serializedFields []*SerializedField

	// Handle TransactionType specially - convert string to uint16
	if txType, ok := tx.Fields["TransactionType"].(string); ok {
		typeCode, found := transactionTypes[txType]
		if !found {
			return nil, fmt.Errorf("unknown transaction type: %s", txType)
		}
		field, err := s.serializeField("TransactionType", typeCode, forSigning)
		if err != nil {
			return nil, err
		}
		if field != nil {
			serializedFields = append(serializedFields, field)
		}
	}

	// Serialize other fields
	for name, value := range tx.Fields {
		if shouldSkipField(name, value) {
			continue
		}

		field, err := s.serializeField(name, value, forSigning)
		if err != nil {
			return nil, err
		}
		if field != nil {
			serializedFields = append(serializedFields, field)
		}
	}

	return serializedFields, nil
}

// shouldSkipField returns true if a field should be skipped during serialization
func shouldSkipField(name string, value any) bool {
	// Skip already handled or non-serializable fields
	if name == "TransactionType" || name == "Hash" {
		return true
	}

	// Skip nil values
	if value == nil {
		return true
	}

	// Skip empty strings
	if str, ok := value.(string); ok && str == "" {
		return true
	}

	// Skip zero values for optional fields
	// Use toUint64 to handle all numeric types (uint32, uint64, etc.)
	if num, ok := toUint64(value); ok && num == 0 {
		if name == "Flags" || name == "DestinationTag" || name == "SourceTag" {
			return true
		}
	}

	return false
}

// SerializeForSigning serializes a transaction for signing (excludes signature fields)
func (s *Serializer) SerializeForSigning(tx *Transaction) ([]byte, error) {
	return s.Serialize(tx, true)
}

// SerializeForSubmission serializes a transaction for network submission
func (s *Serializer) SerializeForSubmission(tx *Transaction) ([]byte, error) {
	return s.Serialize(tx, false)
}

// ComputeSigningHash computes the hash to be signed for a transaction
// Reference: https://xrpl.org/transaction-basics.html#signing-and-submitting-transactions
func (s *Serializer) ComputeSigningHash(tx *Transaction) ([]byte, error) {
	// Serialize the transaction for signing
	serialized, err := s.SerializeForSigning(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// Prepend the signing prefix
	data := make([]byte, 4, 4+len(serialized))
	binary.BigEndian.PutUint32(data, HashPrefixTransactionSig)
	data = append(data, serialized...)

	// Return SHA512Half of the prefixed data
	return Sha512Half(data), nil
}

// ComputeTransactionID computes the transaction ID from a signed transaction blob
func ComputeTransactionID(txBlob []byte) string {
	// Prepend the transaction ID prefix
	data := make([]byte, 4, 4+len(txBlob))
	binary.BigEndian.PutUint32(data, HashPrefixTransactionID)
	data = append(data, txBlob...)

	// Return SHA512Half as hex
	hash := Sha512Half(data)
	return strings.ToUpper(hex.EncodeToString(hash))
}

// ComputeTransactionIDFromHex computes the transaction ID from a hex-encoded signed transaction
func ComputeTransactionIDFromHex(txBlobHex string) (string, error) {
	txBlob, err := hex.DecodeString(txBlobHex)
	if err != nil {
		return "", fmt.Errorf("invalid tx blob hex: %w", err)
	}
	return ComputeTransactionID(txBlob), nil
}
