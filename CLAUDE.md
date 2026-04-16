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

docker-compose で以下の4サービスを構成する。

| サービス | イメージ | 役割 |
|---|---|---|
| api | 自前ビルド（Go） | アプリケーション本体 |
| db | postgres:18.3 | 開発用データベース |
| test-db | postgres:18.3 | テスト用データベース（TDDで使用） |
| atlas | arigaio/atlas:1.2.0 | スキーマ管理（宣言的ワークフロー） |

## スキーマ管理

Atlas（宣言的スキーマ管理）を採用。Go公式レイアウトに従い、非Goファイルはプロジェクトルートに配置する。スキーマ適用は `/schema-apply` スキルを使用すること。

```
atlas.hcl              … Atlas プロジェクト設定（srcは配列で複数ディレクトリを指定）
schemas/               … スキーマ宣言（テーブル定義）
  db/
    schema.hcl         … データベース定義（publicスキーマ）
    tables/            … テーブル定義（*.pg.hcl）
```

テーブル定義ファイルは `schemas/db/tables/` に `<テーブル名>.pg.hcl` の命名規則で配置する。Atlasはディレクトリを再帰的にスキャンしないため、`atlas.hcl` の `src` に `schemas/db` と `schemas/db/tables` の両方を配列で指定している。

```bash
# スキーマ適用（dry-run）
docker compose run --rm atlas schema apply --env local --dry-run

# スキーマ適用
docker compose run --rm atlas schema apply --env local --auto-approve
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

### 4. レビュー

- `code-review` コマンドを使用してコードレビューを実施する
- レビューで指摘された項目をユーザーに提示し、修正が必要な箇所の指示を受ける
- 修正は実装と同様にTDD（`/tdd-workflow`）で行う

### 5. PR作成

- PR作成前のコードpushはユーザーの承認を得てから実行する
- `/commit-commands:commit-push-pr` を使用してPR作成を行う
