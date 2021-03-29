package account

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// TestNewMultisigAccounts is test for NewMultisigAccounts
func TestNewMultisigAccounts(t *testing.T) {
	// t.SkipNow()

	type args struct {
		acnt AccountType
	}
	type want struct {
		ok bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path client",
			args: args{AccountTypeClient},
			want: want{false},
		},
		{
			name: "happy path deposit",
			args: args{AccountTypeDeposit},
			want: want{true},
		},
		{
			name: "happy path payment",
			args: args{AccountTypePayment},
			want: want{true},
		},
		{
			name: "happy path stored",
			args: args{AccountTypeStored},
			want: want{true},
		},
		{
			name: "happy path blanc",
			args: args{""},
			want: want{false},
		},
	}

	// config
	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/account.toml", projPath)
	conf, err := NewAccount(confPath)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	multi := NewMultisigAccounts(conf.Multisigs)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := multi.IsMultisigAccount(tt.args.acnt); res != tt.want.ok {
				t.Errorf("IsMultisigAccount() = %t, want %t", res, tt.want.ok)
			}
		})
	}
}
