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
  domain/       … エンティティ・値オブジェクト・リポジトリIF（外部依存なし）
  application/  … ユースケース・ポート定義（ドメイン層のみに依存）
  infrastructure/ … 外部サービス実装（Claude API, DB等）
  presentation/ … HTTPハンドラ・ルーティング
```

依存方向: presentation → application → domain ← infrastructure

domain層は外部パッケージに一切依存しない。infrastructure層はdomain/applicationのインターフェースを実装する。

## Docker 開発環境

docker-compose で以下の3サービスを構成する。

| サービス | イメージ         | 役割                              |
| -------- | ---------------- | --------------------------------- |
| api      | 自前ビルド（Go） | アプリケーション本体              |
| db       | postgres:18.3    | 開発用データベース                |
| test-db  | postgres:18.3    | テスト用データベース（TDDで使用） |

## スキーマ管理

[Ent](https://entgo.io) の宣言的スキーマを採用。スキーマはGoコードとして `internal/infrastructure/ent/schema/` に定義し、`go generate` でクライアントコードを生成する。マイグレーションは Ent の [auto-migration](https://entgo.io/docs/migrate#auto-migration) を利用する。

```
internal/infrastructure/ent/
  generate.go               … //go:generate ディレクティブ
  schema/<エンティティ>.go  … スキーマ定義（Goコード）
  <生成ファイル>            … go generate で生成・Gitコミット
```

ent CLI は `api` コンテナイメージに `go install entgo.io/ent/cmd/ent@<version>` で同梱している。

```bash
# ent クライアントコード生成
docker compose exec api go generate ./internal/infrastructure/ent

# マイグレーション（auto-migration）はアプリ起動時に client.Schema.Create(ctx) を呼んで反映する
```

## 開発フロー

以下の手順に従って開発を進める。

### 1. Issue取得

- ユーザーが指定したIssue番号からGitHub Issuesを取得する
- GitHubへのアクセスはMCPサーバー（`mcp__github`）を使用する
- `main`ブランチから、Issue番号とタイトルに適した英語名で作業ブランチを作成する
  - 例: `git switch -c "#1_fix_bugs"`

### 2. 設計

- Plan modeを使用して仕様設計を行う
- サブエージェントを利用し、以下の観点を並列で設計する
  - フロントエンド
  - バックエンド
  - インフラ
  - その他（横断的関心事など）
- 設計完了後、ユーザーに内容を提示し承認を得てから次のステップへ進む

### 3. 実装

- `/tdd-workflow` スキルを読み込み、Red→Green→Refactoringのサイクルで実装する
- コミット前に `gofmt -w .` および `go vet ./...` を `api` コンテナ内で実行し、フォーマットと静的解析を通すこと
- コミット時は `/commit-commands:commit` を使用する
- フロントエンド・バックエンド等を並行して実装できる場合はサブエージェントで並列実装する
- `code-simplifier` プラグインを使用してコードの簡潔さ・可読性を維持する
- `security-guideline` プラグインに準拠し、セキュリティを考慮した実装を行う

### 4. PR作成

- PR作成前のコードpushはユーザーの承認を得てから実行する
- `/commit-commands:commit-push-pr` を使用してPR作成を行う

### 5. レビュー

- `/code-review:code-review` プラグインを使用してコードレビューを実施する
- レビューで指摘された項目をユーザーに提示し、修正が必要な箇所の指示を受ける
- 修正は実装と同様にTDD（`/tdd-workflow`）で行う
