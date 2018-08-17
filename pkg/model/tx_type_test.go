package model_test

import (
	"testing"
)

func TestGetTxType(t *testing.T) {
	txTypes, err := db.GetTxType()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txTypes)
}
