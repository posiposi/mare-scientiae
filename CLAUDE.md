# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

過去に読んだ書籍の傾向から次の一冊をおすすめするJSON API。フロントエンドは持たず、curlまたはCLIで実行する。Claude APIを利用してレコメンドロジックを実装する。最終的にはAWS Lambda上でのサーバレス実行を目指す。

## ビルド・静的解析

すべてのコマンドは `api` コンテナ内で実行する。

```bash
# ビルド
docker compose exec api go build ./...

# フォーマット
docker compose exec api gofmt -w .

# vet
docker compose exec api go vet ./...
```

テスト実行は `test-runner` エージェントが `api` コンテナ内で行う。ローカルでの `go test` は使用しない。

## アーキテクチャ

DDD（ドメイン駆動設計）のレイヤードアーキテクチャを採用。実装は必ずTDD（テスト駆動開発）で行う。`/tdd-workflow` スキルを使用してRed→Green→Refactoringのサイクルに従うこと。

```
cmd/
  server/       … HTTPサーバのエントリポイント
  lambda/       … Lambda用エントリポイント（将来）
internal/
  domain/
    model/      … エンティティ・値オブジェクト（外部依存なし）
    repository/ … Repository interface（ドメインが要求する永続化契約）
  usecase/
    interactor/ … ユースケース本体の実装
    port/
      output/   … usecase が要求する外部サービス interface（Claude API 等）
  infrastructure/
    ent/         … Ent スキーマ定義と生成コード
    persistence/ … Repository interface の実装（DB アクセス）
  presentation/
    handler/    … HTTP ハンドラ（消費する usecase interface はハンドラ側で定義）
    dto/        … レスポンス/リクエスト DTO
    router/     … ルーティング定義
```

依存方向: presentation → usecase → domain ← infrastructure

domain 層は外部パッケージに一切依存しない。infrastructure 層は domain のインターフェース（repository）を実装する。

### interface 配置原則

interface の配置は `/interface-placement` スキルに従う。概要:

- 基本は Go idiom に則り **利用側（consumer）** に interface を定義する。handler が呼ぶ usecase interface は handler パッケージ内で定義する
- 例外として **repository interface のみ** `internal/domain/repository/` に配置する（ドメインが要求する永続化契約のため）
- 命名: repository 系は実装 struct 名 + `-er`（例: `BookRepository` → `BookQueryRepositorier`）。usecase / handler 側 interface も `-er` を付ける

### interface 適合チェックのイディオム

`var _ Interface = (*Impl)(nil)` によるコンパイル時の interface 適合チェックは、以下の方針で扱う。

- **追加する**: 実装 struct が**本番コード（非テスト）で** interface 型として利用されている箇所（DI の引数、フィールド、戻り値等）がまだ存在しない場合。チェックが無いとシグネチャ変更時にコンパイルエラーが検出されないため、明示的に宣言する。
- **追加しない**: usecase 層や `cmd/server/` などの本番コードで既に interface 型として利用されている場合。コンパイラが利用箇所で自動的に適合性を検査するため、イディオムは冗長となる。

## Docker 開発環境

docker-compose で以下の3サービスを構成する。

| サービス | イメージ         | 役割                              |
| -------- | ---------------- | --------------------------------- |
| api      | 自前ビルド（Go） | アプリケーション本体              |
| db       | postgres:18.3    | 開発用データベース                |
| test-db  | postgres:18.3    | テスト用データベース（TDDで使用） |

## ホットリロード（air）

