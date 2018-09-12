package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/bookerzzz/grok"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

//TODO:参考に(中国語のサイト)
//https://www.haowuliaoa.com/article/info/11350.html

//transactionの命名ルールについて
//*wire.MsgTxなどは、txがsuffixとして扱われているため、それに習う。e.g. hexTx, hashTxなど

//入金時における、トランザクション作成順序
//[Online]  CreateRawTransaction
//[Offline] SignRawTransactionByHex

// SignRawTransactionWithWalletResult signrawtransactionwithwalletをcallしたresponseの型
type SignRawTransactionWithWalletResult struct {
	Hex      string                              `json:"hex"`
	Complete bool                                `json:"complete"`
	Errors   []SignRawTransactionWithWalletError `json:"errors"`
}

// SignRawTransactionWithWalletError SignRawTransactionWithWalletResult内のerrorオブジェクト
type SignRawTransactionWithWalletError struct {
	Txid      string `json:"txid"`
	Vout      int64  `json:"vout"`
	ScriptSig string `json:"scriptSig"`
	Sequence  int64  `json:"sequence"`
	Error     string `json:"error"`
}

// FundRawTransactionResult fundrawtransactionをcallしたresponseの型
type FundRawTransactionResult struct {
	Hex       string `json:"hex"`
	Fee       int64  `json:"fee"`
	Changepos int64  `json:"changepos"`
}

// ToHex 16進数のstringに変換する
func (b *Bitcoin) ToHex(tx *wire.MsgTx) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", errors.Errorf("tx.Serialize(): error: %s", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

//ToMsgTx 16進数のstringから、wire.MsgTxに変換する
func (b *Bitcoin) ToMsgTx(txHex string) (*wire.MsgTx, error) {
	byteHex, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Errorf("hex.DecodeString(): error: %s", err)
	}

	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(byteHex)); err != nil {
		return nil, err
	}

	return &msgTx, nil
}

// DecodeRawTransaction Hex stringをデコードして、Rawのトランザクションデータに変換する
func (b *Bitcoin) DecodeRawTransaction(hexTx string) (*btcjson.TxRawResult, error) {
	byteHex, err := hex.DecodeString(hexTx)
	if err != nil {
		return nil, errors.Errorf("hex.DecodeString(): error: %s", err)
	}
	resTx, err := b.client.DecodeRawTransaction(byteHex)
	if err != nil {
		return nil, errors.Errorf("client.DecodeRawTransaction(): error: %s", err)
	}

	return resTx, nil
}

// GetRawTransactionByHex Hexからトランザクションを取得する
func (b *Bitcoin) GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error) {

	hashTx, err := chainhash.NewHashFromStr(strHashTx)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %s", strHashTx, err)
	}

	tx, err := b.client.GetRawTransaction(hashTx)
	if err != nil {
		return nil, errors.Errorf("client.GetRawTransaction(hash): error: %s", err)
	}
	//MsgTx()
	//tx.MsgTx()

	return tx, nil
}

// GetTransactionByTxID txIDからトランザクション詳細を取得する
func (b *Bitcoin) GetTransactionByTxID(txID string) (*btcjson.GetTransactionResult, error) {
	// Transaction詳細を取得(必要な情報があるかどうか不明)
	hashTx, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %s", txID, err)
	}
	resTx, err := b.client.GetTransaction(hashTx)
	if err != nil {
		return nil, errors.Errorf("client.GetTransaction(%s): error: %s", hashTx, err)
	}
	//type GetTransactionResult struct {
	//	Amount          float64                       `json:"amount"`
	//	Fee             float64                       `json:"fee,omitempty"`
	//	Confirmations   int64                         `json:"confirmations"`
	//	BlockHash       string                        `json:"blockhash"`
	//	BlockIndex      int64                         `json:"blockindex"`
	//	BlockTime       int64                         `json:"blocktime"`
	//	TxID            string                        `json:"txid"`
	//	WalletConflicts []string                      `json:"walletconflicts"`
	//	Time            int64                         `json:"time"`
	//	TimeReceived    int64                         `json:"timereceived"`
	//	Details         []GetTransactionDetailsResult `json:"details"`
	//	Hex             string                        `json:"hex"`
	//}

	return resTx, nil
}

