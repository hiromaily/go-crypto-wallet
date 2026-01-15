# Fix BTC E2E Pattern 3 Test #{issue_number}

BTC E2Eテスト（パターン3: P2SH-P2WPKH Single-sig）を **regtest環境** で実装・修正する。

## 前提条件

**以下の共通ルールを最初に読み込むこと：**

- @.claude/rules/btc/e2e-script.md - BTC E2E共通ルール（ビルド、検証、エスカレーション、セキュリティ）

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは `scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh` を作成・実行し、発生したエラーを分析・修正します。

> **Note**: このE2Eテストはローカルのregtest（Regression Test）環境で実行されます。
> 実際のBitcoinネットワーク（mainnet/testnet）には接続しません。

### Pattern 3 の技術仕様

| 項目 | 値 |
|------|-----|
| **パターン番号** | 3 |
| **ネットワーク** | **regtest** (ローカル環境) |
| **鍵タイプ** | P2SH-P2WPKH (BIP49 Nested SegWit) |
| **スクリプトタイプ** | Single-sig |
| **アドレス形式** | `3...` (Mainnet), `2...` (regtest/testnet) |
| **署名要件** | Single-sig (1つの署名) |
| **Descriptor** | `sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))` |
| **必要なウォレット** | watch, keygen |
| **環境変数** | `WALLET_ADDRESS_TYPE="p2sh-segwit"` |

### Pattern 1 (P2PKH Single-sig) との違い

| 項目 | Pattern 1 | Pattern 3 |
|------|-----------|-----------|
| 鍵タイプ | BIP44 (Legacy) | BIP49 (Nested SegWit) |
| アドレス形式 | `m.../n...` (P2PKH) | `2...` (P2SH) |
| Descriptor | `pkh(...)` | `sh(wpkh(...))` |
| 環境変数 | `legacy` | `p2sh-segwit` |
| トランザクションサイズ | 大きい | 小さい (SegWit割引) |

### Pattern 8 (P2SH-P2WSH 3-of-3) との違い

| 項目 | Pattern 3 | Pattern 8 |
|------|-----------|-----------|
| 署名要件 | Single-sig | 3-of-3 Multisig |
| Descriptor | `sh(wpkh(...))` | `sh(wsh(sortedmulti(3,...)))` |
| 必要ウォレット | keygen のみ | keygen + sign1 + sign2 |
| fullpubkey交換 | 不要 | 必要 |

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `fix/issue-{issue_number}-btc-e2e-p3`
- **コミットタイプ**: `feat(btc)` (新規スクリプト作成の場合) / `fix(btc)` (修正の場合)
- **スコープ**: BTC E2E Pattern 3

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで実装・修正のみ行います。

## Pattern 3 固有のドキュメント

共通ルールの Required Documentation に加えて、以下を参照：

- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - Pattern 1 スクリプト（Single-sig のベース）
- @scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh - Pattern 8 スクリプト（P2SH-SegWit の参考）
- @config/wallet/account.yaml - Single-sig account設定

## 事前確認: 環境変数

**パターン3では `WALLET_ADDRESS_TYPE="p2sh-segwit"` が必要です。**

スクリプト内で自動設定されますが、確認用：

```bash
echo $WALLET_ADDRESS_TYPE  # "p2sh-segwit" であること
```

> **Note**: 設定ファイルを直接編集しないでください。環境変数で上書きします。
> 詳細は共通ルールの「Configuration File Policy」を参照。

## 実装手順

### Step 1: スクリプト作成

Pattern 1 (`e2e-p1-p2pkh-singlesig.sh`) をベースに、以下を変更：

1. ファイル名: `e2e-p3-p2sh-p2wpkh-singlesig.sh`
2. 環境変数: `WALLET_ADDRESS_TYPE="p2sh-segwit"`
3. ヘッダーコメント: Pattern 3 の仕様に更新
4. アドレス検証ロジック: `2...` 形式の確認

### Step 2: Makefile ターゲット追加

`make/btc_e2e.mk` に以下を追加：

```makefile
###############################################################################
# E2E Testing - Pattern 3: P2SH-P2WPKH Single-sig
###############################################################################
.PHONY: btc-e2e-p3-reset
btc-e2e-p3-reset:
 ./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --reset

.PHONY: btc-e2e-p3
btc-e2e-p3:
 ./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh

.PHONY: btc-e2e-p3-verbose
btc-e2e-p3-verbose:
 ./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --verbose

.PHONY: btc-e2e-p3-ci
btc-e2e-p3-ci:
 ./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --non-interactive

.PHONY: btc-e2e-p3-cleanup
btc-e2e-p3-cleanup:
 ./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --cleanup
```

