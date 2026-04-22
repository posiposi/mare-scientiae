# mare-scientiae

過去に読んだ書籍の傾向から、次の一冊をおすすめする JSON API。フロントエンドは持たず、curl または CLI で実行する。レコメンドロジックは Claude API で実装し、将来的には AWS Lambda 上でのサーバレス実行を目指す。

> [!NOTE]
> 本リポジトリは Go 言語およびクリーンアーキテクチャ（DDD レイヤードアーキテクチャ）のキャッチアップを目的とした学習プロジェクトです。プロダクション利用は想定していません。

## 技術スタック

- Go 1.26
- PostgreSQL 18.3
- [Ent](https://entgo.io)（ORM・スキーママイグレーション）
- Docker Compose
- Claude API（将来的に導入予定）
- AWS Lambda（将来的に導入予定）

## アーキテクチャ

DDD レイヤードアーキテクチャを採用し、依存方向は `presentation → usecase → domain ← infrastructure` を遵守する。実装は TDD（Red → Green → Refactoring）で行う。

```
cmd/
  server/                 … HTTPサーバのエントリポイント（Composition Root）
  lambda/                 … Lambda用エントリポイント（将来）
  ent/                    … Ent マイグレーション実行用CLI
internal/
  domain/
    model/                … エンティティ・値オブジェクト（外部依存なし）
    repository/           … Repository interface + sentinel error
  usecase/
    interactor/           … ユースケース本体の実装
    port/output/          … usecase が要求する外部サービス interface
  infrastructure/
    ent/                  … Ent スキーマ定義と生成コード
    persistence/          … Repository interface の実装
  presentation/
    handler/              … HTTP ハンドラ（consumer interface も同居）
    dto/                  … レスポンス/リクエスト DTO
    router/               … ルーティング定義
```

詳細な設計指針・命名規則・interface 配置原則は [CLAUDE.md](./CLAUDE.md) および `.claude/rules/backend/` を参照。

## 実装済みエンドポイント

| メソッド | パス | 概要 | 対応 Issue |
| --- | --- | --- | --- |
| `GET` | `/v1/books` | 書籍一覧を取得 | [#19](https://github.com/posiposi/mare-scientiae/issues/19) |
| `GET` | `/v1/books/{id}` | UUID をキーに書籍を 1 件取得 | [#24](https://github.com/posiposi/mare-scientiae/issues/24) |

### レスポンス例

```bash
# 一覧取得
curl http://localhost:8080/v1/books

# 特定書籍取得
curl http://localhost:8080/v1/books/87c50373-9d3e-4f71-8a92-b50af0633b1c
```

| ステータス | ケース | レスポンス |
| --- | --- | --- |
| 200 | 正常 | `BookResponse` JSON |
| 400 | UUID 形式でない `id` | `{"error":"invalid_book_id"}` |
| 404 | 該当レコードなし | `{"error":"book_not_found"}` |
| 500 | 内部エラー | `{"error":"internal_server_error"}` |

## これまでの主な実装履歴

| # | タイトル | PR | 内容 |
| --- | --- | --- | --- |
| #13 | BookQueryRepository の追加 | [#18](https://github.com/posiposi/mare-scientiae/pull/18) | Book ドメインと一覧取得 Repository の基盤整備 |
| #14 | Ent 導入 | [#16](https://github.com/posiposi/mare-scientiae/pull/16) | スキーマ管理ツールに Ent を採用 |
| #17 | CI パイプライン整備 | [#21](https://github.com/posiposi/mare-scientiae/pull/21) | ビルド・vet・テストを GitHub Actions で実行 |
| #19 | 書籍一覧取得 API | [#23](https://github.com/posiposi/mare-scientiae/pull/23) | `GET /v1/books` 実装 |
| #24 | 特定書籍取得 API | [#27](https://github.com/posiposi/mare-scientiae/pull/27) | `GET /v1/books/{id}` 実装 |

## 開発環境

`api` / `db` / `test-db` の 3 サービスを Docker Compose で起動する。

```bash
# 初回・依存更新時
docker compose up -d --build

# 起動
docker compose up -d

# ビルド・静的解析・テストはすべて api コンテナ内で実行
docker compose exec api go build ./...
docker compose exec api gofmt -w .
docker compose exec api go vet ./...
docker compose exec api go test ./...

# Ent クライアントコード生成
docker compose exec api go generate ./internal/infrastructure/ent

# マイグレーション実行（auto-migration は append-only）
docker compose exec api go run ./cmd/ent
```

> [!IMPORTANT]
> 現状ホットリロードは導入されていません。コード変更を反映するには `docker compose restart api` が必要です（[#28](https://github.com/posiposi/mare-scientiae/issues/28) で air の導入を検討中）。

## 開発フロー

Issue 取得 → 設計（`/design-workflow`）→ TDD 実装（`/tdd-workflow`）→ PR 作成（`/pull-request-creation`）→ レビュー（`/code-review:code-review`）の 5 段階で進める。詳細は [CLAUDE.md](./CLAUDE.md) を参照。