// GetTxOutByTxID TxOutを指定したトランザクションIDから取得する
func (b *Bitcoin) GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Errorf("chainhash.NewHashFromStr(%s): error: %s", txID, err)
	}

	// Gettxout / txHash *chainhash.Hash, index uint32, mempool bool
	txOutResult, err := b.client.GetTxOut(hash, index, false)
	if err != nil {
		return nil, errors.Errorf("client.GetTxOut(%s, %d, false): error: %s", hash, index, err)
	}

	return txOutResult, nil
	//logger.Infof("TxOut: %v\n", txOut): Output
	//grok.Value(txOut)
	//value *GetTxOutResult = {
	//	BestBlock string = "00000000000000a080461b99935872934fa35bc705453f9f2ad7712b2228e849" 64
	//	Confirmations int64 = 473
	//	Value float64 = 0.65
	//	ScriptPubKey ScriptPubKeyResult = {
	//		Asm string = "OP_HASH160 b24f4d8c8237c73a7299b6e790b309d477bb509c OP_EQUAL" 60
	//		Hex string = "a914b24f4d8c8237c73a7299b6e790b309d477bb509c87" 46
	//		ReqSigs int32 = 1
	//		Type string = "scripthash" 10
	//		Addresses []string = [
	//			0 string = "2N9W3GVam33jQc5FbkLKwMqH7RkvkYK7xvz" 35
	//		]
	//	}
	//	Coinbase bool = false
	//}
}

// CreateRawTransaction Rawトランザクションを作成する
//  => watch only wallet(online)で利用されることを想定
// こちらは多対1用の送金、つまり入金時、集約用アドレスに一括して送信するケースで利用することを想定
// [Noted] 手数料を考慮せず、全額送金しようとすると、SendRawTransaction()で、`min relay fee not met`
func (b *Bitcoin) CreateRawTransaction(sendAddr string, amount btcutil.Amount, inputs []btcjson.TransactionInput) (*wire.MsgTx, error) {
	//TODO:sendAddrの厳密なチェックがセキュリティ的に必要な場面もありそう
	sendAddrDecoded, err := btcutil.DecodeAddress(sendAddr, b.GetChainConf())
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress(%s): error: %s", sendAddr, err)
	}

	// パラメータを作成する
	outputs := make(map[btcutil.Address]btcutil.Amount)
	outputs[sendAddrDecoded] = amount //satoshi
	//lockTime := int64(0) //TODO:Raw locktime ここは何をいれるべき？

	// CreateRawTransaction
	//TODO:ここから下はCreateRawTransactionWithOutput()を呼び出すことでコードの重複を防いだほうがいい。
	return b.CreateRawTransactionWithOutput(inputs, outputs)
}

