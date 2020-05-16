package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ResponseSyncing response of eth_syncing
type ResponseSyncing struct {
	StartingBlock int64 `json:"startingBlock"`
	CurrentBlock  int64 `json:"currentBlock"`
	HighestBlock  int64 `json:"highestBlock"`
	KnownStates   int64 `json:"knownStates"`
	PulledStates  int64 `json:"pulledStates"`
}

// BlockNumber response of eth_blockNumber
type BlockNumber struct {
	Number string
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
//  - any accounts can not be retrieved if no private key in keystore.
//  - when running geth command, keystore option should be used to specify directory
//  - that means, this rpc can be called from cold wallet
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
//  - any accounts can not be retrieved if no private key in keystore.
//  - when running geth command, keystore option should be used to specify directory
//  - that means, this rpc can be called from cold wallet
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

// GetBalance returns the balance of the account of given address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance
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
func (e *Ethereum) GetStoreageAt(hexAddr string, quantityTag QuantityTag) (string, error) {

	var storagePosition string
	err := e.rpcClient.CallContext(e.ctx, &storagePosition, "eth_getStorageAt", hexAddr, "0x0", quantityTag.String())
	if err != nil {
		return "", errors.Wrap(err, "fail to call rpc.CallContext(eth_getStorageAt)")
	}

	return storagePosition, nil
}

// GetTransactionCount returns the number of transactions sent from an address
// - this is used as nonce
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
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
func (e *Ethereum) GetBlockTransactionCountByHash(txHash string) (*big.Int, error) {

	var txCount string
	err := e.rpcClient.CallContext(e.ctx, &txCount, "eth_getBlockTransactionCountByHash", txHash)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getBlockTransactionCountByHash)")
	}
	if txCount == "" {
		e.logger.Debug("transactionCount is blank")
		return nil, errors.New("transactionCount is blank")
	}

	h, err := hexutil.DecodeBig(txCount)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions in a block matching the given block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber
func (e *Ethereum) GetBlockTransactionCountByNumber(number uint64) (*big.Int, error) {

	//convert uint64 to hex
	hexNum := hexutil.EncodeUint64(number)

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
func (e *Ethereum) GetUncleCountByBlockHash(blockHash string) (*big.Int, error) {

	var uncleCount string
	err := e.rpcClient.CallContext(e.ctx, &uncleCount, "eth_getUncleCountByBlockHash", blockHash)
	if err != nil {
		return nil, errors.Errorf("fail to call rpc.CallContext(eth_getUncleCountByBlockHash)")
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

// GetUncleCountByBlockNumber returns the number of uncles in a block from a block matching the given block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblocknumber
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
func (e *Ethereum) GetCode(hexAddr string, quantityTag QuantityTag) (*big.Int, error) {

	var code string
	err := e.rpcClient.CallContext(e.ctx, &code, "eth_getCode", hexAddr, quantityTag.String())
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getCode)")
	}
	e.logger.Debug("code", zap.String("code", code))
	if code == "0x" {
		code = "0x0"
	}

	h, err := hexutil.DecodeBig(code)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// eth_getBlockByHash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbyhash

// GetBlockByNumber returns information about a block by block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func (e *Ethereum) GetBlockByNumber(quantityTag QuantityTag) (*big.Int, error) {

	var lastBlock BlockNumber

	err := e.rpcClient.CallContext(e.ctx, &lastBlock, "eth_getBlockByNumber", quantityTag.String(), false)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_getBlockByNumber)")
	}

	h, err := hexutil.DecodeBig(lastBlock.Number)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// eth_getTransactionByBlockHashAndIndex
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblockhashandindex

// eth_getTransactionByBlockNumberAndIndex
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblocknumberandindex

// eth_pendingTransactions
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_pendingtransactions

// eth_getUncleByBlockHashAndIndex
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblockhashandindex

// eth_getUncleByBlockNumberAndIndex
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblocknumberandindex

//
