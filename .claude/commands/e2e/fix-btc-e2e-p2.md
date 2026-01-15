# Fix BTC E2E Pattern 2 Test #{issue_number}

BTC E2Eテスト（パターン2: P2PKH 2-of-3 Multisig）を **regtest環境** で実行し、エラーを修正する。

## 前提条件

**以下の共通ルールを最初に読み込むこと：**

- @.claude/rules/btc-e2e-script.md - BTC E2E共通ルール（ビルド、検証、エスカレーション、セキュリティ）

## パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `{issue_number}` | Optional | GitHub issue番号。指定時はgit-workflowに従う |

## 概要

このコマンドは`scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh`を **Bitcoin Core regtest環境** で実行し、発生したエラーを分析・修正します。

> **Note**: このE2Eテストはローカルのregtest（Regression Test）環境で実行されます。
> 実際のBitcoinネットワーク（mainnet/testnet）には接続しません。

### Pattern 2 の技術仕様

| 項目 | 値 |
|------|-----|
| **パターン番号** | 2 |
| **ネットワーク** | **regtest** (ローカル環境) |
| **鍵タイプ** | P2PKH (BIP44 Legacy) |
| **スクリプトタイプ** | 2-of-3 Multisig (P2SH wrapped) |
| **アドレス形式** | `2...` (regtest/testnet P2SH) |
| **署名要件** | 2-of-3 (任意の2つの署名で完了) |
| **Descriptor** | `sh(multi(2, [fp/44'/0'/0']xpub1/0/*, xpub2/0/*, xpub3/0/*))` |
| **必要なウォレット** | watch, keygen, sign1, sign2 |
| **環境変数** | `WALLET_ADDRESS_TYPE="legacy"` |

### Pattern 1 (Single-sig) との違い

| 項目 | Pattern 1 | Pattern 2 |
|------|-----------|-----------|
| 署名要件 | Single-sig (1つ) | 2-of-3 Multisig |
| 必要ウォレット | keygen のみ | keygen + sign1 + sign2 |
| アドレス形式 | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey交換 | 不要 | 必要 |
| account設定 | `account_singlesig.yaml` | `account_2of3.yaml` |

### Pattern 8 (3-of-3) との違い

| 項目 | Pattern 2 | Pattern 8 |
|------|-----------|-----------|
| 鍵タイプ | BIP44 (Legacy) | BIP49 (P2SH-SegWit) |
| 署名要件 | 2-of-3 | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | `sh(wsh(sortedmulti(3,...)))` |
| 署名フロー | 2回署名で完了 | 3回署名必要 |

### issue番号が指定された場合

`git-workflow`スキルを読み込み、以下の設定で作業してください：

- **ブランチ名**: `fix/issue-{issue_number}-btc-e2e-p2`
- **コミットタイプ**: `fix(btc)`
- **スコープ**: BTC E2E Pattern 2

→ 詳細は @.claude/skills/git-workflow/SKILL.md を参照

### issue番号が指定されない場合

ブランチ作成・PR作成なしで、ローカルで修正のみ行います。

## Pattern 2 固有のドキュメント

共通ルールの Required Documentation に加えて、以下を参照：

- @scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh - 実行対象のスクリプト
- @scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh - Pattern 1 スクリプト（Single-sig部分の参考）
- @scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh - Pattern 8 スクリプト（Multisig部分の参考）
- @config/wallet/account_2of3.yaml - 2-of-3 マルチシグ設定

## 事前確認: 環境変数

**パターン2では `WALLET_ADDRESS_TYPE="legacy"` が必要です。**

スクリプト内で自動設定されますが、確認用：

```bash
echo $WALLET_ADDRESS_TYPE  # "legacy" であること
```

> **Note**: 設定ファイルを直接編集しないでください。環境変数で上書きします。
> 詳細は共通ルールの「Configuration File Policy」を参照。

## 実行手順

### Step 1: E2Eテストを実行

```bash
# フルリセットして実行（推奨）
make btc-e2e-p2pkh-2of3-reset

# または既存状態から実行
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh

# デバッグ出力付き
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --verbose
```

> **Note**: ビルドと検証コマンドは共通ルールを参照。