// CreateRawTransactionWithOutput 出金時に1対多で送信する場合に利用するトランザクションを作成する
func (b *Bitcoin) CreateRawTransactionWithOutput(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (*wire.MsgTx, error) {
	lockTime := int64(0) //TODO:Raw locktime ここは何をいれるべき？

	// CreateRawTransaction
	msgTx, err := b.client.CreateRawTransaction(inputs, outputs, &lockTime)
	if err != nil {
		return nil, errors.Errorf("btcutil.CreateRawTransaction(): error: %s", err)
	}

	return msgTx, nil
}

// FundRawTransaction 送信したい金額に応じて、自動的にutxoを算出してくれる
//  現時点で使う予定無し
func (b *Bitcoin) FundRawTransaction(hex string) (*FundRawTransactionResult, error) {
	//fundrawtransaction
	//https://bitcoincore.org/en/doc/0.16.2/rpc/rawtransactions/fundrawtransaction/
	//"{\"changePosition\":2}"

	//hex
	bHex, err := json.Marshal(hex)
	if err != nil {
		return nil, errors.Errorf("json.Marchal(hex): error: %s", err)
	}

	//fee rate
	feePerKb, err := b.EstimateSmartFee()
	if err != nil {
		return nil, errors.Errorf("BTC.EstimateSmartFee(): error: %s", err)
	}

	bFeeRate, err := json.Marshal(struct {
		FeeRate float64 `json:"feeRate"`
	}{
		FeeRate: feePerKb,
	})
	if err != nil {
		return nil, errors.Errorf("json.Marchal(feeRate): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex, bFeeRate})
	//rawResult, err := b.client.RawRequest("fundrawtransaction", []json.RawMessage{bHex})
	if err != nil {
		//error: -4: Insufficient funds
		return nil, errors.Errorf("json.RawRequest(fundrawtransaction): error: %s", err)
	}

	fundRawTransactionResult := FundRawTransactionResult{}
	err = json.Unmarshal([]byte(rawResult), &fundRawTransactionResult)
	if err != nil {
		return nil, errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	//logger.Infof("fundRawTransactionResult: %v: %s\n", fundRawTransactionResult, fundRawTransactionResult.Hex)

	return &fundRawTransactionResult, nil
}

// SignRawTransaction *wire.MsgTxからRawのトランザクションに署名する
// こちらは一連の流れを調査するために、CreateRawTransaction()の戻りに合わせたI/F (実際の運用で利用されることはないはず)
// 秘密鍵を保持している側のwallet(つまりcold wallet)で実行することを想定
func (b *Bitcoin) SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	//署名
	//FIXME:SignRawTransactionWithWallet()はエラーが出て使えないので、暫定対応としてoptionでdeprecatedを無効化する
	//restart bitcoind with -deprecatedrpc=signrawtransaction
	//FIXME:一旦コメントアウト
	//if b.Version() >= enum.BTCVer17 {
	//	msgTx, isSigned, err := b.SignRawTransactionWithWallet(tx)
	//	if err != nil {
	//		return nil, false, errors.Errorf("BTC.SignRawTransactionWithWallet() error: %s", err)
	//	}
	//	return msgTx, isSigned, nil
	//}
	return b.signRawTransaction(tx)
}

// signRawTransaction *wire.MsgTxからRawのトランザクションに署名する
// Deprecated
func (b *Bitcoin) signRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	//署名
	msgTx, isSigned, err := b.client.SignRawTransaction(tx)
	if err != nil {
		return nil, false, errors.Errorf("client.SignRawTransaction(): error: %s", err)
	}

	//Debug
	if !isSigned {
		logger.Debug("トランザクションHEXの結果比較スタート")
		////compare tx hex
		//b.debugCompareTx(tx, msgTx)

		//isSignedがfalseの場合、inputのtxとoutputのtxはまったく一緒だった
		//ここから、関連するprivate keyは出力できないか？pubkeyscriptから抽出できるかもしれない。
		//grok.Value(tx.TxOut)
		//for _, txout := range tx.TxOut {
		//	hexPubKey := hex.EncodeToString(txout.PkScript)
		//	logger.Debug("hexPubKey:", hexPubKey)
		//}
		//grok.Value(tx.TxIn)

		//WIF: cR5KECrkYF7RKi7v8DezyUCBdV1nYQF99gEapcCnzTRdXjWyiPgt
		//func (c *Client) SignRawTransaction3(tx *wire.MsgTx,
		//inputs []btcjson.RawTxInput,
		//	privKeysWIF []string) (*wire.MsgTx, bool, error) {

		//Debug
		msgTx, isSigned, err = b.client.SignRawTransaction3(tx, nil, []string{"cNw2Kd7fJLcyXJ9UvgEAUahwb6aoiwpXG8kSPELmC5Q42QemP6vq","cR5KECrkYF7RKi7v8DezyUCBdV1nYQF99gEapcCnzTRdXjWyiPgt"})
		if err != nil {
			return nil, false, errors.Errorf("client.SignRawTransaction3(): error: %s", err)
		}
		b.debugCompareTx(tx, msgTx)
	}

	//Multisigの場合、これによって署名が終了したか判断するはず
	//if !isSigned {
	//	return nil, errors.New("SignRawTransaction() can not sign on given transaction")
	//}

	return msgTx, isSigned, nil
}

func (b *Bitcoin) debugCompareTx(tx1, tx2 *wire.MsgTx){
	hexTx1, err := b.ToHex(tx1)
	if err != nil {
		logger.Debugf("btc.ToHex(tx1): error: %s", err)
	}
	hexTx2, err := b.ToHex(tx2)
	if err != nil {
		logger.Debugf("btc.ToHex(tx2): error: %s", err)
	}
	if hexTx1 == hexTx2 {
		logger.Debug("トランザクションHEXの結果が同じであった。これでは意味がない。")
	}
}

// SignRawTransactionWithWallet Ver17から利用可能なSignRawTransaction
// FIXME:まだエラーが出て使えない
// restart bitcoind with -deprecatedrpc=signrawtransaction
func (b *Bitcoin) SignRawTransactionWithWallet(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	//hex tx
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return nil, false, errors.Errorf("BTC.ToHex(tx): error: %s", err)
	}

	input1, err := json.Marshal(hexTx)
	if err != nil {
		return nil, false, errors.Errorf("json.Marchal(txHex): error: %s", err)
	}

	rawResult, err := b.client.RawRequest("signrawtransactionwithwallet", []json.RawMessage{input1})
	if err != nil {
		return nil, false, errors.Errorf("json.RawRequest(signrawtransactionwithwallet): error: %s", err)
	}
	//SignRawTransactionWithWalletResult
	signRawTxResult := SignRawTransactionWithWalletResult{}
	err = json.Unmarshal([]byte(rawResult), &signRawTxResult)
	if err != nil {
		return nil, false, errors.Errorf("json.Unmarshal(): error: %s", err)
	}
	if len(signRawTxResult.Errors) != 0 {
		//FIXME:error: Input not found or already spent
		//もし、解決のためにパラメータprevtxsが必要となると、非常にめんどくさい。。。listunspentで取得可能だが。。。
		//[debug]
		grok.Value(signRawTxResult)
		//value SignRawTransactionWithWalletResult = {
		//	Hex string = "020000000138c14167c131202d81d054bce8c6726c87ca072f14b03ec121755f5983b1c0b70000000000ffffffff0600093d..." 486
		//	Complete bool = false
		//	Errors []SignRawTransactionWithWalletError = [
		//		0 SignRawTransactionWithWalletError = {
		//			Txid string = "b7c0b183595f7521c13eb0142f07ca876c72c6e8bc54d0812d2031c16741c138" 64
		//			Vout int64 = 0
		//			ScriptSig string = "" 0
		//			Sequence int64 = 4294967295
		//			Error string = "Input not found or already spent" 32
		//		}
		//	]
		//}

		return nil, false, errors.Errorf("json.RawRequest(signrawtransactionwithwallet): error: %s", signRawTxResult.Errors[0].Error)
	}

	msgTx, err := b.ToMsgTx(signRawTxResult.Hex)
	if err != nil {
		return nil, false, errors.Errorf("BTC.ToMsgTx(hex): error: %s", err)
	}

	return msgTx, signRawTxResult.Complete, nil
}

// SignRawTransactionByHex HexからRawトランザクションを生成し、署名する
// 秘密鍵を保持している側のwallet(つまりcold wallet)で実行することを想定
func (b *Bitcoin) SignRawTransactionByHex(hex string) (string, bool, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return "", false, err
	}

	//署名
	signedTx, isSigned, err := b.SignRawTransaction(msgTx)
	if err != nil {
		return "", false, err
	}

	//Hexに変換
	hexTx, err := b.ToHex(signedTx)
	if err != nil {
		return "", false, errors.Errorf("BTC.ToHex(msgTx): error: %s", err)
	}

	//return signedTx, nil
	return hexTx, isSigned, nil
}

