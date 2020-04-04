package key_test

import (
	"flag"
	"os"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cpacia/bchutil"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	. "github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/service"
)

var (
	wlt      *service.Wallet
	confPath = flag.String("conf", "../../data/toml/local_watch_only.toml", "Path for configuration toml file")
)

func setup() {
	if wlt != nil {
		return
	}

	flag.Parse()

	var err error
	wlt, err = service.InitialSettings(*confPath)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	wlt.Done()
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func generateKeys(bSeed []byte, coinType enum.CoinType, t *testing.T) {
	// key生成
	keyData := NewKey(coinType, wlt.BTC.GetChainConf())

	priv, _, err := keyData.CreateAccount(bSeed, enum.AccountTypeClient)
	if err != nil {
		t.Fatal(err)
	}

	walletKeys, err := keyData.CreateKeysWithIndex(priv, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, k := range walletKeys {
		t.Logf("[Address] %s", k.Address)
		t.Logf("[P2shSegwit] %s", k.P2shSegwit)
		t.Logf("[FullPubKey] %s", k.FullPubKey)
		t.Logf("[WIF] %s", k.WIF)

		if coinType == enum.BCH {
			//Decode
			decodeForBCH(k.Address, t)
		}
	}

}

func decodeForBCH(addr string, t *testing.T) {
	//prefix, ok := bchutil.Prefixes[wlt.BTC.GetChainConf().Name]
	//if !ok {
	//	t.Fatal("*chaincfg.Params is wrong")
	//}
	resPrefix, val, err := bchutil.DecodeCashAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("prefix: %s, val: %v", resPrefix, val)
}

func TestKeyIntegrationBTC(t *testing.T) {
	t.SkipNow()

	// アドレスを生成
	bSeed, err := GenerateSeed()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("[seed] %s", SeedToString(bSeed))
	generateKeys(bSeed, enum.BTC, t)
}

func TestKeyIntegrationBTC2(t *testing.T) {
	t.SkipNow()

	testSeed := "ggqMLyaZ7pwOXRdH8N2CWBf7L9gS/P8/p7oJdjp9M8U="
	bSeed, err := SeedToByte(testSeed)
	if err != nil {
		t.Fatal(err)
	}

	generateKeys(bSeed, enum.BTC, t)

	//expecting always when using same seed
	//key_test.go:70: [Address] n34rJwoiPM6igHLqLUpfZmftc31u42gagZ
	//key_test.go:71: [P2shSegwit] 2Mv6uWSWBYsoLkfCwr8nD6vMBDtjH4kpem2
	//key_test.go:72: [FullPubKey] 023d7e14cd8c3f3682b7c6f83f6df44415e8e35f7d5abe11fa1a076494d7c11830
	//key_test.go:73: [WIF] cUMVTCv4qP9uGc3xGa8ixkAX2WgCfoKYcto2fm2vKo4bzko9GDn6
	//key_test.go:70: [Address] mvDrRKeV7hivvJGbLiiWncgfszADJcgqun
	//key_test.go:71: [P2shSegwit] 2Mxmxh72CopAt7VNLWybHrb1nHgSmvwFBeG
	//key_test.go:72: [FullPubKey] 03b1929341eed9f6b1ab8199f10ea3406cbf3739c19450f9c895f349614228bcd7
	//key_test.go:73: [WIF] cTi8dgxNr6b9SWiUdzBhGMA1czMfHAhCd9ToTctRsXUWX9ifn4ek
	//key_test.go:70: [Address] miE5tPU3ftywzKtgwgKm5Rwx1ifAiKEk3c
	//key_test.go:71: [P2shSegwit] 2MuUNXmMbbogw25Ce1mernMd5xxDNdiSDgZ
	//key_test.go:72: [FullPubKey] 02480debf8b3d46ab2a339b2d793d46d5a238233490c6a7e0ef6416bffb6cdb377
	//key_test.go:73: [WIF] cPRNxLu3EXhrc8qfDAJnewRqYqavuKaLt6GTsJk8HgaHQBZpXmfM
	//key_test.go:70: [Address] mpJTndTvPvQmHHfJPpuXPZDnNf988ZPhDE
	//key_test.go:71: [P2shSegwit] 2N3jMJVFYvm1QzpNYVY8h1yfCEjpxF1HMYZ
	//key_test.go:72: [FullPubKey] 02106d20e9e88f70f190c2bc08d65d38483538859dc2327330eabd8bc537d52af8
	//key_test.go:73: [WIF] cV8FxLpDhqVcDVgSUnbEk1FVuXz2wqxRQvwGC3wXdxJU3W6JRxfk
	//key_test.go:70: [Address] mth7KjuxpcTBK8PfQ8bbfkgo8txDY1knBV
	//key_test.go:71: [P2shSegwit] 2N6DYNCZw2LjGYvXXHqmCs1D6Pe68BSa5F8
	//key_test.go:72: [FullPubKey] 02d477c90cd169eb3762d9c76f175a68192e6fc477194b4e3beb23a7622113c9ae
	//key_test.go:73: [WIF] cRQYRjvmneQfkVqWezDjFJ76H8J2R2DuwFsJnQgVbA4FaszLWbD6
	//key_test.go:70: [Address] mnFqztrm4yu3Pc7vSfobVjrkjy6XSqzeBn
	//key_test.go:71: [P2shSegwit] 2NGTJCpRcMFy1QK8pEh1UxHeMriV5GoAQ9s
	//key_test.go:72: [FullPubKey] 0388c1f97dddf9da68c6461df89a97c873d976bac16db2d8c807f6fac93e899543
	//key_test.go:73: [WIF] cVCUspxA2upT8NM7p12oZgSAirT7vjc3r5aZiEC7MwNEiqL91KLH
	//key_test.go:70: [Address] mpRTEdoAzTb8i5bx3K3zpzLV2VzYtiBpRW
	//key_test.go:71: [P2shSegwit] 2Mzh4UbxNkB6xA7VZyZmufVuPy1sszs53q8
	//key_test.go:72: [FullPubKey] 0367bff62a67d7ae8e09f1d287942ca745d49020ed958dd79f2f9b13e340acac37
	//key_test.go:73: [WIF] cVUuDFQFbAKjWbu6NiP1KJhG7CxzuEPFr8tvcWzgVZfUZmxSRxYL
	//key_test.go:70: [Address] mm8eKi1tV7ScRj694GuxFCbd2CNxkfz7nL
	//key_test.go:71: [P2shSegwit] 2N53gBAC1ub2JMC2BPnGDrC2qF9zsoF98LD
	//key_test.go:72: [FullPubKey] 03d18c7b49f4facdaf1983c9597cf49d9f455e34351b7fb7297fc574ab99677497
	//key_test.go:73: [WIF] cQbECBZ41gMy9HGn1Jpt1TG6ycm5SuciEiXtfZpJ1Cc6fPH5aCRK
	//key_test.go:70: [Address] n1CF17732FwusDcabjzNDrHuKTZHXmQpqy
	//key_test.go:71: [P2shSegwit] 2NCMjLjqnAYSkopNm57NQLBst3UiRarsPTA
	//key_test.go:72: [FullPubKey] 035d29a14480f4ccd07834039c1104225ace0ac31af88b34a17a6f1e0d62b9268a
	//key_test.go:73: [WIF] cRCqapDroRmjMmZHjHyqYYRUNKGdNKHRKWtGma4oo4JcB3b4r1EF
	//key_test.go:70: [Address] mgkd6VtDcntYnfbd8eqCEfE68FmrQ2HMuy
	//key_test.go:71: [P2shSegwit] 2Mvbe5Nyns2SXW6TvZyoSABUQ9Z3aqnCLTV
	//key_test.go:72: [FullPubKey] 02e5a7b2968c6c717e9a780b47606a56cf07c8d507182e33407641314b49dd5c1c
	//key_test.go:73: [WIF] cRWEXKcHFFAnkyN8He86uavPgPy44YJm6TdiUrA2SS6g9cF3K8qD
}

func TestKeyIntegrationBCH(t *testing.T) {
	//t.SkipNow()

	//BCHのconfとして利用する
	//FIXME: 一旦コメントアウトしたので、bch側から実行できるようにする
	//wlt.BTC.OverrideChainParamsByBCH()

	testSeed := "ggqMLyaZ7pwOXRdH8N2CWBf7L9gS/P8/p7oJdjp9M8U="
	bSeed, err := SeedToByte(testSeed)
	if err != nil {
		t.Fatal(err)
	}

	generateKeys(bSeed, enum.BCH, t)
}

//For Debug
func TestCheckBTCPrivateKey(t *testing.T) {
	conf := wlt.BTC.GetChainConf()

	testSeed := "ggqMLyaZ7pwOXRdH8N2CWBf7L9gS/P8/p7oJdjp9M8U="
	bSeed, err := SeedToByte(testSeed)
	if err != nil {
		t.Fatal(err)
	}

	// BTC key生成
	keyData := NewKey(enum.BTC, conf)

	priv, _, err := keyData.CreateAccount(bSeed, enum.AccountTypeClient)
	if err != nil {
		t.Fatal(err)
	}

	child, err := keyData.GetExtendedKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	experimentalKey(child, t)
}

//For only Debug
func experimentalKey(child *hdkeychain.ExtendedKey, t *testing.T) {
	conf := wlt.BTC.GetChainConf()

	// Private Key
	privateKey, err := child.ECPrivKey()
	if err != nil {
		t.Fatal(err)
	}

	// WIF　(compress: true) => bitcoin coreでは圧縮したアドレスを表示する
	wif, err := btcutil.NewWIF(privateKey, conf, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("WIF String:  %s", wif.String())

	// Address(P2PKH)
	address, err := child.Address(conf)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("P2PKH Address String:  %s", address.String())

	// Address(P2PKH) BTC
	//keyData := key.NewKey(enum.BTC, conf)

	//serializedKey := privateKey.PubKey().SerializeCompressed()
	//pubKeyAddr, err := btcutil.NewAddressPubKey(serializedKey, conf)
	//log.Println("address.String()", address.String())       //mySBc7pWWXjBUmAtjBY3sCdgnPAvAzwCoA
	//log.Println("pubKeyAddr.String()", pubKeyAddr.String()) //022c70901aac621c4436c4cb1f2daa8b9a6ff2c9d707b3f2639319d902679e1dfd
	//log.Println("pubKeyAddr.AddressPubKeyHash().String()", pubKeyAddr.AddressPubKeyHash().String()) //mySBc7pWWXjBUmAtjBY3sCdgnPAvAzwCoA
	//log.Println("getFullPubKey(privateKey)", getFullPubKey(privateKey)) //pubKeyAddr.String()とは微妙に異なる。。
	//log.Println(" ")
}
