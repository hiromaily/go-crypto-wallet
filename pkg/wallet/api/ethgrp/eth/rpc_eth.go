package eth

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// ResponseSyncing response of eth_syncing
type ResponseSyncing struct {
	StartingBlock int64 `json:"startingBlock"`
	CurrentBlock  int64 `json:"currentBlock"`
	HighestBlock  int64 `json:"highestBlock"`
	KnownStates   int64 `json:"knownStates"`
	PulledStates  int64 `json:"pulledStates"`
}

// Syncing returns sync status or bool
//  - return false if not syncing (it means syncing is done)
//  - there seems 2 different responses
func (e *Ethereum) Syncing() (*ResponseSyncing, bool, error) {
	var (
		resIF  interface{}
		resMap map[string]string
	)

	err := e.rpcClient.CallContext(e.ctx, &resIF, "eth_syncing")
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call client.CallContext(eth_syncing)")
	}

	// try to cast to bool
	bRes, ok := resIF.(bool)
	if !ok {
		// interface can't not be casted to type map
		//resMap, ok = resIF.(map[string]string)
		err := e.rpcClient.CallContext(e.ctx, &resMap, "eth_syncing")
		if err != nil {
			return nil, false, errors.Wrap(err, "fail to call client.CallContext(eth_syncing)")
		}
		//grok.Value(resMap)
		//value map[string]string = [
		//	startingBlock string = "0x606c" 6
		//	currentBlock string = "0x95ac" 6
		//	highestBlock string = "0x294545" 8
		//	knownStates string = "0x2084c" 7
		//	pulledStates string = "0x1eb12" 7
		//]

		startingBlock, err := hexutil.DecodeBig(resMap["startingBlock"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		currentBlock, err := hexutil.DecodeBig(resMap["currentBlock"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		highestBlock, err := hexutil.DecodeBig(resMap["highestBlock"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		knownStates, err := hexutil.DecodeBig(resMap["knownStates"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		pulledStates, err := hexutil.DecodeBig(resMap["pulledStates"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}

		resSync := ResponseSyncing{
			StartingBlock: startingBlock.Int64(),
			CurrentBlock:  currentBlock.Int64(),
			HighestBlock:  highestBlock.Int64(),
			KnownStates:   knownStates.Int64(),
			PulledStates:  pulledStates.Int64(),
		}

		return &resSync, true, nil
	}
	return nil, bRes, nil
}

// ProtocolVersion returns the current ethereum protocol version
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_protocolversion
// - returns like 65
func (e *Ethereum) ProtocolVersion() (uint64, error) {

	var resProtocolVer string
	err := e.rpcClient.CallContext(e.ctx, &resProtocolVer, "eth_protocolVersion")
	if err != nil {
		return 0, errors.Wrapf(err, "fail to call rpc.CallContext(eth_protocolVersion) error: %s", err)
	}
	h, err := e.DecodeBig(resProtocolVer)
	if err != nil {
		return 0, err
	}

	return h.Uint64(), err
}

// Coinbase returns the client coinbase address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_coinbase
// Note:
//  - any accounts can not be retrieved if no private key in keystore in node server.
//  - when running geth command, keystore option should be used to specify directory
//  - that means, this rpc can be called from cold wallet which has private key
func (e *Ethereum) Coinbase() (string, error) {

	var resAddr string
	err := e.rpcClient.CallContext(e.ctx, &resAddr, "eth_coinbase")
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_coinbase)")
	}
	return resAddr, err
}

// Accounts returns a list of addresses owned by client
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_accounts
// https://github.com/ethereum/go-ethereum/wiki/Managing-your-accounts
// Note:
//  - any accounts can not be retrieved if no private key in keystore in node server.
//  - when running geth command, keystore option should be used to specify directory
//  - that means, this rpc can be called from cold wallet which has private key
// - private key is stored and read from node server
// - result is same as ListAccounts()
func (e *Ethereum) Accounts() ([]string, error) {

	var accounts []string
	err := e.rpcClient.CallContext(e.ctx, &accounts, "eth_accounts")
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_accounts)")
	}

	return accounts, nil
}

// BlockNumber returns the number of most recent block
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber
// - returns like `2706141`
// - result may be a bit changeable, may be better to call several times.
func (e *Ethereum) BlockNumber() (*big.Int, error) {

	var resBlockNumber string
	err := e.rpcClient.CallContext(e.ctx, &resBlockNumber, "eth_blockNumber")
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_blockNumber)")
	}
	h, err := hexutil.DecodeBig(resBlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// EnsureBlockNumber calls BlockNumber() several times
func (e *Ethereum) EnsureBlockNumber(loopCount int) (*big.Int, error) {
	latestBlockNumber := new(big.Int)
	for i := 0; i < loopCount; i++ {
		if i != 0 {
			time.Sleep(500 * time.Millisecond)
		}
		num, err := e.BlockNumber()
		if err != nil {
			return nil, err
		}
		if latestBlockNumber.Uint64() < num.Uint64() {
			latestBlockNumber = latestBlockNumber.SetUint64(num.Uint64())
		}
	}
	return latestBlockNumber, nil
}

// GetBalance returns the balance of the account of given address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance
// - `QuantityTagEarliest` must NOT be used
// - On goerli testnet, balance can be found just after sending coins
// - TODO: which quantityTag should be used `QuantityTagLatest` or `QuantityTagPending`
func (e *Ethereum) GetBalance(hexAddr string, quantityTag QuantityTag) (*big.Int, error) {

	var balance string
	err := e.rpcClient.CallContext(e.ctx, &balance, "eth_getBalance", hexAddr, quantityTag.String())
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call rpc.CallContext(eth_getBalance) quantityTag: %s", quantityTag.String())
	}
	h, err := hexutil.DecodeBig(balance)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// GetStoreageAt returns the value from a storage position at a given address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getstorageat
// - `QuantityTagEarliest` must NOT be used
// - always returns `0x0000000000000000000000000000000000000000000000000000000000000000`
// - how this function can be used??
//func (e *Ethereum) GetStoreageAt(hexAddr string, quantityTag QuantityTag) (string, error) {
//
//	var storagePosition string
//	err := e.rpcClient.CallContext(e.ctx, &storagePosition, "eth_getStorageAt", hexAddr, "0x0", quantityTag.String())
//	if err != nil {
//		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_getStorageAt)")
//	}
//
//	return storagePosition, nil
//}

// GetTransactionCount returns the number of transactions sent from an address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
// - this is used as nonce
// - `QuantityTagEarliest` must NOT be used
// - after sending coin from this address, result is counted??
// - generated new address is always 0
func (e *Ethereum) GetTransactionCount(hexAddr string, quantityTag QuantityTag) (*big.Int, error) {

	var transactionCount string
	err := e.rpcClient.CallContext(e.ctx, &transactionCount, "eth_getTransactionCount", hexAddr, quantityTag.String())
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getTransactionCount)")
	}
	h, err := hexutil.DecodeBig(transactionCount)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in a block from a block matching the given block hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbyhash
// - block hash can be found from https://www.etherchain.org/block/2706436 by block number
// - but how it is found for Goerli Testnet??
// FIXME: this RPC doesn't return anything
//func (e *Ethereum) GetBlockTransactionCountByBlockHash(blockHash string) (*big.Int, error) {
//
//	var txCount string
//	err := e.rpcClient.CallContext(e.ctx, &txCount, "eth_getBlockTransactionCountByHash", blockHash)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getBlockTransactionCountByHash)")
//	}
//	if txCount == "" {
//		e.logger.Debug("transactionCount is blank")
//		return nil, errors.New("transactionCount is blank")
//	}
//
//	h, err := hexutil.DecodeBig(txCount)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
//	}
//
//	return h, nil
//}

// GetBlockTransactionCountByNumber returns the number of transactions in a block matching the given block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber
// transaction count of block, it's possible transaction is 0
// - this number is fixed
func (e *Ethereum) GetBlockTransactionCountByNumber(blockNumber uint64) (*big.Int, error) {

	//convert uint64 to hex
	hexNum := hexutil.EncodeUint64(blockNumber)

	var txCount string
	err := e.rpcClient.CallContext(e.ctx, &txCount, "eth_getBlockTransactionCountByNumber", hexNum)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getBlockTransactionCountByNumber)")
	}
	if txCount == "" {
		e.logger.Debug("transactionCount is blank")
		return nil, nil
	}

	h, err := hexutil.DecodeBig(txCount)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// GetUncleCountByBlockHash eturns the number of uncles in a block from a block matching the given block hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblockhash
//func (e *Ethereum) GetUncleCountByBlockHash(blockHash string) (*big.Int, error) {
//
//	var uncleCount string
//	err := e.rpcClient.CallContext(e.ctx, &uncleCount, "eth_getUncleCountByBlockHash", blockHash)
//	if err != nil {
//		return nil, errors.Errorf("fail to call rpc.CallContext(eth_getUncleCountByBlockHash)")
//	}
//	if uncleCount == "" {
//		e.logger.Debug("uncleCount is blank")
//		return nil, errors.New("uncleCount is blank")
//	}
//
//	h, err := e.DecodeBig(uncleCount)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
//	}
//
//	return h, nil
//}

// GetUncleCountByBlockNumber returns the number of uncles in a block from a block matching the given block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblocknumber
// - this func would not be used anywhere
func (e *Ethereum) GetUncleCountByBlockNumber(blockNumber uint64) (*big.Int, error) {

	//convert int64 to hex
	blockHexNumber := hexutil.EncodeUint64(blockNumber)

	var uncleCount string
	err := e.rpcClient.CallContext(e.ctx, &uncleCount, "eth_getUncleCountByBlockNumber", blockHexNumber)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getUncleCountByBlockNumber)")
	}
	if uncleCount == "" {
		e.logger.Debug("uncleCount is blank")
		return nil, errors.New("uncleCount is blank")
	}

	h, err := e.DecodeBig(uncleCount)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// GetCode returns code at a given address ???
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getcode
// - always returns 0
//func (e *Ethereum) GetCode(hexAddr string, quantityTag QuantityTag) (*big.Int, error) {
//
//	var code string
//	err := e.rpcClient.CallContext(e.ctx, &code, "eth_getCode", hexAddr, quantityTag.String())
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getCode)")
//	}
//	e.logger.Debug("code", zap.String("code", code))
//	if code == "0x" {
//		code = "0x0"
//	}
//
//	h, err := hexutil.DecodeBig(code)
//	if err != nil {
//		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
//	}
//
//	return h, nil
//}

