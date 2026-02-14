# XRP (Ripple) Reference

> **⚠️ IMPORTANT: Architecture Change**
>
> **DEPRECATED**: `apps/xrpl-grpc-server/` and `proto/xrpapi/` are no longer used.
>
> **Current Approach**: Native Go implementation using xrpl-go libraries. All XRP functionality is now directly implemented in the Go codebase without gRPC dependencies.
>
> **Active Specification**: See `.kiro/specs/xrp-transaction-flow-alignment/` for current implementation details.

## Overview

| 項目 | 値 |
|------|-----|
| トランザクションモデル | アカウント型 |
| アドレス形式 | r... (Base58) |
| 特殊機能 | Destination Tag, アカウントリザーブ |
| 実装方式 | Native Go (xrpl-go) |

## Directory Structure

### Use Case Layer

```
internal/application/usecase/
├── keygen/xrp/
│   ├── generate_key.go          # 鍵生成
│   ├── generate_key_offline.go  # オフライン鍵生成
│   └── sign_transaction.go      # トランザクション署名
├── sign/xrp/
│   └── sign_transaction.go      # トランザクション署名
└── watch/xrp/
    ├── create_transaction.go    # トランザクション作成
    ├── create_multisig_tx.go    # マルチシグトランザクション作成
    ├── add_multisig_signature.go # マルチシグ署名追加
    ├── submit_multisig_tx.go    # マルチシグトランザクション送信
    ├── monitor_transaction.go   # トランザクション監視
    ├── send_transaction.go      # トランザクション送信
    ├── set_regular_key.go       # レギュラーキー設定
    └── set_signer_list.go       # 署名者リスト設定
```

### Infrastructure Layer

```
internal/infrastructure/api/xrp/
├── xrp.go                       # クライアント初期化
├── xrpapi.go                    # API実装
├── xrpapi_account.go            # アカウントAPI
├── xrpapi_account_test.go       # アカウントAPIテスト
├── xrpapi_address.go            # アドレスAPI
├── xrpapi_address_test.go       # アドレスAPIテスト
├── xrpapi_tx.go                 # トランザクションAPI
├── xrpapi_tx_test.go            # トランザクションAPIテスト
├── connection.go                # 接続管理
├── request_response.go          # リクエスト/レスポンス定義
├── public_account.go            # 公開アカウントAPI
├── public_account_test.go       # 公開アカウントAPIテスト
├── public_ledger.go             # 台帳API
├── public_path_orderbook.go     # パス/オーダーブックAPI
├── public_payment_channel.go    # ペイメントチャネルAPI
├── public_server_info.go        # サーバー情報
├── public_server_info_test.go   # サーバー情報テスト
├── public_subscription.go       # サブスクリプションAPI
├── public_transaction.go        # 公開トランザクションAPI
├── public_utility.go            # ユーティリティAPI
├── admin_keygen.go              # 管理者鍵生成
├── admin_keygen_test.go         # 管理者鍵生成テスト
├── admin_logging_data.go        # 管理者ログデータ
├── admin_peer.go                # 管理者ピア管理
├── admin_server_control.go      # 管理者サーバー制御
├── admin_status_debugging.go    # 管理者ステータス/デバッグ
├── balance.go                   # 残高
├── amount.go                    # 金額
├── transaction.go               # トランザクション
├── transaction_test.go          # トランザクションテスト
├── converter.go                 # 型変換
├── errors.go                    # エラー定義
├── types.go                     # 型定義
├── util.go                      # ユーティリティ
├── util_test.go                 # ユーティリティテスト
├── protogen/                    # [DEPRECATED] gRPC生成コード
│   ├── account.pb.go
│   ├── account_grpc.pb.go
│   ├── address.pb.go
│   ├── address_grpc.pb.go
│   ├── transaction.pb.go
│   └── transaction_grpc.pb.go
├── testutil/
│   └── xrp.go                  # テストユーティリティ
└── xrplgo/                      # Native Go実装 (xrpl-go)
    ├── client.go                # クライアント
    ├── client_test.go           # クライアントテスト
    ├── account.go               # アカウント操作
    ├── ledger.go                # 台帳操作
    └── transaction.go           # トランザクション操作
```

