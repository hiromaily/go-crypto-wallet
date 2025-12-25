# internal/wallet ディレクトリのリファクタリング提案

## 現状分析

### 現在の構造

```
internal/
├── wallet/                    # 独立したディレクトリ
│   ├── keygener.go           # Keygener インターフェース
│   ├── signer.go             # Signer インターフェース
│   ├── watcher.go            # Watcher インターフェース
│   ├── btcwallet/            # BTC 実装
│   ├── ethwallet/            # ETH 実装
│   └── xrpwallet/            # XRP 実装
└── interface-adapters/
    └── cli/                   # CLI アダプター
```

### 問題点

1. **アーキテクチャ上の位置づけが不明確**
   - `internal/wallet/` は実質的にアダプター層として機能しているが、独立したディレクトリとして存在
   - Clean Architecture の観点から、`interface-adapters` 層に属すべき

2. **依存関係の複雑さ**
   - アプリケーション層（use case）とインフラ層（API クライアント）の両方に依存
   - コイン固有の実装が混在

3. **型アサーションの多用**
   - CLI からコイン固有の機能にアクセスするために型アサーションを使用
   - 例：`v.BTC`, `v.ETH` など

## リファクタリング提案

### 推奨案: `internal/interface-adapters/wallet/` に移動

#### 理由

1. **アーキテクチャの明確化**
   - `internal/wallet/` は CLI と use case の間のアダプター層として機能
   - `interface-adapters` 層に移動することで、役割が明確になる

2. **一貫性の向上**
   - 既に `internal/interface-adapters/cli/` が存在
   - 同じ層内で関連するアダプターをまとめる

3. **依存関係の整理**
   - アダプター層として明確に位置づけられる
   - 将来的な拡張（HTTP アダプターなど）も同じ層に配置可能

#### 新しい構造

```
internal/
└── interface-adapters/
    ├── cli/                   # CLI コマンド実装
    ├── wallet/                # Wallet アダプター（移動後）
    │   ├── interfaces.go      # Keygener, Signer, Watcher インターフェース
    │   ├── btc/               # BTC 実装
    │   │   ├── keygen.go
    │   │   ├── sign.go
    │   │   └── watch.go
    │   ├── eth/               # ETH 実装
    │   │   ├── keygen.go
    │   │   ├── sign.go
    │   │   └── watch.go
    │   └── xrp/               # XRP 実装
    │       ├── keygen.go
    │       ├── sign.go
    │       └── watch.go
    └── http/                  # HTTP ハンドラー（既存）
```

## 実装手順

### Phase 1: 準備

1. **新しいディレクトリ構造の作成**
   ```bash
   mkdir -p internal/interface-adapters/wallet/{btc,eth,xrp}
   ```

2. **インターフェースファイルの作成**
   - `internal/interface-adapters/wallet/interfaces.go` を作成
   - `keygener.go`, `signer.go`, `watcher.go` の内容を統合

3. **コイン別実装の移動**
   - `internal/wallet/btcwallet/` → `internal/interface-adapters/wallet/btc/`
   - `internal/wallet/ethwallet/` → `internal/interface-adapters/wallet/eth/`
   - `internal/wallet/xrpwallet/` → `internal/interface-adapters/wallet/xrp/`

### Phase 2: コード修正

1. **パッケージ名の変更**
   ```go
   // Before
   package btcwallet
   
   // After
   package btc
   ```

2. **インポートパスの更新**
   - すべての参照箇所でインポートパスを更新
   - `internal/wallet` → `internal/interface-adapters/wallet`

3. **型名の調整（オプション）**
   ```go
   // Before
   type BTCKeygen struct { ... }
   
   // After (より明確な命名)
   type Keygen struct { ... }
   // または
   type BTCKeygen struct { ... } // パッケージ名が btc なので、BTC プレフィックスは不要かも
   ```

### Phase 3: 型アサーションの改善（オプション）

現在、CLI からコイン固有の機能にアクセスするために型アサーションを使用：

