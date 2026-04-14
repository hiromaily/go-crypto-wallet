---
task_type: feature-add
description: 新機能追加タスク用のコンテキスト
version: 1.0.0
---

# Task: Feature Add (新機能追加)

## When to Use

- 新しいUse Caseを追加する場合
- 新しいCLIコマンドを追加する場合
- 新しいAPIエンドポイントを追加する場合
- 既存機能を大幅に拡張する場合
- 新しい暗号通貨のサポートを追加する場合

## Required Context Documents

### Must Read (必須)

以下のドキュメントは必ず読み込んでください：

| Document | Path | Purpose |
|----------|------|---------|
| Architecture | `docs/guidelines/architecture.md` | Clean Architecture原則、レイヤー分離 |
| Coding Standards | `docs/guidelines/coding-standards.md` | コーディング規約、命名規則 |
| Internal Guidelines | `internal/AGENTS.md` | internalパッケージ構造 |
| Workflow Guidelines | `docs/guidelines/workflow.md` | 検証ステップ、PRフロー |

### Conditional Read (条件付き)

条件に応じて追加で読み込むドキュメント：

| Condition | Document | Path |
|-----------|----------|------|
| DB変更が必要 | Database | `docs/database/db-management.md` |
| 複数チェーン対応 | Multi-Chain | `docs/guidelines/multi-chain.md` |
| セキュリティ関連機能 | Core Principles | `docs/guidelines/core.md` |
| コード生成が必要 | Code Generation | `docs/guidelines/code-generation.md` |
| 新規テスト戦略 | Testing | `docs/guidelines/testing.md` |

## Task-Specific Rules

### 1. レイヤー分離の遵守

```
Domain Layer (internal/domain/)
    ↑ 依存しない
Application Layer (internal/application/)
    ↑ Portsで抽象化
Infrastructure Layer (internal/infrastructure/)
```

- **Domain**: ビジネスロジック、エンティティ（外部依存なし）
- **Application**: Use Case、DTO、Ports（インターフェース定義）
- **Infrastructure**: 外部サービス実装（DB、API、ファイルIO）
- **Interface Adapters**: CLI、HTTP、ウォレットアダプター

### 2. ファイル配置規則

```
新しいUse Case:
  internal/application/usecase/{wallet-type}/{coin}/{usecase-name}.go

新しいPort (Interface):
  internal/application/ports/{port-name}.go

新しいCLIコマンド:
  internal/interface-adapters/cli/{wallet-type}/{coin}/{command}.go

新しいRepository実装:
  internal/infrastructure/repository/{entity}/{repository-name}.go
```

### 3. 依存性注入

```go
// ❌ Bad: 具象クラスに直接依存
type MyUseCase struct {
    repo *MySQLRepository
}

// ✅ Good: インターフェースに依存
type MyUseCase struct {
    repo ports.Repository
}

func NewMyUseCase(repo ports.Repository) *MyUseCase {
    return &MyUseCase{repo: repo}
}
```

### 4. エラーハンドリング

```go
// ✅ エラーラップの標準パターン
if err != nil {
    return fmt.Errorf("failed to create transaction: %w", err)
}
```

### 5. セキュリティ考慮

- **NEVER**: 秘密鍵やセンシティブ情報をログに出力しない
- **ALWAYS**: 鍵操作はオフラインウォレット（keygen/sign）で行う
- **CHECK**: 新機能がオフライン操作に影響しないか確認

## Pre-Task Checklist

- [ ] 影響するレイヤー（domain/application/infrastructure）を特定したか
- [ ] 既存の類似実装を確認したか
- [ ] 必要なインターフェース（Ports）を特定したか
- [ ] テストファイルの配置場所を確認したか
- [ ] セキュリティ影響を評価したか
- [ ] DIコンテナへの登録が必要か確認したか

## Verification Commands

```bash
# 必須の検証コマンド
make go-lint      # リンターチェック
make tidy         # 依存関係の整理
make check-build  # ビルド確認
make go-test       # テスト実行

# 追加の検証（該当する場合）
make go-test-integration  # 統合テスト
make sqlc                # SQLCコード生成（DB変更時）
make mockery             # モック生成（新インターフェース追加時）
```

## Feature Add Workflow

### Step 1: 要件分析

1. 機能要件を明確にする
2. 影響するレイヤーを特定
3. 既存の類似実装を確認

### Step 2: 設計

