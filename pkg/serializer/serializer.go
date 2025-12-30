package serializer

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
)

// Serializer is an interface for serializing and deserializing data
type Serializer interface {
	EncodeToString(data any) (string, error)
	DecodeFromString(str string, data any) error
}

// GobSerializer implements Serializer using gob encoding with base64
type GobSerializer struct{}

// NewGobSerializer creates a new GobSerializer
func NewGobSerializer() Serializer {
	return &GobSerializer{}
}

// EncodeToString encodes data to string using gob and base64
func (*GobSerializer) EncodeToString(data any) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// DecodeFromString decodes string to data using base64 and gob
func (*GobSerializer) DecodeFromString(str string, data any) error {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// JSONSerializer implements Serializer using JSON encoding
type JSONSerializer struct{}

// NewJSONSerializer creates a new JSONSerializer
func NewJSONSerializer() Serializer {
	return &JSONSerializer{}
}

// EncodeToString encodes data to string using JSON
func (*JSONSerializer) EncodeToString(data any) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// DecodeFromString decodes string to data using JSON
func (*JSONSerializer) DecodeFromString(str string, data any) error {
	return json.Unmarshal([]byte(str), data)
}

// defaultSerializer is the default serializer (gob) for backward compatibility
var defaultSerializer Serializer = NewGobSerializer()

// SetDefaultSerializer sets the default serializer
func SetDefaultSerializer(s Serializer) {
	defaultSerializer = s
}

// GetDefaultSerializer returns the default serializer
func GetDefaultSerializer() Serializer {
	return defaultSerializer
}
