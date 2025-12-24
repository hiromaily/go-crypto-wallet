# Overview (日本語)

このリポジトリは、Bitcoin、Bitcoin Cash、Ethereum、ERC-20 Token、Rippleなどの暗号通貨を扱うウォレットシステムです。
セキュリティを重視した設計により、3つの異なるタイプのウォレット（Watch Wallet、Keygen Wallet、Sign Wallet）を分離し、それぞれが異なる役割を担います。

## ウォレットタイプ

### 1. Watch Wallet（ウォッチウォレット）

**特徴:**

- **オンライン環境**で動作
- **公開鍵のみ**を保持（秘密鍵は保持しない）
- ブロックチェーンノードにアクセス可能

**主な機能:**

- 未署名トランザクションの作成
- 署名済みトランザクションの送信
- トランザクション状態の監視
- アドレスのインポート（Keygen Walletからエクスポートされた公開鍵アドレス）

**主要なCLIコマンド:**

- `watch import address` - アドレスのインポート
- `watch create deposit` - 入金用未署名トランザクションの作成
- `watch create payment` - 出金用未署名トランザクションの作成
- `watch create transfer` - アカウント間送金用未署名トランザクションの作成
- `watch send` - 署名済みトランザクションの送信
- `watch monitor` - トランザクションや残高の監視
- `watch api` - ブロックチェーンノードへのAPI呼び出し（BTC/BCH/ETH/XRP固有）

### 2. Keygen Wallet（キージェンウォレット）

**特徴:**

- **オフライン環境**で動作（コールドウォレット）
- アカウントの鍵管理を担当
- HD Wallet（階層的決定性ウォレット）に基づく鍵生成

**主な機能:**

- アカウント用シードの生成
- HD Walletに基づく鍵の生成
- マルチシグアドレスの生成（アカウント設定に基づく）
- 公開鍵アドレスのエクスポート（CSV形式、Watch Walletにインポート用）
- 未署名トランザクションへの**1回目の署名**（マルチシグの場合）

**主要なCLIコマンド:**

- `keygen create seed` - シードの生成
- `keygen create hdkey` - HD鍵の生成
- `keygen create multisig` - マルチシグアドレスの生成
- `keygen export address` - 公開鍵アドレスのエクスポート
- `keygen import full-pubkey` - Sign Walletからエクスポートされたフル公開鍵のインポート
- `keygen sign` - 未署名トランザクションへの署名（1回目）
- `keygen api` - ウォレット管理API（BTC/BCH/ETH固有）

### 3. Sign Wallet（署名ウォレット）

**特徴:**

- **オフライン環境**で動作（コールドウォレット）
- 認証オペレーターが使用
- 各オペレーターに独自の認証アカウントとSign Walletアプリが提供される

**主な機能:**

- 認証アカウント用シードの生成
- 認証アカウント用HD鍵の生成
- フル公開鍵アドレスのエクスポート（CSV形式、Keygen Walletにインポート用）
- 未署名トランザクションへの**2回目以降の署名**（マルチシグの場合）

**主要なCLIコマンド:**

- `sign create seed` - シードの生成
- `sign create hdkey` - HD鍵の生成
- `sign export fullpubkey` - フル公開鍵アドレスのエクスポート
- `sign import privkey` - 秘密鍵のインポート
- `sign sign` - 未署名トランザクションへの署名（2回目以降）

## トランザクション実行フロー

### 非マルチシグアドレスの場合

```text
1. Watch Wallet: 未署名トランザクションを作成
   └─> watch create deposit/payment/transfer
   └─> トランザクションファイル（未署名）が生成される

2. Keygen Wallet: 1回目の署名を実行
   └─> keygen sign -file <tx_file>
   └─> 署名済みトランザクションファイルが生成される

3. Watch Wallet: 署名済みトランザクションを送信
   └─> watch send -file <signed_tx_file>
   └─> ブロックチェーンに送信され、トランザクションIDが返される
```

### マルチシグアドレスの場合

