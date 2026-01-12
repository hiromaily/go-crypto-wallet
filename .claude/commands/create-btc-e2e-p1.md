# Create BTC E2E Pattern 1 Script #{issue_number}

BTC E2Eテスト（パターン1: P2PKH Single-sig）スクリプトを新規作成する。

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは `scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh` を新規作成します。
Pattern 8 (P2SH-P2WSH 3-of-3) の既存スクリプトを参考に、Single-sig 用にシンプル化します。

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `feat/issue-{issue_number}-btc-e2e-p1`
- **コミットタイプ**: `feat(btc)`
- **スコープ**: BTC E2E Pattern 1

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで開発のみ行います。

## 必須ドキュメント（最初に読み込む）

以下のドキュメントを最初に読み込んでください：

- @docs/crypto/btc/e2e_transaction_patterns.md - **全パターンの詳細（必須）**
- @scripts/operation/btc/README.md - E2Eスクリプトの使い方
- @scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh - **参考実装（Pattern 8）**
- @scripts/operation/common.sh - 共通ユーティリティ

## GitHub Issue

このコマンドに対応する Issue: **#325**

Issue URL: <https://github.com/hiromaily/go-crypto-wallet/issues/325>

## 技術仕様: パターン1 (P2PKH Single-sig)

### アドレス形式

```
Address Type: P2PKH (BIP44 Legacy)
Address Format: 1... (Mainnet), m.../n... (Testnet/Regtest)
Descriptor: pkh([fingerprint/44'/0'/account']xpub.../0/*)
```

### 鍵派生パス

```
BIP44: m/44'/0'/account'/change/index
```

### 署名フロー（Single-sig）

```
┌─────────────────────────────────────────────────────────┐
│                  SINGLE-SIG FLOW                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign with single key                │
│          ↓                                              │
│  3. Watch Wallet: Broadcast transaction                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**重要**: Sign1, Sign2 ウォレットは不要です。

### Config 設定

パターン1を実行するには、**Watch と Keygen の設定で `address_type: "legacy"` が必要**です。

| ファイル | 設定値 |
|----------|--------|
| `config/wallet/btc_watch.yaml` | `address_type: "legacy"` |
| `config/wallet/btc_keygen.yaml` | `address_type: "legacy"` |

**注意**: Sign1, Sign2 の設定は Single-sig では使用しません。

```yaml
# 正しい設定例（パターン1用）
address_type: "legacy"
```

## ディレクトリ構造

### 変更後の構造

```
scripts/operation/btc/
├── e2e/                          # NEW: E2E テスト専用ディレクトリ
│   ├── e2e-p2sh-p2wsh-3of3.sh   # MOVED: 既存の Pattern 8 スクリプト
│   └── e2e-p2pkh-singlesig.sh   # NEW: Pattern 1 スクリプト（これを作成）
├── create-bitcoind-wallet.sh
├── ...
└── README.md
```

## Acceptance Criteria

### 1. ディレクトリ構造の整理

- [ ] `scripts/operation/btc/e2e/` ディレクトリを作成
- [ ] 既存の `e2e-p2sh-p2wsh-3of3.sh` を `scripts/operation/btc/e2e/` に移動
- [ ] 関連ドキュメントのパス参照を更新:
  - `make/btc.mk` (5箇所)
  - `scripts/operation/btc/README.md` (6箇所)
  - `docs/crypto/btc/e2e_transaction_patterns.md` (3箇所)
  - 移動したスクリプト内の Usage コメント

### 2. 新規スクリプト作成

- [ ] `scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh` を作成
  - Single-sig ワークフローを実装（Keygen のみで署名完了）
  - 既存の `common.sh` ユーティリティを再利用
  - P2PKH (BIP44) アドレスタイプを使用
  - `--verbose`, `--cleanup`, `--reset`, `--non-interactive` オプションをサポート

### 3. Makefile ターゲット

- [ ] `make/btc.mk` の既存ターゲットのパスを更新
- [ ] 新規ターゲットを追加:
  - `btc-e2e-p2pkh` - 基本実行
  - `btc-e2e-p2pkh-reset` - 完全リセットからの実行
  - `btc-e2e-p2pkh-verbose` - 詳細出力
  - `btc-e2e-p2pkh-ci` - CI/CD 用非対話モード
  - `btc-e2e-p2pkh-cleanup` - クリーンアップ

### 4. ドキュメント更新

- [ ] `docs/crypto/btc/e2e_transaction_patterns.md` を更新
  - Pattern 1 のステータスを `✅ e2e-p2pkh-singlesig.sh` に変更

## 実装手順

### Step 1: ディレクトリ作成とファイル移動

```bash
# e2e ディレクトリを作成
mkdir -p scripts/operation/btc/e2e