### Step 3: E2Eテストを実行

```bash
# フルリセットして実行（推奨）
make btc-e2e-p3-reset

# デバッグ出力付き
./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --verbose
```

> **Note**: ビルドと検証コマンドは共通ルールを参照。

### Step 4: エラー分析

エラーが発生したフェーズを特定し、対応するコードを調査：

| Phase | 関連コード | 説明 |
|-------|-----------|------|
| Prerequisites | CLI commands | `watch`, `keygen` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | HD Key derivation | `internal/application/usecase/keygen/` |
| Descriptor Export | BIP49 Descriptor | `internal/infrastructure/wallet/key/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

## 技術仕様: P2SH-P2WPKH (Nested SegWit)

### アドレス構造

```
P2SH-P2WPKH アドレス:
┌─────────────────────────────────────────────────────────────┐
│  P2SH wrapper                                                │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  P2WPKH (Native SegWit)                                  │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Public Key Hash                                     │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Descriptor 形式

```
sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))
           └─ BIP49 derivation path
```

### 署名フロー（Single-sig）

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (sign with single key)
    ↓
Watch Wallet (broadcast)
```

## Pattern 3 固有のエラー

共通エラー（No utxo, RPC接続等）は共通ルールを参照。以下はPattern 3固有のエラー：

### address_type 不一致

**症状**: `m...` や `n...` で始まるアドレスが生成される（P2PKHアドレス）

**原因**: `address_type` が `legacy` になっている

**解決策**:

```bash
# 環境変数確認
echo $WALLET_ADDRESS_TYPE  # "p2sh-segwit" であること

# スクリプト内の設定確認
grep "WALLET_ADDRESS_TYPE" scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh
```

### Descriptor 形式エラー

**症状**: Descriptor export/import 時にエラー

**原因**: BIP44 形式（`pkh(...)`）になっている

**確認**:

```bash
# Descriptor ファイル確認
cat data/descriptor/btc/payment_descriptors.json

# 期待される形式
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
# → "sh(wpkh([...]xpub.../0/*))"
```

### key_type 自動派生の確認

**確認**: `address_type` から `key_type` が正しく派生されているか

| address_type | 期待される key_type |
|--------------|-------------------|
| `p2sh-segwit` | `bip49` |

関連コード: `internal/domain/address/types.go` の `AddrType.ToKeyType()`

### トランザクション署名エラー

**症状**: 署名時に witness 関連のエラー

**原因**: SegWit トランザクションの witness データ処理の問題

**確認**:

```bash
# PSBT の分析
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  analyzepsbt "${psbt_hex}"
```

## デバッグ用コマンド

### 状態確認

```bash
# Bitcoin ノード状態
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# ウォレット残高
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO 一覧
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# アドレス情報確認（P2SH-P2WPKH かどうか）
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
# → "isscript": true, "iswitness": false, "script": "witness_v0_keyhash"
```

### Descriptor 確認

```bash
# Keygen の Descriptor 一覧
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true

# Watch の Descriptor 一覧
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  listdescriptors
```

## 関連コード（Go）

| パス | 役割 |
|------|------|
| `internal/application/usecase/keygen/btc/` | 鍵生成ユースケース |
| `internal/application/usecase/watch/btc/` | Watch walletユースケース |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC実装 |
| `internal/infrastructure/wallet/key/descriptor/` | Descriptor処理 |
| `internal/domain/address/types.go` | address_type → key_type 変換 |
| `pkg/config/loader.go` | 設定ローダー |

## ドキュメント更新

スクリプト作成後、以下のドキュメントを更新：

1. `scripts/operation/btc/e2e/README.md` - スクリプト一覧に追加
2. `docs/crypto/btc/e2e_transaction_patterns.md` - 実装ステータス更新
3. `.claude/rules/btc/e2e-script.md` - パターン一覧に追加

## 注意事項

### 他パターンへの影響を避ける

- パターン3固有の修正は `P2SH-P2WPKH Single-sig` 関連コードに限定
- 共通コードを修正する場合は、他パターン（特に1, 2, 8）への影響を確認
- 共通関数を修正する場合は単体テストで回帰を確認

> **Note**: ビルドルール、セキュリティは共通ルールを参照。

## クリーンアップ

```bash
# コンテナ停止のみ
./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh --cleanup

# 完全リセット（データ含む）
make btc-e2e-p3-reset
```
