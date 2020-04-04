package btc_test

import "testing"

// TestLockUnspent
func TestLockUnspent(t *testing.T) {
	//1.　ロックされているものがあれば、最初に解除
	if err := wlt.BTC.UnlockAllUnspentTransaction(); err != nil {
		t.Fatalf("UnlockAllUnspentTransaction() error: %v", err)
	}
	//2. Unspentのutxoを取得
	list, err := wlt.BTC.Client().ListUnspentMin(6)
	if err != nil {
		t.Fatalf("ListUnspentMin() error: %v", err)
	}
	if len(list) == 0 {
		t.Logf("unspent data is zero. first of all, test data should be created for testing by `%s`", "https://testnet.manu.backend.hamburg/faucet")
		return
	}
	firstLen := len(list)
	t.Logf("target tx: %v", list[0])
	//3. 最初の1つのみ、Lock
	if wlt.BTC.LockUnspent(list[0]) != nil {
		t.Fatalf("LockUnspent() error: %v", err)
	}
	//4. 再度Unspentで取得
	list2, err := wlt.BTC.Client().ListUnspentMin(6)
	if err != nil {
		t.Fatalf("2nd ListUnspentMin() error: %v", err)
	}
	//5. 3.にてロックされたことによって、len(list)が減算されたことを確認
	if len(list2) != firstLen-1 {
		t.Fatalf("length is wrong. first:%d, second:%d", firstLen, len(list2))
	}
}