### CLI Layer

```
internal/interface-adapters/cli/watch/api/xrp/
├── api.go                       # API呼び出し
└── sendcoin.go                  # 送金コマンド
```

### Protocol Buffers [DEPRECATED]

```
proto/rippleapi/  [DEPRECATED - No longer used]
├── account.proto                # [Deprecated]
├── address.proto                # [Deprecated]
└── transaction.proto            # [Deprecated]
```

**Note**: Protocol buffers are no longer used. XRP operations now use native Go implementation with xrpl-go libraries.

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

## Native Go Communication

### Direct XRP Ledger Access

XRPは`xrpl-go`ライブラリを使用して直接通信：

```
[Wallet] --xrpl-go (WebSocket/HTTP)--> [XRPL]
```

実装: `internal/infrastructure/api/xrp/xrplgo/`

**使用ライブラリ**:

- **xrpscan/xrpl-go** (v0.2.11): WebSocket通信、トランザクション送信
- **Peersyst/xrpl-go** (予定): ネイティブGoトランザクション署名

### 旧gRPCアーキテクチャ [DEPRECATED]

```
[Wallet] --gRPC--> [xrpl-grpc-server] --WebSocket--> [XRPL]  [NO LONGER USED]
```

**Deprecated Components**:

- `apps/xrpl-grpc-server/` - TypeScript gRPCサーバー [廃止]
- `proto/rippleapi/*.proto` - Protocol Buffers定義 [廃止]
- `make proto`, `make proto-ts` - コード生成コマンド [不要]

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

```yaml
# config/wallet/xrp/watch.yaml
coin_type: xrp
network_type: mainnet  # mainnet, testnet

ripple:
  # XRP Ledger接続設定 (Native Go xrpl-go)
  ws_url: wss://s.altnet.rippletest.net:51233  # WebSocket URL
  # または
  http_url: https://s.altnet.rippletest.net:51234  # HTTP URL
```

## Testing

```bash
# XRP Use Caseテスト
go test ./internal/application/usecase/watch/xrp/...
go test ./internal/application/usecase/keygen/xrp/...
go test ./internal/application/usecase/sign/xrp/...

# XRP Infrastructureテスト
go test ./internal/infrastructure/api/xrp/xrp/...
```

## Related Documentation

- [XRP README](../../../../docs/crypto/xrp/README.md) - 詳細なXRPドキュメント
- [XRP Transaction Flow Alignment Spec](.kiro/specs/xrp-transaction-flow-alignment/) - 現在の実装仕様
- [Code Generation](../../guidelines/code-generation.md) - コード生成ガイド

## Common Operations

| 操作 | Wallet | Use Case |
|------|--------|----------|
| アドレス生成 | Keygen | `keygen/xrp/generate_key.go` |
| トランザクション作成 | Watch | `watch/xrp/create_transaction.go` |
| 署名 | Keygen/Sign | `{keygen,sign}/xrp/sign_transaction.go` |
| 送信 | Watch | `watch/xrp/send_transaction.go` |
| 残高確認 | Watch | Infrastructure gRPC API |

## ⚠️ Implementation Notes

1. **Sequence管理**: トランザクション送信前に最新のSequenceを取得
2. **Drops変換**: 金額操作時はDrops単位で計算
3. **アカウントリザーブ**: 新規アカウントには10XRP以上の初期送金が必要
4. **Destination Tag**: 取引所への送金時は必須の場合が多い
5. **Native Go実装**: gRPC依存を削除し、xrpl-goライブラリを直接使用

## Comparison with Other Chains

| 項目 | XRP | ETH | BTC |
|------|-----|-----|-----|
| モデル | アカウント | アカウント | UTXO |
| 通信 | WebSocket (xrpl-go) | JSON-RPC | RPC |
| 順序管理 | Sequence | Nonce | UTXO |
| 手数料 | Drops (固定的) | Gas (変動) | sat/vB |
| 最低残高 | 10 XRP | なし | なし |