# 既存スクリプトを移動
git mv scripts/operation/btc/e2e-p2sh-p2wsh-3of3.sh scripts/operation/btc/e2e/
```

### Step 2: パス参照を更新

以下のファイルのパス参照を更新：

1. **make/btc.mk** - 5箇所

   ```makefile
   # Before
   ./scripts/operation/btc/e2e-p2sh-p2wsh-3of3.sh
   # After
   ./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh
   ```

2. **scripts/operation/btc/README.md** - 6箇所

3. **docs/crypto/btc/e2e_transaction_patterns.md** - 3箇所

4. **scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh** - Usage コメント

### Step 3: 新規スクリプト作成

`e2e-p2pkh-singlesig.sh` の骨格：

```bash
#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 1: P2PKH Single-sig
# This script automates the complete Bitcoin workflow for single-sig P2PKH transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/crypto/btc/e2e_transaction_patterns.md - E2E transaction patterns
#
# Transaction Pattern:
#   Pattern 1: BTC P2PKH Single-sig
#   - Address Type: P2PKH (BIP44 Legacy)
#   - Address Format: `1...` (Mainnet), `m.../n...` (Testnet/Regtest)
#   - Signature Requirement: Single-sig (Keygen only)
#   - Descriptor: pkh([fingerprint/44'/0'/0']xpub.../0/*)
#
# Required Config Settings:
#   - config/wallet/btc_watch.yaml:  address_type: "legacy"
#   - config/wallet/btc_keygen.yaml: address_type: "legacy"

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Source common utilities
# shellcheck source=../../common.sh
source "${SCRIPT_DIR}/../../common.sh"

# Configuration
COIN="btc"
ENCRYPTED="false"
SIGN_WALLET_NUM=0  # Single-sig: no additional sign wallets needed
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# ... 以下、Pattern 8 を参考にシンプル化して実装
```

### Step 4: Single-sig 固有の実装ポイント

Pattern 8 (3-of-3 Multisig) との主な違い：

| 項目 | Pattern 8 (Multisig) | Pattern 1 (Single-sig) |
|------|---------------------|------------------------|
| 署名ウォレット数 | 3 (keygen, sign1, sign2) | 1 (keygen only) |
| Descriptor形式 | `sh(wsh(sortedmulti(...)))` | `pkh(...)` |
| address_type | `"p2sh-segwit"` | `"legacy"` |
| fullpubkey交換 | 必要 | 不要 |
| PSBT署名回数 | 3回 | 1回 |

### Step 5: Makefile ターゲット追加

`make/btc.mk` に追加：

```makefile
###############################################################################
# E2E Testing - Pattern 1: P2PKH Single-sig
###############################################################################
# Run Bitcoin E2E workflow Pattern 1 from completely fresh state (recommended)
.PHONY: btc-e2e-p2pkh-reset
btc-e2e-p2pkh-reset:
 ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --reset

# Run complete Bitcoin end-to-end workflow Pattern 1
.PHONY: btc-e2e-p2pkh
btc-e2e-p2pkh:
 ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh

# Run Bitcoin E2E workflow Pattern 1 with verbose output
.PHONY: btc-e2e-p2pkh-verbose
btc-e2e-p2pkh-verbose:
 ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --verbose