1. **Domain**: 必要なエンティティ/値オブジェクトを設計
2. **Ports**: 必要なインターフェースを定義
3. **Use Case**: ビジネスロジックフローを設計
4. **Infrastructure**: 外部依存の実装方針を決定

### Step 3: 実装順序

```
1. Domain層（エンティティ、値オブジェクト）
   ↓
2. Application層（Ports、DTO、Use Case）
   ↓
3. Infrastructure層（Repository、API Client）
   ↓
4. Interface Adapters層（CLI、HTTP Handler）
   ↓
5. DI（internal/di/container.go）
```

### Step 4: テスト作成

```
Unit Tests:
  - {package}_test.go ファイルに配置
  - Use Caseのテストはモックを使用

Integration Tests:
  - internal/integration_test/ に配置
  - または test/integration/ に配置
```

### Step 5: 検証

```bash
# 全ての検証コマンドを実行
make go-lint && make tidy && make check-build && make go-test
```

### Step 6: PRドラフト

```bash
git add .
git commit -m "feat: add {feature description}

- {implementation detail 1}
- {implementation detail 2}

Related to #{issue_number}"

gh pr create --title "Feature: {feature title}"
```

## Testing Requirements

### 必須

- [ ] Use Caseの単体テスト（モック使用）
- [ ] 新しいドメインロジックのテスト
- [ ] エラーケースのテスト

### 推奨

- [ ] 統合テスト
- [ ] エッジケースのテスト
- [ ] パフォーマンステスト（必要な場合）

## Examples

### Example 1: 新しいUse Case追加

```
User: "BTCにMuSig2署名のUse Caseを追加して"

Agent Actions:
1. docs/guidelines/architecture.md を読み込み
2. internal/application/usecase/sign/btc/ の既存実装を確認
3. docs/guidelines/multi-chain.md を読み込み（BTC固有の考慮事項）
4. Portsにインターフェース定義があるか確認
5. 必要に応じて新しいPortを定義
6. Use Case実装
7. テスト作成
8. DI登録
9. 検証コマンド実行
```

### Example 2: 新しいCLIコマンド追加

```
User: "残高確認のCLIコマンドを追加して"

Agent Actions:
1. docs/guidelines/architecture.md を読み込み
2. internal/interface-adapters/cli/ の既存実装パターンを確認
3. 既存のUse Caseを使用するか、新規作成が必要か判断
4. CLIコマンド実装（cobra使用）
5. 必要に応じてUse Case追加
6. テスト作成
7. 検証コマンド実行
```

### Example 3: 新しい暗号通貨サポート

```
User: "Cosmosチェーンのサポートを追加して"

Agent Actions:
1. docs/guidelines/architecture.md を読み込み
2. docs/guidelines/multi-chain.md を読み込み
3. 既存チェーン（BTC/ETH/XRP）の実装パターンを調査
4. Domain層のエンティティ定義
5. Infrastructure層のAPI Client実装
6. Use Case実装
7. CLIコマンド実装
8. 設定ファイル追加
9. テスト作成
10. ドキュメント更新
```

## Architecture Reference

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│           Interface Adapters            │
│    (CLI, HTTP, Wallet Adapters)         │
├─────────────────────────────────────────┤
│            Application                  │
│    (Use Cases, Ports, DTOs)             │
├─────────────────────────────────────────┤
│              Domain                     │
│    (Entities, Value Objects)            │
├─────────────────────────────────────────┤
│           Infrastructure                │
│  (DB, API, Storage, External Services)  │
└─────────────────────────────────────────┘

依存の方向: 外側 → 内側 (Domain層は何にも依存しない)
```

### Directory Mapping

| Layer | Directory |
|-------|-----------|
| Domain | `internal/domain/` |
| Application | `internal/application/` |
| Infrastructure | `internal/infrastructure/` |
| Interface Adapters | `internal/interface-adapters/` |
| DI | `internal/di/` |

## Related Commands

- `/fix-issue #{issue_number}` - Issue対応時のワークフロー
- `gh issue create` - 機能要求のIssue作成

## Related Documents

- [Architecture Guidelines](../../guidelines/architecture.md) - Clean Architecture詳細
- [Coding Standards](/guidelines/coding-conventions) - コーディング規約詳細
- `internal/AGENTS.md` - internalパッケージ詳細
- [Multi-Chain Support](../../guidelines/multi-chain.md) - マルチチェーン対応詳細
