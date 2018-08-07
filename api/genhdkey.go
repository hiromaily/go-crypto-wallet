package api

import (
	"errors"
	"log"

	"encoding/base64"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
)

//TODO:WIP まだ検証段階
//ちゃんとドキュメント(BIP)を読まないと罠にはまる。。。
//BIP32, BIP44

// getSeed() is to return seed []byte
func getSeed() ([]byte, error) {
	log.Println("[Info] generate seed")
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return nil, err
	}
	base64seed := base64.StdEncoding.EncodeToString(seed)
	log.Println("[Debug] generated seed(string):", base64seed)
	log.Println(" ")

	return seed, nil
}

// generateMasterPrivateKey() is to generate private master key
func generateMasterPrivateKey(seed []byte, conf *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	log.Println("[Info] generate private master key")

	key, err := hdkeychain.NewMaster(seed, conf)
	if err != nil {
		return nil, err
	}

	// Debug
	if key.IsPrivate() {
		log.Println("[Debug] private key was generated")
	} else {
		log.Println("[Error] Unexpectedly, public key was generated")
		err = errors.New("Unexpectedly, public key was generated")
		return nil, err
	}

	log.Println("[Debug] Private Master Key String:", key.String())
	// Address
	keyHash, err := key.Address(conf)
	if err != nil {
		return nil, err
	}

	log.Println("[Debug] Private Master Key Address (??):", keyHash.EncodeAddress())
	// TODO: PrivateでもPublicでも、アドレスが変わらない。ひょっとしたら、private keyはアドレス化できないかも。
	log.Println(" ")

	//おまけ
	//stringからmaster keyを生成
	//key2, err2 := hdkeychain.NewKeyFromString(key.String())
	//if err2 != nil {
	//	err = err2
	//	log.Println(err)
	//	return
	//}

	return key, nil
}

// generateMasterPublicKey() is to generate public master key
func generateMasterPublicKey(masterKey *hdkeychain.ExtendedKey, conf *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	log.Println("[Info] generate public master key")

	key, err := masterKey.Neuter()
	if err != nil {
		return nil, err
	}

	// Debug
	if key.IsPrivate() {
		log.Println("[Error] Unexpectedly, private key was generated")
		err = errors.New("Unexpectedly, private key was generated")
		return nil, err
	}
	log.Println("[Debug] public key was generated")

	log.Println("[Debug] Public Master Key String:", key.String())
	// Address
	keyHash, err := key.Address(conf)
	if err != nil {
		return nil, err
	}

	log.Println("[Debug] Public Master Key Address:", keyHash.EncodeAddress())
	log.Println(" ")

	return key, nil
}

//TODO: generateHierarchicalChild()
func generateHierarchicalChild(masterKey *hdkeychain.ExtendedKey, layerIndex, indexFrom, indexTo uint32) error {
	log.Println("[Info] generate child master key")

	// m/0H
	//TODO:これは何にも使えない？？
	//HardenedKeyStart = 0x80000000
	//account0, err := masterKey.Child(hdkeychain.HardenedKeyStart + 0)
	account0, err := masterKey.Child(hdkeychain.HardenedKeyStart + layerIndex)
	if err != nil {
		return err
	}

	// Debug ここはヘッダ的な?? でも普通に使えそう
	if account0.IsPrivate() {
		log.Println("[Debug] private key was generated")
	} else {
		log.Println("[Error] Unexpectedly, public key was generated")
		err = errors.New("Unexpectedly, public key was generated")
		return err
	}
	log.Println("[Debug] Private Child Key String:", account0.String())
	log.Println(" ")

	// to public
	account0PubKey, err := account0.Neuter()
	if err != nil {
		return err
	}
	if account0PubKey.IsPrivate() {
		log.Println("[Error] Unexpectedly, private key was generated")
		err = errors.New("Unexpectedly, private key was generated")
		return err
	}
	log.Println("[Debug] public key was generated")

	log.Println("[Debug] Public Child Key String:", account0PubKey.String())
	log.Println(" ")

	// Generate Child
	//m/0H/0-10
	//未着手
	//var act0 *hdkeychain.ExtendedKey //for only debug
	//for i := indexFrom; i <= indexTo; i++ {
	//	act, err := account0.Child(i)
	//	if err != nil {
	//		log.Println(err)
	//		break
	//	}
	//	//debug
	//	if i==0{
	//		act0 = act
	//	}
	//	log.Printf("[Debug] Private Child(%d) Key String: %s\n", i, act.String())
	//}

	//TODO:こちらから生成したアドレスでも送金に使えるかチェックする必要がある

	return nil
}

// GenerateHDKey HDウォレットキー及びSeedを生成する
//TODO: WIP
//func (b *Bitcoin) GenerateHDKey(paramSeed string) (*btcutil.WIF, string, error) {
func (b *Bitcoin) GenerateHDKey(paramSeed string) error {
	var (
		seed []byte
		err  error
	)

	// 1.seed
	if paramSeed != "" {
		log.Println("[Info] use seed argument")
		//Decode
		unbase64, err := base64.StdEncoding.DecodeString(paramSeed)
		if err != nil {
			log.Fatalf("[Error] base64.StdEncoding.DecodeString(): error: %v", err)
		}
		log.Printf("[Debug] seed is %s from argument\n", unbase64)
		seed = unbase64
	} else {
		// Generate a random seed at the recommended length.
		seed, err = getSeed()
		if err != nil {
			return err
		}
	}

	// 2.Generate a new private master key using the seed.
	masterPrivateKey, err := generateMasterPrivateKey(seed, b.GetChainConf())
	if err != nil {
		return err
	}

	// 3.Generate a new public master key using the master private key.
	//masterPublicKey, err := generateMasterPublicKey(masterPrivateKey, *isTest)
	_, err = generateMasterPublicKey(masterPrivateKey, b.GetChainConf())
	if err != nil {
		return err
	}

	// 4.Childを作成
	generateHierarchicalChild(masterPrivateKey, 0, 0, 10)

	return nil
}
