# ETH (Ethereum) Reference

## Overview

| 項目 | 値 |
|------|-----|
| トランザクションモデル | アカウント型 |
| 通信方式 | JSON-RPC |
| アドレス形式 | 0x... (40文字の16進数) |
| 特殊機能 | ERC-20トークン、スマートコントラクト |

## Directory Structure

### Use Case Layer

```
internal/application/usecase/
├── keygen/eth/
│   ├── import_private_key.go    # 秘密鍵インポート
│   └── sign_transaction.go      # トランザクション署名
├── sign/eth/
│   └── sign_transaction.go      # トランザクション署名
└── watch/eth/
    ├── create_transaction.go    # トランザクション作成
    ├── monitor_transaction.go   # トランザクション監視
    └── send_transaction.go      # トランザクション送信
```

### Infrastructure Layer

```
internal/infrastructure/api/ethereum/
├── connection.go              # RPC接続
├── api-interface.go           # インターフェース定義
├── eth/
│   ├── ethereum.go            # クライアント初期化
│   ├── client.go              # HTTPクライアント
│   ├── balance.go             # 残高取得
│   ├── key.go                 # 鍵管理
│   ├── transaction.go         # トランザクション
│   ├── converters.go          # 型変換
│   ├── rpc_eth.go             # eth_* RPC
│   ├── rpc_eth_gas.go         # Gas関連RPC
│   ├── rpc_eth_tx.go          # トランザクションRPC
│   ├── rpc_admin.go           # admin_* RPC
│   ├── rpc_net.go             # net_* RPC
│   ├── rpc_personal.go        # personal_* RPC
│   ├── rpc_miner.go           # miner_* RPC
│   └── rpc_web3.go            # web3_* RPC
├── erc20/
│   └── erc20.go               # ERC-20トークン操作
└── ethtx/
    ├── ethtx.go               # トランザクション構築
    └── converters.go          # 型変換
```

### CLI Layer

```
internal/interface-adapters/cli/
├── keygen/api/eth/
│   ├── api.go                 # API呼び出し
│   └── importrawkey.go        # 秘密鍵インポート
└── watch/api/eth/
    ├── api.go                 # API呼び出し
    ├── clientversion.go       # クライアントバージョン
    ├── netversion.go          # ネットワークバージョン
    ├── nodeinfo.go            # ノード情報
    └── syncing.go             # 同期状態
```

## Key Concepts

### Account Model

ETHはアカウント型トランザクションモデルを使用：

```
アカウント {
    Address: 0x...
    Balance: Wei単位
    Nonce: トランザクション数
    Code: コントラクトの場合
    Storage: コントラクトの場合
}

トランザクション {
    From: 送信元アドレス
    To: 送信先アドレス
    Value: 送金額 (Wei)
    Nonce: アカウントのNonce
    GasPrice: Gas単価
    GasLimit: 最大Gas量
    Data: コントラクト呼び出しデータ
}
```

### Nonce Management

Nonceはトランザクションの順序を保証：

```go
// ❌ Bad: Nonceを固定値で使用
nonce := uint64(0)

// ✅ Good: 最新のNonceを取得
nonce, err := client.PendingNonceAt(ctx, address)
```

**注意点**:
- Nonceは連続している必要がある
- 同じNonceで新しいトランザクションを送ると上書き（置換）
- ペンディングトランザクションがある場合は考慮が必要

### Gas Management

```go
// Gas Price取得
gasPrice, err := client.SuggestGasPrice(ctx)

// Gas Limit推定
gasLimit, err := client.EstimateGas(ctx, callMsg)

// 手数料 = GasPrice × GasUsed
// 最大手数料 = GasPrice × GasLimit
```

### Wei/Ether Conversion

```
1 Ether = 10^18 Wei
1 Gwei = 10^9 Wei

// 変換例
0.1 ETH = 100,000,000,000,000,000 Wei (10^17)
```

## ERC-20 Tokens

### Token Operations