`api` サービスは [air](https://github.com/air-verse/air) によるホットリロードで起動する。`./` 配下の `.go` ファイル変更を検知して自動で再ビルド・再起動する。

- 設定: `.air.toml`（ビルドコマンドは `go build -o ./tmp/main ./cmd/server`）
- 監視除外: `tmp`、`.git`、`.claude`、`schemas`、`docs`、`*_test.go`
- 一時ファイル: `tmp/main`（バイナリ）、`build-errors.log`（ビルド失敗ログ）。いずれも `.gitignore` 済み

ビルド失敗時は `build-errors.log` または `docker compose logs api` で確認する。

## スキーマ管理

[Ent](https://entgo.io) の宣言的スキーマを採用。スキーマはGoコードとして `internal/infrastructure/ent/schema/` に定義し、`go generate` でクライアントコードを生成する。マイグレーションは Ent の [auto-migration](https://entgo.io/docs/migrate#auto-migration) を `cmd/ent` の手動CLIから実行する。将来的にテーブル数が増えた段階で Atlas versioned migrations に移行する。

```
cmd/ent/main.go             … マイグレーション実行用CLI（Schema.Create）
internal/infrastructure/ent/
  generate.go               … //go:generate ディレクティブ
  schema/<エンティティ>.go  … スキーマ定義（Goコード）
  <生成ファイル>            … go generate で生成・Gitコミット
```

ent CLI は `api` コンテナイメージに `go install entgo.io/ent/cmd/ent@<version>` で同梱している。

```bash
# ent クライアントコード生成
docker compose exec api go generate ./internal/infrastructure/ent

# マイグレーション実行（auto-migration は append-only）
# DATABASE_URL は docker-compose.yml で api コンテナに注入済み
docker compose exec api go run ./cmd/ent
```

## 開発フロー

以下の手順に従って開発を進める。

### 1. Issue取得

- ユーザーが指定したIssue番号からGitHub Issuesを取得する
- GitHubへのアクセスはMCPサーバー（`mcp__github`）を使用する
- `main`ブランチから、Issue番号とタイトルに適した英語名で作業ブランチを作成する
  - 例: `git switch -c "#1_fix_bugs"`

### 2. 設計

- `/design-workflow` スキルを読み込み、設計フェーズを実行する
- Plan modeを使用し、Issue内容から必要な設計観点を動的に決定してサブエージェントで並列設計する
- 統合Planをユーザーが承認するまで実装フェーズには進まない

### 3. 実装

- `/tdd-workflow` スキルを読み込み、Red→Green→Refactoringのサイクルで実装する
- コミット前に `gofmt -w .` および `go vet ./...` を `api` コンテナ内で実行し、フォーマットと静的解析を通すこと
- コミット時は `/commit-commands:commit` を使用する
  - コミットの粒度は**最低限の機能単位または修正項目**一つずつで行うこと
  - コミットメッセージは日本語で記述するようにコマンド実行時に指示すること
- フロントエンド・バックエンド等を並行して実装できる場合はサブエージェントで並列実装する
- `code-simplifier` プラグインを使用してコードの簡潔さ・可読性を維持する
- `security-guideline` プラグインに準拠し、セキュリティを考慮した実装を行う

### 4. PR作成

- `/pull-request-creation` スキルを読み込み、PR作成フェーズを実行する
- 内部で `/commit-commands:commit-push-pr` を呼び出す
  - git pushおよびPR作成前にユーザー承認を得ること
  - コミットメッセージ・PRタイトル・PR本文は日本語で作成するようにコマンド実行時に指示すること
  - PRタイトルにはIssue番号を含めない

### 5. レビュー

- `/code-review:code-review` プラグインを使用してコードレビューを実施する
- レビューで指摘された項目をユーザーに提示し、修正が必要な箇所の指示を受ける
- 修正は実装と同様にTDD（`/tdd-workflow`）で行う

#### ユーザー承認についての備考

- 4.PR作成と5.レビューはユーザー承認を得た場合、連続で実施する
  - 5.レビューのプラグイン実行にユーザー承認は**不要**

## サブエージェント

### コードベース調査

- コードベースの調査を行う際は`code-analyzer`サブエージェントを利用する
  - フロントエンド、バックエンド等、複数の領域が異なる項目で調査が必要な場合はサブエージェントを並列で起動する
