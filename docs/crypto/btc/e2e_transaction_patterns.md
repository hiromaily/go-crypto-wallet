# E2E Transaction Patterns Guide

このドキュメントは、Bitcoin/Bitcoin Cash におけるトランザクションの組み合わせパターンについて説明します。鍵の種類とマルチシグの有無によって、様々なE2Eワークフローパターンが存在します。

## 目次

1. [概要](#概要)
2. [サポートされる鍵の種類](#サポートされる鍵の種類)
3. [署名パターン](#署名パターン)
4. [E2Eワークフローマトリックス](#e2eワークフローマトリックス)
5. [各パターンの詳細](#各パターンの詳細)
6. [アカウントタイプと署名要件](#アカウントタイプと署名要件)
7. [実装状態](#実装状態)
8. [E2Eスクリプトの対応表](#e2eスクリプトの対応表)

---

## 概要

Bitcoinトランザクションは、以下の2つの主要な軸で分類できます：

1. **鍵の種類（アドレスタイプ）** - どのBIPに基づいてアドレスを生成するか
2. **署名パターン** - シングルシグかマルチシグか

これらの組み合わせにより、様々なE2Eワークフローが必要となります。

---

## サポートされる鍵の種類

### Bitcoin (BTC)

| アドレスタイプ | BIP | Prefix (Mainnet) | Prefix (Testnet) | 説明 |
|---------------|-----|------------------|------------------|------|
| **P2PKH** (Legacy) | BIP44 | `1...` | `m.../n...` | 従来のPay-to-Public-Key-Hash |
| **P2SH-P2WPKH** | BIP49 | `3...` | `2...` | SegWit wrapped in P2SH |
| **P2WPKH** (Native SegWit) | BIP84 | `bc1q...` | `tb1q...` | Native SegWit |
| **P2TR** (Taproot) | BIP86 | `bc1p...` | `tb1p...` | Taproot (推奨) |

### Bitcoin Cash (BCH)

| アドレスタイプ | Prefix | 説明 |
|---------------|--------|------|
| **CashAddr** | `bitcoincash:q...` | Bitcoin Cash専用フォーマット |
| **Legacy** | `1...` | レガシーフォーマット（互換性用） |

### 鍵派生パス

| 標準 | パス | 用途 |
|------|------|------|
| BIP44 | `m/44'/0'/account'/change/index` | P2PKH (Legacy) |
| BIP49 | `m/49'/0'/account'/change/index` | P2SH-P2WPKH |
| BIP84 | `m/84'/0'/account'/change/index` | P2WPKH (Native SegWit) |
| BIP86 | `m/86'/0'/account'/change/index` | P2TR (Taproot) |

---

## 署名パターン

### Single-Sig（シングルシグ）

1つの秘密鍵で署名するパターン。

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

**特徴:**

- シンプルで高速
- 1回の署名で完了
- 秘密鍵が1つのため、紛失・漏洩リスクが集中

### Multi-Sig（マルチシグ）

複数の秘密鍵で署名するパターン。M-of-N（N個の鍵のうちM個の署名が必要）。

#### 3-of-3 マルチシグ

```
┌─────────────────────────────────────────────────────────┐
│                  3-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign (1st signature)                │
│          ↓                                              │
│  3. Sign1 Wallet: Sign (2nd signature)                 │
│          ↓                                              │
│  4. Sign2 Wallet: Sign (3rd signature)                 │
│          ↓                                              │
│  5. Watch Wallet: Broadcast transaction                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 2-of-3 マルチシグ

```
┌─────────────────────────────────────────────────────────┐
│                  2-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign (1st signature)                │
│          ↓                                              │
│  3. Sign1 Wallet: Sign (2nd signature)                 │
│          ↓                                              │
│  4. Watch Wallet: Broadcast transaction                 │
│                                                         │
│  (Sign2 Wallet は不要 - 2つの署名で完了)                │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### MuSig2（署名集約）

Schnorr署名ベースの集約署名プロトコル。N-of-N マルチシグがシングルシグと同じサイズになる。

```
┌─────────────────────────────────────────────────────────┐
│                    MUSIG2 FLOW                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Round 1: Nonce Generation (並列実行可能)              │
│  ├─ Keygen Wallet: Generate nonce                       │
│  ├─ Sign1 Wallet: Generate nonce                        │
│  └─ Sign2 Wallet: Generate nonce                        │
│          ↓                                              │
│  Round 2: Signing (順次実行)                           │
│  ├─ Keygen Wallet: Create partial signature             │
│  ├─ Sign1 Wallet: Create partial signature              │
│  └─ Sign2 Wallet: Create partial signature              │
│          ↓                                              │
│  Aggregation:                                           │
│  └─ Watch Wallet: Aggregate & broadcast                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**MuSig2の利点:**

- トランザクションサイズが30-50%削減
- プライバシー向上（シングルシグと見分けがつかない）
- 手数料削減

---

## E2Eワークフローマトリックス

### BTC パターンマトリックス

| パターン | 鍵タイプ | 署名パターン | アドレスフォーマット | E2Eスクリプト対応 |
|----------|---------|-------------|---------------------|-------------------|
| 1 | P2PKH (BIP44) | Single-sig | `1...` | 🔶 手動テスト |
| 2 | P2PKH (BIP44) | 2-of-3 Multisig | `3...` (P2SH wrapped) | ❌ 未対応 |
| 3 | P2SH-P2WPKH (BIP49) | Single-sig | `3...` | 🔶 手動テスト |
| 4 | P2SH-P2WPKH (BIP49) | 2-of-3 Multisig | `3...` | ❌ 未対応 |
| 5 | P2WPKH (BIP84) | Single-sig | `bc1q...` | 🔶 手動テスト |
| 6 | P2WSH (BIP84) | 2-of-3 Multisig | `bc1q...` | ❌ 未対応 |
| 7 | P2WSH (BIP84) | 3-of-3 Multisig | `bc1q...` | ❌ 未対応 |
| **8** | **P2SH-P2WSH** | **3-of-3 Multisig** | **`3...`** | **✅ e2e-workflow.sh** |
| 9 | P2TR (BIP86) | Single-sig | `bc1p...` | 🔶 手動テスト |
| 10 | P2TR (BIP86) | MuSig2 (N-of-N) | `bc1p...` | 🔜 開発中 |
| 11 | P2TR (BIP86) | Tapscript (M-of-N) | `bc1p...` | 🔜 開発中 |

### BCH パターンマトリックス

| パターン | 鍵タイプ | 署名パターン | アドレスフォーマット | E2Eスクリプト対応 |
|----------|---------|-------------|---------------------|-------------------|
| 1 | CashAddr | Single-sig | `bitcoincash:q...` | 🔶 手動テスト |
| **2** | **CashAddr** | **3-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e-workflow.sh** |
| 3 | CashAddr | 2-of-3 Multisig | `bitcoincash:p...` | ❌ 未対応 |

---

## 各パターンの詳細

### パターン 8: BTC P2SH-P2WSH 3-of-3 Multisig（現在のE2E）

**現在の `scripts/operation/btc/e2e-workflow.sh` で実装されているパターン**

```
アドレスタイプ: P2SH-P2WSH (BIP49 wrapped SegWit)
署名要件: 3-of-3 (Keygen + Sign1 + Sign2)
Descriptor: sh(wsh(sortedmulti(3, xpub1, xpub2, xpub3)))
```

**ワークフロー:**

1. Keygen/Sign1/Sign2 で Seed を生成
2. Keygen で HD Key を生成（各アカウント10個）
3. Sign1/Sign2 で HD Key を生成
4. Sign1/Sign2 から fullpubkey をエクスポート
5. Keygen に fullpubkey をインポート
6. Keygen で Descriptor をエクスポート
7. Watch に Descriptor をインポート
8. Test UTXO を生成（regtest）
9. 未署名トランザクション作成 → 3回署名 → ブロードキャスト

### パターン 2: BCH CashAddr 3-of-3 Multisig（現在のE2E）

**現在の `scripts/operation/bch/e2e-workflow.sh` で実装されているパターン**

```
アドレスタイプ: CashAddr P2SH
署名要件: 3-of-3 (Keygen + Sign1 + Sign2)
アドレス形式: bitcoincash:p... (P2SH multisig)
```

**ワークフロー:**

1. Keygen/Sign1/Sign2 で Seed を生成
2. Keygen で HD Key を生成
3. Sign1/Sign2 で HD Key を生成
4. Sign1/Sign2 から fullpubkey をエクスポート
5. Keygen に fullpubkey をインポート
6. Keygen で Multisig アドレスを作成
7. Keygen からアドレスをエクスポート
8. Watch にアドレスをインポート
9. Test UTXO を生成（regtest）
10. 未署名トランザクション作成 → 3回署名 → ブロードキャスト

### パターン 9: BTC P2TR Single-sig（Taproot）

```
アドレスタイプ: P2TR (BIP86)
署名要件: Single-sig (Keygen のみ)
Descriptor: tr([fingerprint/86'/0'/0']xpub.../0/*)
```

**シンプルなワークフロー:**

1. Keygen で Seed を生成
2. Keygen で BIP86 HD Key を生成
3. Keygen から Taproot アドレスをエクスポート
4. Watch に Taproot アドレスをインポート
5. 未署名トランザクション作成 → 1回署名（Schnorr）→ ブロードキャスト

### パターン 10: BTC P2TR MuSig2（開発中）

```
アドレスタイプ: P2TR (BIP86)
署名要件: N-of-N MuSig2 (全員の署名が必要)
Descriptor: tr(musig(xpub1, xpub2, xpub3))
```

**2ラウンドプロトコル:**

1. Round 1: 各ウォレットでノンスを生成
2. Round 2: 各ウォレットで部分署名を作成
3. Watch で署名を集約してブロードキャスト

---

## アカウントタイプと署名要件

| アカウント | 用途 | 推奨署名パターン | 理由 |
|-----------|------|-----------------|------|
| **client** | 顧客入金アドレス | Single-sig | 顧客側での操作が必要なため |
| **deposit** | 入金集約 | Multisig (2-of-3 または 3-of-3) | セキュリティ強化 |
| **payment** | 支払い | Multisig (2-of-3 または 3-of-3) | 承認フロー |
| **stored** | 長期保管 | Multisig (3-of-3) | 最高レベルのセキュリティ |

---

## 実装状態

### 鍵タイプの実装状態

| 鍵タイプ | BTC | BCH |
|---------|-----|-----|
| P2PKH (Legacy) | ✅ 実装済み | N/A |
| P2SH-P2WPKH (BIP49) | ✅ 実装済み | N/A |
| P2WPKH (BIP84) | ✅ 実装済み | N/A |
| P2TR (BIP86) | ✅ 実装済み | N/A |
| CashAddr | N/A | ✅ 実装済み |

### 署名パターンの実装状態

| 署名パターン | BTC | BCH |
|-------------|-----|-----|
| Single-sig | ✅ 実装済み | ✅ 実装済み |
| 2-of-3 Multisig | ⚠️ 部分的 | ⚠️ 部分的 |
| 3-of-3 Multisig | ✅ 実装済み | ✅ 実装済み |
| MuSig2 | 🔜 開発中 | N/A |

### Descriptor サポート

| 機能 | BTC | BCH |
|------|-----|-----|
| Descriptor Export | ✅ 実装済み | ❌ 未対応 |
| Descriptor Import | ✅ 実装済み | ❌ 未対応 |
| Bitcoin Core連携 | ✅ 実装済み | N/A |

---

## E2Eスクリプトの対応表

### 現在利用可能なE2Eスクリプト

| スクリプト | コイン | パターン | 署名要件 |
|-----------|--------|---------|---------|
| `scripts/operation/btc/e2e-workflow.sh` | BTC | P2SH-P2WSH Multisig | 3-of-3 |
| `scripts/operation/bch/e2e-workflow.sh` | BCH | CashAddr Multisig | 3-of-3 |

### 今後追加予定のE2Eスクリプト

| スクリプト（予定） | コイン | パターン | 署名要件 | 優先度 |
|-------------------|--------|---------|---------|--------|
| `e2e-singlesig.sh` | BTC | P2WPKH/P2TR Single-sig | 1 | 高 |
| `e2e-musig2.sh` | BTC | P2TR MuSig2 | N-of-N | 中 |
| `e2e-2of3.sh` | BTC | P2WSH 2-of-3 | 2-of-3 | 低 |
| `e2e-tapscript.sh` | BTC | P2TR Script Path | M-of-N | 低 |

---

## クイックリファレンス

### BTC アドレス種別の見分け方

| プレフィックス | 種類 | BIP | SegWit |
|---------------|------|-----|--------|
| `1...` | P2PKH | BIP44 | ❌ |
| `3...` | P2SH または P2SH-P2WPKH | BIP16/BIP49 | △ |
| `bc1q...` | P2WPKH または P2WSH | BIP84 | ✅ |
| `bc1p...` | P2TR (Taproot) | BIP86 | ✅ |

### BCH アドレス種別の見分け方

| プレフィックス | 種類 | Multisig |
|---------------|------|----------|
| `bitcoincash:q...` | P2PKH | ❌ |
| `bitcoincash:p...` | P2SH | ✅ |

### トランザクションサイズ比較

| パターン | Weight | vBytes | 備考 |
|----------|--------|--------|------|
| P2PKH Single-sig (1-in, 2-out) | ~680 | ~170 | Legacy |
| P2WPKH Single-sig (1-in, 2-out) | ~440 | ~110 | Native SegWit |
| P2TR Single-sig (1-in, 2-out) | ~396 | ~99 | Taproot |
| 2-of-3 P2WSH Multisig | ~1,100 | ~275 | Traditional Multisig |
| 2-of-3 MuSig2 (P2TR) | ~560 | ~140 | Signature Aggregation |

---

## 関連ドキュメント

- [BTC Technical Reference](./README.md) - Bitcoin技術リファレンス
- [Taproot User Guide](./TAPROOT_GUIDE.md) - Taprootの使い方
- [MuSig2 User Guide](./musig2_guide.md) - MuSig2の使い方
- [Descriptor Examples](./descriptor_examples.md) - Descriptorの例
- [PSBT Developer Guide](./psbt_developer_guide.md) - PSBT開発ガイド
- [BCH E2E Workflow](../../../scripts/operation/bch/README.md) - BCH E2Eワークフロー

---

**ドキュメントバージョン:** 1.0
**最終更新:** 2026-01-12
**メンテナー:** go-crypto-wallet team