// SendRawTransaction Rawトランザクションを送信する
// こちらは一連の流れを調査するために、CreateRawTransaction()の戻りに合わせたI/F (実際の運用で利用されることはないはず)
// オンラインで実行される必要がある
func (b *Bitcoin) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	//送信
	hash, err := b.client.SendRawTransaction(tx, true)
	if err != nil {
		//feeを1Satoshiで試してみたら、
		//-26: 66: min relay fee not metが出た
		return nil, errors.Errorf("client.SendRawTransaction(): error: %s", err)
	}

	return hash, nil
}

// SendTransactionByHex 外部から渡されたバイト列からRawトランザクションを送信する
// オンラインで実行される必要がある
func (b *Bitcoin) SendTransactionByHex(hex string) (*chainhash.Hash, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, err
	}

	//送信
	hash, err := b.SendRawTransaction(msgTx)
	//hash, err := b.client.SendRawTransaction(msgTx, true)
	if err != nil {
		return nil, errors.Errorf("BTC.SendRawTransaction(): error: %s", err)
	}

	//txID
	//hash.String()

	return hash, nil
}

// SendTransactionByByte 外部から渡されたバイト列からRawトランザクションを送信する
// オンラインで実行される必要がある
func (b *Bitcoin) SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error) {

	//[]byte to wireTx
	wireTx := new(wire.MsgTx)
	r := bytes.NewBuffer(rawTx)

	if err := wireTx.Deserialize(r); err != nil {
		return nil, errors.Errorf("wireTx.Deserialize(): error: %s", err)
	}

	//送信
	hash, err := b.SendRawTransaction(wireTx)
	//hash, err := b.client.SendRawTransaction(wireTx, true)
	if err != nil {
		return nil, errors.Errorf("BTC.SendRawTransaction(): error: %v", err)
	}

	return hash, nil
}