// eth_getBlockByHash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbyhash
// - no need, use GetBlockByNumber()

// BlockNumber response of eth_blockNumber
type BlockNumber struct {
	Number string
}

// BlockRawInfo is block raw info
type BlockRawInfo struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             string   `json:"size"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	Transactions     []string `json:"transactions"`
	Uncles           []string `json:"uncles"`
}

// BlockInfo is block info
type BlockInfo struct {
	Number           *big.Int `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            *big.Int `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       *big.Int `json:"difficulty"`
	TotalDifficulty  *big.Int `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             *big.Int `json:"size"`
	GasLimit         *big.Int `json:"gasLimit"`
	GasUsed          *big.Int `json:"gasUsed"`
	Timestamp        *big.Int `json:"timestamp"`
	Transactions     []string `json:"transactions"`
	Uncles           []string `json:"uncles"`
}

// GetBlockByNumber returns information about a block by block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func (e *Ethereum) GetBlockByNumber(blockNumber uint64) (*BlockInfo, error) {
	//convert int64 to hex
	blockHexNumber := hexutil.EncodeUint64(blockNumber)

	//var lastBlock BlockNumber
	var blockRawInfo BlockRawInfo

	err := e.rpcClient.CallContext(e.ctx, &blockRawInfo, "eth_getBlockByNumber", blockHexNumber, false)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getBlockByNumber)")
	}

	return convertBlockRawInfo(&blockRawInfo), nil
}