```go
// 現在の実装
if v, ok := wallet.(*btcwallet.BTCKeygen); ok {
    v.BTC.GetChainConf()
}
```

改善案：

1. **インターフェース拡張**
   ```go
   type Keygener interface {
       // 既存のメソッド
       GenerateSeed() ([]byte, error)
       // ...
       
       // コイン固有の機能へのアクセス
       GetChainConfig() *chaincfg.Params
       GetAPI() interface{} // または具体的な型
   }
   ```

2. **ヘルパー関数の導入**
   ```go
   // internal/interface-adapters/wallet/btc/helpers.go
   func GetBitcoiner(keygen Keygener) (bitcoin.Bitcoiner, error) {
       if btcKeygen, ok := keygen.(*btc.Keygen); ok {
           return btcKeygen.BTC, nil
       }
       return nil, fmt.Errorf("not a BTC keygen")
   }
   ```

### Phase 4: 後方互換性の維持（移行期間）

移行期間中は、`internal/wallet` に型エイリアスを残す：

```go
// internal/wallet/keygener.go (後方互換性)
package wallet

import (
    walletadapter "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// Type aliases for backward compatibility
type Keygener = walletadapter.Keygener
type Signer = walletadapter.Signer
type Watcher = walletadapter.Watcher
```

### Phase 5: テストと検証

1. **すべてのテストを実行**
   ```bash
   make gotest
   ```

2. **ビルド確認**
   ```bash
   make check-build
   ```

3. **リント確認**
   ```bash
   make lint-fix
   ```

4. **インポートパスの確認**
   ```bash
   grep -r "internal/wallet" --include="*.go" .
   ```

### Phase 6: クリーンアップ

1. **後方互換性の型エイリアスを削除**
   - すべての参照が更新されたことを確認後、`internal/wallet/` を削除

2. **ドキュメントの更新**
   - `AGENTS.md` の更新
   - `docs/directory_structure.md` の更新

## 代替案

### 案2: アプリケーション層のファサードとして配置

```
internal/
└── application/
    └── facade/              # ファサード層
        └── wallet/          # Wallet ファサード
```

**メリット:**
- 複数の use case を統合する役割が明確

**デメリット:**
- インフラ層（API クライアント）への依存がアプリケーション層に混入
- Clean Architecture の原則に反する可能性

### 案3: 現状維持 + ドキュメント改善

**メリット:**
- 変更が最小限
- リスクが低い

**デメリット:**
- アーキテクチャの不明確さが残る
- 将来的な拡張時に問題になる可能性

## 推奨事項

**推奨案1（`internal/interface-adapters/wallet/` への移動）を推奨**

理由：
1. Clean Architecture の原則に最も適合
2. アーキテクチャの位置づけが明確になる
3. 将来的な拡張（HTTP アダプターなど）も同じ層に配置可能
4. 既存の `interface-adapters` 層との一貫性が保たれる

## 移行チェックリスト

- [ ] Phase 1: 新しいディレクトリ構造の作成
- [ ] Phase 2: コード修正（パッケージ名、インポートパス）
- [ ] Phase 3: 型アサーションの改善（オプション）
- [ ] Phase 4: 後方互換性の型エイリアス追加
- [ ] Phase 5: テストと検証
  - [ ] すべてのテストがパス
  - [ ] ビルドが成功
  - [ ] リントエラーなし
- [ ] Phase 6: クリーンアップ
  - [ ] 後方互換性の型エイリアス削除
  - [ ] ドキュメント更新

## 注意事項

1. **段階的な移行**
   - 一度にすべてを変更せず、段階的に移行
   - 各フェーズでテストを実行

2. **後方互換性**
   - 移行期間中は型エイリアスで後方互換性を維持
   - すべての参照が更新されたことを確認してから削除

3. **型アサーション**
   - 型アサーションの使用を減らす方向で改善を検討
   - ただし、完全に排除する必要はない（実用的な判断）

4. **テストカバレッジ**
   - 移行前後でテストカバレッジを確認
   - 不足しているテストを追加