# Run Bitcoin E2E workflow Pattern 1 in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-p2pkh-ci
btc-e2e-p2pkh-ci:
 ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --non-interactive

# Cleanup Bitcoin E2E test environment Pattern 1
.PHONY: btc-e2e-p2pkh-cleanup
btc-e2e-p2pkh-cleanup:
 ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --cleanup
```

### Step 6: 検証

```bash
# 1. シェルスクリプトの lint
make shfmt

# 2. E2E テスト実行
make btc-e2e-p2pkh-reset
```

## 試行回数の上限とエスカレーション

**実装→テストのサイクルが5回を超えた場合、進捗を整理して報告してください。**

### エスカレーション条件

以下のいずれかに該当する場合、5回未満でも報告を検討：

- 同じエラーが繰り返し発生する
- Bitcoin仕様の深い理解が必要と判断される
- Go コード側の修正が必要になった場合

### 進捗報告フォーマット

```markdown
## 進捗報告

### 完了した項目
- [ ] ディレクトリ構造の整理
- [ ] パス参照の更新
- [ ] 新規スクリプト作成
- [ ] Makefile ターゲット追加
- [ ] ドキュメント更新

### 現在の状態
[現在の状態を説明]

### 発生している問題
[問題があれば記載]

### 次のステップ
[次に必要なアクション]
```

## 必須スキル

1. `git-workflow` - ブランチ管理とコミット
2. `shell-scripts` - シェルスクリプトの作成
3. `makefile-update` - Makefile の修正
4. `go-development` - Go コードの修正が必要な場合

## よくあるエラーと解決策

### "No utxo" エラー

**原因**: UTXOが見つからない、またはアドレスがsolvableでない

**解決策**:

1. Descriptorが正しくインポートされているか確認
2. `getaddressinfo`でsolvable/ismine状態を確認
3. ブロック生成(101ブロック)が完了しているか確認

```bash
# デバッグ: アドレス情報確認
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
  getaddressinfo "<address>"
```

### RPC接続エラー

**原因**: Dockerコンテナが起動していない

**解決策**:

```bash
# コンテナ状態確認
docker compose -f compose.btc.yaml ps

# ログ確認
docker compose -f compose.btc.yaml logs btc-watch
```

### 署名エラー

**原因**: 秘密鍵がインポートされていない、またはアドレスタイプの不一致

**解決策**:

1. `address_type: "legacy"` が設定されているか確認
2. Keygen ウォレットで秘密鍵がインポートされているか確認

## 関連ドキュメント

| ドキュメント | 内容 |
|-------------|------|
| `docs/crypto/btc/e2e_transaction_patterns.md` | 全11パターンの説明（**必須参照**） |
| `docs/crypto/btc/descriptor_examples.md` | Descriptorの例 |
| `docs/crypto/btc/psbt_developer_guide.md` | PSBT開発ガイド |
| `docs/crypto/btc/README.md` | BTC技術リファレンス |
| `scripts/operation/btc/README.md` | E2Eスクリプトの詳細 |

## 関連コード（Go）

| パス | 役割 |
|------|------|
| `internal/application/usecase/keygen/btc/` | 鍵生成ユースケース |
| `internal/application/usecase/watch/btc/` | Watch walletユースケース |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC実装 |
| `internal/domain/wallet/key/` | 鍵ドメインモデル |

## 注意事項

### ビルドルール

**Go コード修正後は必ずバイナリを再ビルドしてください。**

```bash
make build-all
```

### 既存スクリプトへの影響を避ける

- Pattern 8 (`e2e-p2sh-p2wsh-3of3.sh`) の動作を壊さないこと
- 共通コードを修正する場合は、他パターンへの影響を確認
- `common.sh` を拡張する場合は後方互換性を維持

### セキュリティ

- 秘密鍵をログに出力しない
- テスト用のパスフレーズ/RPCクレデンシャルは本番で使用しない
- `docs/standards/security.md`を参照

## クリーンアップ

```bash
# コンテナ停止のみ
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --cleanup

# 完全リセット（データ含む）
./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh --reset
```