```text
1. Watch Wallet: 未署名トランザクションを作成
   └─> watch create deposit/payment/transfer
   └─> トランザクションファイル（未署名）が生成される

2. Keygen Wallet: 1回目の署名を実行
   └─> keygen sign -file <tx_file>
   └─> 1回署名済みトランザクションファイルが生成される

3. Sign Wallet #1: 2回目の署名を実行
   └─> sign sign -file <tx_file_signed1>
   └─> 2回署名済みトランザクションファイルが生成される

4. Sign Wallet #2: 3回目の署名を実行（必要に応じて）
   └─> sign sign -file <tx_file_signed2>
   └─> 3回署名済みトランザクションファイルが生成される
   └─> （マルチシグの設定に応じて、必要な署名数まで繰り返し）

5. Watch Wallet: 署名済みトランザクションを送信
   └─> watch send -file <fully_signed_tx_file>
   └─> ブロックチェーンに送信され、トランザクションIDが返される
```

## 鍵生成フロー

マルチシグアドレスを生成するための鍵生成フロー:

```text
1. Keygen Wallet: アカウント用シードを生成
   └─> keygen create seed

2. Keygen Wallet: HD鍵を生成
   └─> keygen create hdkey --account <account>

3. Sign Wallet #1: 認証アカウント用シードを生成
   └─> sign create seed

4. Sign Wallet #1: 認証アカウント用HD鍵を生成
   └─> sign create hdkey

5. Sign Wallet #1: フル公開鍵をエクスポート
   └─> sign export fullpubkey
   └─> CSVファイルが生成される

6. Keygen Wallet: Sign Walletからエクスポートされたフル公開鍵をインポート
   └─> keygen import full-pubkey -file <fullpubkey.csv>

7. Keygen Wallet: マルチシグアドレスを生成
   └─> keygen create multisig --account <account>
   └─> マルチシグアドレスが生成される

8. Keygen Wallet: 公開鍵アドレスをエクスポート
   └─> keygen export address
   └─> CSVファイルが生成される

9. Watch Wallet: Keygen Walletからエクスポートされたアドレスをインポート
   └─> watch import address -file <address.csv>
   └─> Watch Walletがアドレスを監視可能になる
```

## トランザクションタイプ

### Deposit（入金）

クライアントアドレスに送金されたコインを、オフライン管理の安全なアドレス（コールドウォレット）に集約するトランザクション。

```bash
# Watch Walletで未署名トランザクションを作成
watch create deposit

# Keygen Walletで署名
keygen sign -file <tx_file>

# Watch Walletで送信
watch send -file <signed_tx_file>
```

### Payment（出金）

ユーザーの出金リクエストに基づき、指定されたアドレスにコインを送金するトランザクション。

```bash
# Watch Walletで未署名トランザクションを作成
watch create payment

# Keygen Walletで1回目の署名
keygen sign -file <tx_file>

# Sign Walletで2回目以降の署名（マルチシグの場合）
sign sign -file <tx_file_signed1>

# Watch Walletで送信
watch send -file <fully_signed_tx_file>
```

### Transfer（送金）

内部アカウント間でのコイン送金を行うトランザクション。

```bash
# Watch Walletで未署名トランザクションを作成
watch create transfer --account1 <sender> --account2 <receiver> --amount <amount>

# Keygen Walletで署名
keygen sign -file <tx_file>

# Watch Walletで送信
watch send -file <signed_tx_file>
```

## セキュリティ設計

- **秘密鍵の分離**: Watch Walletは秘密鍵を一切保持せず、公開鍵のみを管理
- **オフライン運用**: Keygen WalletとSign Walletはオフライン環境で動作し、ネットワークから隔離
- **マルチシグ**: 複数の署名が必要なマルチシグアドレスにより、単一障害点を排除
- **役割分離**: 鍵生成、署名、送信の各機能を異なるウォレットに分離

## 対応コイン

- **Bitcoin (BTC)**
- **Bitcoin Cash (BCH)**
- **Ethereum (ETH)**
- **ERC-20 Token**
- **Ripple (XRP)**

各コインに対して、上記の3つのウォレットタイプが実装されています。
