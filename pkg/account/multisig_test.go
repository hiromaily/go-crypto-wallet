package account

import (
	"testing"

	"github.com/stretchr/testify/require"

	configutil "github.com/hiromaily/go-crypto-wallet/pkg/config/testutil"
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
	confPath := configutil.GetConfigFilePath("account.toml")
	// projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	// confPath := fmt.Sprintf("%s/data/config/account.toml", projPath)
	conf, err := NewAccount(confPath)
	require.NoError(t, err, "fail to create config")

	multi := NewMultisigAccounts(conf.Multisigs)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := multi.IsMultisigAccount(tt.args.acnt)
			require.Equal(t, tt.want.ok, res, "IsMultisigAccount() result mismatch")
		})
	}
}
