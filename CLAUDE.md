# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

過去に読んだ書籍の傾向から次の一冊をおすすめするJSON API。フロントエンドは持たず、curlまたはCLIで実行する。Claude APIを利用してレコメンドロジックを実装する。最終的にはAWS Lambda上でのサーバレス実行を目指す。

## ビルド・テスト

```bash
# ビルド
go build ./...

# テスト全実行
go test ./...

# 単一パッケージのテスト
go test ./internal/domain/book/

# 特定テスト関数の実行
go test ./internal/domain/book/ -run TestBookEntity

# テスト（カバレッジ付き）
go test -cover ./...

# フォーマット
gofmt -w .

# vet
go vet ./...
```

## アーキテクチャ

DDD（ドメイン駆動設計）のレイヤードアーキテクチャを採用。TDDで実装する。

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
| atlas | arigaio/atlas | スキーマ管理（宣言的ワークフロー） |

## スキーマ管理

Atlas（宣言的スキーマ管理）を採用。Go公式レイアウトに従い、非Goファイルはプロジェクトルートに配置する。

```
atlas.hcl              … Atlas プロジェクト設定
schemas/               … スキーマ宣言（テーブル定義）
  db/
    schema.hcl         … データベース定義（publicスキーマ）
    tables/            … テーブル定義
```