// SequentialTransaction [Debug用]: 一連の未署名トランザクション作成から送信までの流れ
// TODO:Testに移動するか
func (b *Bitcoin) SequentialTransaction(hex string) (*chainhash.Hash, *btcutil.Tx, error) {
	// Hexからトランザクションを取得
	msgTx, err := b.ToMsgTx(hex)
	if err != nil {
		return nil, nil, err
	}

	//署名(オフライン)
	signedTx, isSigned, err := b.SignRawTransaction(msgTx)
	if err != nil {
		return nil, nil, err
	}
	if !isSigned {
		return nil, nil, errors.New("BTC.SignRawTransaction() can not sign on given transaction or multisig may be required")
	}

	//送金(オンライン)
	hash, err := b.SendRawTransaction(signedTx)
	if err != nil {
		return nil, nil, err
	}

	//txを取得
	resTx, err := b.GetRawTransactionByHex(hash.String())
	if err != nil {
		return nil, nil, err
	}

	return hash, resTx, nil
}

// Sign 署名を行う without Bitcoin Core [WIP]
//FIXME: これはColdWallet内で必要となるが、BitcoinCoreの機能が必要ないので、あれば実装しておきたい
func (b *Bitcoin) Sign(tx *wire.MsgTx, strPrivateKey string) (string, error) {
	// Key
	wif, err := btcutil.DecodeWIF(strPrivateKey)
	if err != nil {
		return "", err
	}
	privKey := wif.PrivKey

	// SignatureScript
	for idx, val := range tx.TxIn {
		//
		script, err := txscript.SignatureScript(tx, idx, val.SignatureScript, txscript.SigHashAll, privKey, false)
		if err != nil {
			return "", err
		}
		tx.TxIn[idx].SignatureScript = script
	}
	//TODO: isSignedかどうかをどうチェックするか
	//TODO: TxInごとに異なるKeyの場合は難しい

	//Hexに変換
	hexTx, err := b.ToHex(tx)
	if err != nil {
		return "", err
	}

	return hexTx, nil
}