func convertBlockRawInfo(raw *BlockRawInfo) *BlockInfo {
	return &BlockInfo{
		Number:           decodeString(raw.Number),
		Hash:             raw.Hash,
		ParentHash:       raw.ParentHash,
		Nonce:            decodeString(raw.Nonce),
		Sha3Uncles:       raw.Sha3Uncles,
		LogsBloom:        raw.LogsBloom,
		TransactionsRoot: raw.TransactionsRoot,
		StateRoot:        raw.StateRoot,
		Miner:            raw.Miner,
		Difficulty:       decodeString(raw.Difficulty),
		TotalDifficulty:  decodeString(raw.TotalDifficulty),
		ExtraData:        raw.ExtraData,
		Size:             decodeString(raw.Size),
		GasLimit:         decodeString(raw.GasLimit),
		GasUsed:          decodeString(raw.GasUsed),
		Timestamp:        decodeString(raw.Timestamp),
		Transactions:     raw.Transactions,
		Uncles:           raw.Uncles,
	}
}

func decodeString(val string) *big.Int {
	decoded, err := hexutil.DecodeBig(val)
	if err != nil {
		return new(big.Int).SetUint64(0)
	}
	return decoded
}

// eth_getUncleByBlockHashAndIndex
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblockhashandindex

// eth_getUncleByBlockNumberAndIndex
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblocknumberandindex
