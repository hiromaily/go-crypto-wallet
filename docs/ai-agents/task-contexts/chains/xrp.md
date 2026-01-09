# XRP (Ripple) Reference

## Overview

| 項目 | 値 |
|------|-----|
| トランザクションモデル | アカウント型 |
| 通信方式 | gRPC (ripple-lib-server経由) |
| アドレス形式 | r... (Base58) |
| 特殊機能 | Destination Tag, アカウントリザーブ |

## Directory Structure

### Use Case Layer

```
internal/application/usecase/
├── keygen/xrp/
│   ├── generate_key.go          # 鍵生成
│   └── sign_transaction.go      # トランザクション署名
├── sign/xrp/
│   └── sign_transaction.go      # トランザクション署名
└── watch/xrp/
    ├── create_transaction.go    # トランザクション作成
    ├── monitor_transaction.go   # トランザクション監視
    └── send_transaction.go      # トランザクション送信
```

### Infrastructure Layer

```
internal/infrastructure/api/ripple/
├── connection.go                # gRPC接続
└── xrp/
    ├── ripple.go                # クライアント初期化
    ├── ripppleapi.go            # API実装
    ├── rippleapi_account.go     # アカウントAPI
    ├── rippleapi_address.go     # アドレスAPI
    ├── rippleapi_tx.go          # トランザクションAPI
    ├── public_account.go        # 公開アカウントAPI
    ├── public_ledger.go         # 台帳API
    ├── public_server_info.go    # サーバー情報
    ├── public_transaction.go    # 公開トランザクションAPI
    ├── admin_keygen.go          # 管理者鍵生成
    ├── balance.go               # 残高
    ├── amount.go                # 金額
    ├── transaction.go           # トランザクション
    ├── converter.go             # 型変換
    ├── errors.go                # エラー定義
    └── types.go                 # 型定義

# gRPC生成ファイル
    ├── account_grpc.pb.go
    ├── account.pb.go
    ├── address_grpc.pb.go
    ├── address.pb.go
    ├── transaction_grpc.pb.go
    └── transaction.pb.go
```

### CLI Layer

```
internal/interface-adapters/cli/watch/api/xrp/
├── api.go                       # API呼び出し
└── sendcoin.go                  # 送金コマンド
```

### Protocol Buffers

```
proto/rippleapi/
├── account.proto                # アカウントサービス定義
├── address.proto                # アドレスサービス定義
└── transaction.proto            # トランザクションサービス定義
```

## Key Concepts

### Account Model

XRPはアカウント型モデルを使用：

```
アカウント {
    Address: r...
    Balance: Drops単位
    Sequence: トランザクション番号
    Flags: アカウント設定フラグ
}

トランザクション {
    Account: 送信元アドレス
    Destination: 送信先アドレス
    Amount: Drops単位
    Fee: Drops単位
    Sequence: アカウントのSequence
    DestinationTag: オプション
}
```

### Drops/XRP Conversion

```
1 XRP = 1,000,000 Drops

// 変換例
0.1 XRP = 100,000 Drops
```

### Sequence Number

ETHのNonceに相当：

```go
// トランザクション送信前にSequenceを取得
sequence, err := client.GetAccountSequence(ctx, address)
```

### Account Reserve

XRPアカウントには最低残高要件があります：

```
基本リザーブ: 10 XRP (アカウント作成に必要)
オーナーリザーブ: 2 XRP × オブジェクト数
```

### Destination Tag

取引所等で受取人を識別するための数値タグ：

```go
type Payment struct {
    Destination    string
    DestinationTag *uint32  // オプション
    Amount         string
}
```

## gRPC Communication

### ripple-lib-server

XRPはripple-lib-serverを経由してgRPCで通信：

```
[Wallet] --gRPC--> [ripple-lib-server] --WebSocket--> [XRPL]
```

サーバー実装: `apps/ripple-lib-server/`

### Proto定義

```protobuf
// proto/rippleapi/transaction.proto
service TransactionService {
    rpc CreatePayment(CreatePaymentRequest) returns (CreatePaymentResponse);
    rpc SignTransaction(SignTransactionRequest) returns (SignTransactionResponse);
    rpc SubmitTransaction(SubmitTransactionRequest) returns (SubmitTransactionResponse);
}
```

### コード生成

```bash
# protobuf/gRPCコード生成
make proto-gen
```

## Implementation Patterns

### トランザクション作成 (Watch Wallet)

```go
// internal/application/usecase/watch/xrp/create_transaction.go
type CreateTransactionUseCase struct {
    xrpClient ports.XRPClient
}

func (u *CreateTransactionUseCase) Execute(ctx context.Context, req *CreateTxRequest) (*CreateTxResponse, error) {
    // 1. Sequence番号を取得
    // 2. 手数料を取得/設定
    // 3. トランザクションを構築
    // 4. 署名用データを返す
}
```

### トランザクション署名 (Keygen/Sign Wallet)

```go
// internal/application/usecase/keygen/xrp/sign_transaction.go
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
# config/wallet/xrp_watch.toml
coin_type = "xrp"
network_type = "mainnet"  # mainnet, testnet

[ripple]
# ripple-lib-server接続設定
host = "localhost"
port = 50051
```

## Testing

```bash
# XRP Use Caseテスト
go test ./internal/application/usecase/watch/xrp/...
go test ./internal/application/usecase/keygen/xrp/...
go test ./internal/application/usecase/sign/xrp/...

# XRP Infrastructureテスト
go test ./internal/infrastructure/api/ripple/xrp/...
```

## Related Documentation

- [XRP README](../../../../docs/crypto/xrp/README.md) - 詳細なXRPドキュメント
- [Code Generation](../../guidelines/code-generation.md) - protobuf生成

## Common Operations

| 操作 | Wallet | Use Case |
|------|--------|----------|
| アドレス生成 | Keygen | `keygen/xrp/generate_key.go` |
| トランザクション作成 | Watch | `watch/xrp/create_transaction.go` |
| 署名 | Keygen/Sign | `{keygen,sign}/xrp/sign_transaction.go` |
| 送信 | Watch | `watch/xrp/send_transaction.go` |
| 残高確認 | Watch | Infrastructure gRPC API |

## ripple-lib-server

### 起動

```bash
cd apps/ripple-lib-server
yarn install
yarn start
```

### 構成

```
apps/ripple-lib-server/
├── src/
│   ├── grpc/           # gRPCサービス実装
│   ├── ripple/         # XRPL接続
│   └── index.ts        # エントリーポイント
├── package.json
└── tsconfig.json
```

## ⚠️ Implementation Notes

1. **gRPC依存**: XRP操作にはripple-lib-serverの起動が必要
2. **Sequence管理**: トランザクション送信前に最新のSequenceを取得
3. **Drops変換**: 金額操作時はDrops単位で計算
4. **アカウントリザーブ**: 新規アカウントには10XRP以上の初期送金が必要
5. **Destination Tag**: 取引所への送金時は必須の場合が多い
6. **Proto更新**: `.proto` ファイルを変更した場合は `make proto-gen` でコード再生成

## Comparison with Other Chains

| 項目 | XRP | ETH | BTC |
|------|-----|-----|-----|
| モデル | アカウント | アカウント | UTXO |
| 通信 | gRPC | JSON-RPC | RPC |
| 順序管理 | Sequence | Nonce | UTXO |
| 手数料 | Drops (固定的) | Gas (変動) | sat/vB |
| 最低残高 | 10 XRP | なし | なし |