```go
// internal/infrastructure/api/ethereum/erc20/erc20.go
type ERC20Client struct {
    // ...
}

// 残高取得
func (c *ERC20Client) BalanceOf(ctx context.Context, tokenAddr, accountAddr common.Address) (*big.Int, error)

// 転送
func (c *ERC20Client) Transfer(ctx context.Context, tokenAddr, to common.Address, amount *big.Int) (*types.Transaction, error)
```

### Contract ABI

```
contracts/token.abi  # ERC-20 ABI定義
```

## Implementation Patterns

### トランザクション作成 (Watch Wallet)

```go
// internal/application/usecase/watch/eth/create_transaction.go
type CreateTransactionUseCase struct {
    ethClient ports.EthereumClient
}

func (u *CreateTransactionUseCase) Execute(ctx context.Context, req *CreateTxRequest) (*CreateTxResponse, error) {
    // 1. Nonceを取得
    // 2. GasPriceを取得/推定
    // 3. GasLimitを推定
    // 4. トランザクションを構築
    // 5. 署名用データを返す
}
```

### トランザクション署名 (Keygen/Sign Wallet)

```go
// internal/application/usecase/keygen/eth/sign_transaction.go
type SignTransactionUseCase struct {
    keyStore ports.KeyStore
}

func (u *SignTransactionUseCase) Execute(ctx context.Context, req *SignTxRequest) (*SignTxResponse, error) {
    // 1. 秘密鍵を取得
    // 2. トランザクションに署名
    // 3. 署名済みトランザクションを返す
}
```

## Configuration

```toml
# config/wallet/eth_watch.toml
coin_type = "eth"
network_type = "mainnet"  # mainnet, goerli, sepolia

[ethereum]
host = "localhost"
port = 8545
# または Infura等のURL
url = "https://mainnet.infura.io/v3/YOUR_KEY"
```

## Testing

```bash
# ETH Use Caseテスト
go test ./internal/application/usecase/watch/eth/...
go test ./internal/application/usecase/keygen/eth/...
go test ./internal/application/usecase/sign/eth/...

# ETH Infrastructureテスト
go test ./internal/infrastructure/api/ethereum/eth/...
go test ./internal/infrastructure/api/ethereum/erc20/...
```

## Related Documentation

- [ETH README](../../../../docs/crypto/eth/README.md) - 詳細なETHドキュメント
- [ERC-20 Guide](../../../../docs/crypto/eth/ERC20.md) - ERC-20トークン操作

## Common Operations

| 操作 | Wallet | Use Case |
|------|--------|----------|
| アドレス生成 | Keygen | 秘密鍵から導出 |
| トランザクション作成 | Watch | `watch/eth/create_transaction.go` |
| 署名 | Keygen/Sign | `{keygen,sign}/eth/sign_transaction.go` |
| 送信 | Watch | `watch/eth/send_transaction.go` |
| 残高確認 | Watch | Infrastructure API |
| ERC-20転送 | Watch | `erc20/` API |

## JSON-RPC Methods

| カテゴリ | メソッド | 用途 |
|---------|---------|------|
| eth | eth_getBalance | 残高取得 |
| eth | eth_getTransactionCount | Nonce取得 |
| eth | eth_sendRawTransaction | トランザクション送信 |
| eth | eth_estimateGas | Gas推定 |
| eth | eth_gasPrice | GasPrice取得 |
| eth | eth_getTransactionReceipt | レシート取得 |
| net | net_version | ネットワークID |
| web3 | web3_clientVersion | クライアント情報 |

## ⚠️ Implementation Notes

1. **Nonce管理**: トランザクション送信時は必ず最新のNonceを取得
2. **Gas推定**: `eth_estimateGas` は失敗する可能性があるため、フォールバック値を用意
3. **Wei変換**: 金額操作時はWei単位で計算、表示時にEther変換
4. **EIP-1559**: 新しいガス方式（maxFeePerGas, maxPriorityFeePerGas）への対応を確認
5. **チェーンID**: 署名時に正しいチェーンIDを使用（リプレイ攻撃防止）

