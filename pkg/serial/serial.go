package serial

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// EncodeToString is binary encoder
func EncodeToString(data interface{}) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// DecodeFromString is binary decoder
func DecodeFromString(str string, data interface{}) error {
	//u := User{}
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