### Step 2: エラー分析

エラーが発生したフェーズを特定し、対応するコードを調査：

| Phase | 関連コード | 説明 |
|-------|-----------|------|
| Prerequisites | CLI commands | `watch`, `keygen`, `sign1`, `sign2` |
| Infrastructure | Docker/compose | `compose.btc.yaml`, `compose.yaml` |
| Wallet Setup | Bitcoin RPC | `createwallet`, `loadwallet` |
| Key Generation | HD Key derivation | `internal/application/usecase/keygen/` |
| Multisig Setup | Descriptor export/import | `internal/application/usecase/watch/` |
| UTXO Generation | Bitcoin Core RPC | `generatetoaddress`, `deriveaddresses` |
| Transaction Flow | PSBT signing | `internal/infrastructure/wallet/api/btc/` |

### Step 3: コードを修正

エラー種別に応じて適切なスキルをロード（共通ルールの Related Skills 参照）。

## 技術仕様: 署名フロー（2-of-3）

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature) ← ここで完了
    ↓
Watch Wallet (broadcast)

※ Sign2 は不要 - 2つの署名で2-of-3が満たされる
```

### 重要な技術ポイント

1. **2-of-3 vs 3-of-3**
   - Pattern 2: 任意の2つの署名で完了
   - Pattern 8: 全3つの署名が必要
   - 署名フローの制御が異なる

2. **Descriptor形式**
   - `sh(multi(2, ...))` - P2SH wrapper + 2-of-3 multisig
   - `sortedmulti` ではなく `multi` を使用可能
   - 鍵の順序に注意

3. **HD Key派生パス**
   - BIP44: `m/44'/0'/account'/change/index`
   - Pattern 8 (BIP49) とは異なる派生パス

## Pattern 2 固有のエラー

共通エラー（No utxo, RPC接続等）は共通ルールを参照。以下はPattern 2固有のエラー：

### Descriptor形式エラー

**症状**: Descriptor export/import 時にエラー

**原因**: BIP49 形式になっている（P2SH-SegWit）

**解決策**: `address_type` が `legacy` であることを確認

```bash
echo $WALLET_ADDRESS_TYPE  # "legacy" であること
```

### アドレス形式が異なる

**症状**: `m...` や `n...` で始まるアドレスが生成される（Single-sig P2PKHアドレス）

**原因**: Single-sig アドレスが生成されている

**解決策**:

- `account_2of3.yaml` の `multisig` セクションを確認
- fullpubkey のインポートが成功しているか確認
- `required: 2` が設定されているか確認

### 署名が足りない / 署名が多すぎる

**症状**: トランザクション送信時に「署名が不完全」または「署名が多すぎる」エラー

**原因**:

- 2回目の署名が正しく適用されていない
- 3回目の署名を行っている（不要）

**確認**:

```bash
# PSBTの署名状態を確認
btc_cli "btc-watch" analyzepsbt "${psbt_hex}"
```

### fullpubkey インポートエラー

**症状**: Multisig セットアップ時にエラー

**原因**: fullpubkey の形式不一致または順序の問題

**解決策**:

1. `sign1`, `sign2` から正しくエクスポートされているか確認
2. `keygen` へのインポート順序を確認
3. 関連コード: `internal/infrastructure/wallet/key/fullpubkey/`

## 関連コード（Go）

| パス | 役割 |
|------|------|
| `internal/application/usecase/keygen/btc/` | 鍵生成ユースケース |
| `internal/application/usecase/watch/btc/` | Watch walletユースケース |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC実装 |
| `internal/infrastructure/wallet/key/fullpubkey/` | fullpubkey処理 |
| `internal/domain/wallet/key/` | 鍵ドメインモデル |

## 注意事項

### 他パターンへの影響を避ける

- パターン2固有の修正は`P2PKH 2-of-3`関連コードに限定
- 共通コードを修正する場合は、他パターン（特に1, 8）への影響を確認
- 共通関数を修正する場合は単体テストで回帰を確認

> **Note**: ビルドルール、セキュリティは共通ルールを参照。

## クリーンアップ

```bash
# コンテナ停止のみ
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --cleanup

# 完全リセット（データ含む）
./scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh --reset
```
