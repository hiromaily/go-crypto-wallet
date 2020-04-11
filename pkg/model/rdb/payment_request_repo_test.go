package rdb_test

import (
	"testing"

	. "github.com/hiromaily/go-bitcoin/pkg/model"
)

func TestGetPaymentRequest(t *testing.T) {
	paymentRequests, err := db.GetPaymentRequestAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(paymentRequests)
}

func TestInsertPaymentRequest(t *testing.T) {

	paymentRequests := []PaymentRequest{
		{
			AddressFrom: "2MuQ83G8hmCnz1bSiqKx4koKbNCptL39k24",
			AccountFrom: "hiroki",
			AddressTo:   "2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz",
			Amount:      "0.5",
			//IsDone:      false,
		},
		{
			AddressFrom: "2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf",
			AccountFrom: "yasui",
			AddressTo:   "2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz",
			Amount:      "0.35",
			//IsDone:      false,
		},
		{
			AddressFrom: "2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf",
			AccountFrom: "yasui",
			AddressTo:   "2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV",
			Amount:      "0.45",
			//IsDone:      false,
		},
	}

	err := db.InsertPaymentRequest(paymentRequests, nil, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePaymentIDOnPaymentRequest(t *testing.T) {

	affected, err := db.UpdatePaymentIDOnPaymentRequest(999, []int64{1, 2, 3}, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(affected)
}

func TestUpdateIsDoneOnPaymentRequest(t *testing.T) {
	affected, err := db.UpdateIsDoneOnPaymentRequest(999, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(affected)
}
