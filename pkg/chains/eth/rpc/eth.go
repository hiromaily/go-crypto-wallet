package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mitchellh/mapstructure"

	"github.com/hiromaily/go-crypto-wallet/pkg/debug"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// Syncing calls eth_syncing and returns the sync status.
//   - Returns (nil, false, nil) when the node is fully synced.
//   - Returns (*ResponseSyncing, true, nil) when the node is still syncing.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_syncing
func (c *rpcClient) Syncing(ctx context.Context) (*ResponseSyncing, bool, error) {
	var result any

	err := c.caller.CallContext(ctx, &result, "eth_syncing")
	if err != nil {
		return nil, false, fmt.Errorf("fail to call client.CallContext(eth_syncing): %w", err)
	}

	// try to cast to bool first
	if v, ok := result.(bool); ok {
		return nil, v, nil
	} else if v2, ok2 := result.(map[string]any); ok2 {
		anyMap := make(map[string]int64)
		for key, val := range v2 {
			if strVal, ok3 := val.(string); ok3 {
				var decoded *big.Int
				decoded, err = hexutil.DecodeBig(strVal)
				if err != nil {
					continue
				}
				anyMap[key] = decoded.Int64()
			}
		}

		resSyncing := ResponseSyncing{}
		if err = mapstructure.Decode(anyMap, &resSyncing); err != nil {
			return nil, false, err
		}
		debug.Debug(resSyncing)
		return &resSyncing, true, nil
	}

	return nil, false, errors.New("unexpected response, need to be implemented")
}

// ProtocolVersion returns the current Ethereum protocol version.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_protocolversion
func (c *rpcClient) ProtocolVersion(ctx context.Context) (uint64, error) {
	var resProtocolVer string
	err := c.caller.CallContext(ctx, &resProtocolVer, "eth_protocolVersion")
	if err != nil {
		return 0, fmt.Errorf("fail to call rpc.CallContext(eth_protocolVersion) error: %s: %w", err, err)
	}
	h, err := decodeBig(resProtocolVer)
	if err != nil {
		return 0, err
	}

	return h.Uint64(), err
}

// Coinbase returns the client coinbase address.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_coinbase
func (c *rpcClient) Coinbase(ctx context.Context) (string, error) {
	var resAddr string
	err := c.caller.CallContext(ctx, &resAddr, "eth_coinbase")
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(eth_coinbase): %w", err)
	}
	return resAddr, err
}

// Accounts returns a list of addresses owned by the client.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_accounts
func (c *rpcClient) Accounts(ctx context.Context) ([]string, error) {
	var accounts []string
	err := c.caller.CallContext(ctx, &accounts, "eth_accounts")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_accounts): %w", err)
	}

	return accounts, nil
}

// BlockNumber returns the number of the most recent block.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber
func (c *rpcClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	var resBlockNumber string
	err := c.caller.CallContext(ctx, &resBlockNumber, "eth_blockNumber")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_blockNumber): %w", err)
	}
	h, err := hexutil.DecodeBig(resBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// EnsureBlockNumber calls BlockNumber loopCount times (with 500 ms between calls) and
// returns the highest observed block number.
func (c *rpcClient) EnsureBlockNumber(ctx context.Context, loopCount int) (*big.Int, error) {
	latestBlockNumber := new(big.Int)
	for i := range loopCount {
		if i != 0 {
			time.Sleep(500 * time.Millisecond)
		}
		num, err := c.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}
		if latestBlockNumber.Uint64() < num.Uint64() {
			latestBlockNumber = latestBlockNumber.SetUint64(num.Uint64())
		}
	}
	return latestBlockNumber, nil
}

// GetBalance returns the balance of the given address.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance
func (c *rpcClient) GetBalance(ctx context.Context, hexAddr string, quantityTag QuantityTag) (*big.Int, error) {
	var balance string
	err := c.caller.CallContext(ctx, &balance, "eth_getBalance", hexAddr, quantityTag.String())
	if err != nil {
		return nil, fmt.Errorf(
			"fail to call rpc.CallContext(eth_getBalance) quantityTag: %s: %w",
			quantityTag.String(), err)
	}
	h, err := hexutil.DecodeBig(balance)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// GetTransactionCount returns the number of transactions sent from an address (nonce).
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
func (c *rpcClient) GetTransactionCount(
	ctx context.Context, hexAddr string, quantityTag QuantityTag,
) (*big.Int, error) {
	var transactionCount string
	err := c.caller.CallContext(
		ctx, &transactionCount, "eth_getTransactionCount", hexAddr, quantityTag.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_getTransactionCount): %w", err)
	}
	h, err := hexutil.DecodeBig(transactionCount)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions in a block by block number.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber
func (c *rpcClient) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64) (*big.Int, error) {
	hexNum := hexutil.EncodeUint64(blockNumber)

	var txCount string
	err := c.caller.CallContext(ctx, &txCount, "eth_getBlockTransactionCountByNumber", hexNum)
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_getBlockTransactionCountByNumber): %w", err)
	}
	if txCount == "" {
		logger.Debug("transactionCount is blank")
		return nil, nil
	}

	h, err := hexutil.DecodeBig(txCount)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// GetUncleCountByBlockNumber returns the number of uncles in a block by block number.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblocknumber
func (c *rpcClient) GetUncleCountByBlockNumber(ctx context.Context, blockNumber uint64) (*big.Int, error) {
	blockHexNumber := hexutil.EncodeUint64(blockNumber)

	var uncleCount string
	err := c.caller.CallContext(ctx, &uncleCount, "eth_getUncleCountByBlockNumber", blockHexNumber)
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_getUncleCountByBlockNumber): %w", err)
	}
	if uncleCount == "" {
		logger.Debug("uncleCount is blank")
		return nil, errors.New("uncleCount is blank")
	}

	h, err := decodeBig(uncleCount)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// GetBlockByNumber returns information about a block by block number.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func (c *rpcClient) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*BlockInfo, error) {
	blockHexNumber := hexutil.EncodeUint64(blockNumber)

	var blockRawInfo BlockRawInfo
	err := c.caller.CallContext(ctx, &blockRawInfo, "eth_getBlockByNumber", blockHexNumber, false)
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_getBlockByNumber): %w", err)
	}

	return convertBlockRawInfo(&blockRawInfo), nil
}
