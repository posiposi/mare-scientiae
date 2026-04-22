---
name: interface-placement
description: Go のレイヤードアーキテクチャ（DDD）で interface をどのパッケージに配置するかを判断するためのプロジェクト固有原則。基本は Go idiom に従い「利用側（consumer）」に配置し、例外として repository interface のみ `internal/domain/repository/` に置く。新しい interface を定義するとき、layer 間の依存境界を設計するとき、「この interface はどこに置くのが正しいか」で迷ったとき、usecase / handler / infrastructure で抽象を追加するとき、DI 対象の interface を切り出すとき、port を新設するときに必ず使用する。判断を飛ばしてコピー＆ペーストで既存場所に置かず、このスキルを参照して配置先を決めること。
---

# Interface Placement Principle

本プロジェクト `mare-scientiae` における interface 配置の判断原則。Go のレイヤードアーキテクチャで「どのパッケージに interface を宣言するか」を決定するために使用する。

## 核となる原則

### 原則 A: interface は**利用側（consumer）**に配置する

Go の慣用「Accept interfaces, return structs」に従い、interface はそれを**呼び出す側**のパッケージで宣言する。実装側は暗黙的（implicit satisfaction）に満たす。

- handler が usecase を呼ぶ → usecase の interface は handler パッケージ内で定義
- usecase が外部サービス（Claude API 等）を呼ぶ → 外部サービスの interface は usecase 側（`internal/usecase/port/output/`）で定義
- infrastructure の impl 側には interface を宣言しない

### 原則 B: repository interface のみ `internal/domain/repository/` に配置する（例外）

repository は DDD における「ドメインが要求する永続化契約」であり、**ドメイン概念の所有物**である。そのため利用側である usecase に interface を寄せず、domain 層で一元管理する。

理由:
- 複数 usecase から同じ repository が再利用される想定。どの usecase も「消費者」であり、単独の所有者と見做せない
- repository interface はエンティティとライフサイクルを共にする。model と同じツリー（`internal/domain/`）に置くことで、ドメイン変更時の影響が見えやすい
- DDD 文献および一般的な Go DDD 実装の慣例に合致

CQRS は interface レベルでは分離（`BookQueryRepositorier` と `BookCommandRepositorier` を別 interface に）するが、**同一ファイル** `book_repository.go` 内に同居させる。impl 側の struct は `BookRepository` として Query/Command を単一の型で担う。

## ディレクトリ配置チャート

| interface の役割 | 宣言場所 | 理由 |
| --- | --- | --- |
| handler が呼ぶ usecase | `internal/presentation/handler/<xxx>_handler.go` | handler が消費者 |
| usecase が呼ぶ外部サービス（Claude API、メール等） | `internal/usecase/port/output/<xxx>.go` | usecase が消費者、port 化して明示 |
| usecase が呼ぶ repository | `internal/domain/repository/<xxx>_repository.go` | **例外ルール**（原則 B） |
| cmd/server が呼ぶ usecase（DI 配線） | 定義不要。interactor struct を直接渡す | handler 側 interface を cmd/server で再宣言する必要はない |

## 具体例（本プロジェクトの既存コード）

### 原則 A の例: handler が定義する usecase interface

`internal/presentation/handler/book_handler.go`:

```go
package handler

type ListBooksUsecaser interface {
    Execute(ctx context.Context) ([]*model.Book, error)
}

type BookHandler struct {
    listBooksUsecase ListBooksUsecaser
}
```

実装側 `internal/usecase/interactor/list_books.go` は interface を import しない:

```go
package interactor

type ListBooksInteractor struct { ... }

func (u *ListBooksInteractor) Execute(ctx context.Context) ([]*model.Book, error) { ... }
```

`cmd/server/main.go` で DI:

```go
listBooksInteractor := interactor.NewListBooksInteractor(bookRepo)
bookHandler := handler.NewBookHandler(listBooksInteractor) // implicit satisfaction
```

### 原則 B の例: domain が所有する repository interface

`internal/domain/repository/book_repository.go`:

```go
package repository

type BookQueryRepositorier interface {
    FindAll(ctx context.Context) ([]*model.Book, error)
}

// 将来:
// type BookCommandRepositorier interface { Save(...); Delete(...) }
```

実装側 `internal/infrastructure/persistence/book_repository.go`:

```go
package persistence

type BookRepository struct { client *ent.Client }

func (r *BookRepository) FindAll(ctx context.Context) ([]*model.Book, error) { ... }
// Command メソッドもここに同居
```

usecase 側は `repository.BookQueryRepositorier` を import して依存する:

```go
package interactor

import "helloworld/internal/domain/repository"

type ListBooksInteractor struct {
    repo repository.BookQueryRepositorier
}
```

## 命名規則

- **repository 系 interface**: 実装 struct 名 + `-er` サフィックス
  - 実装 `BookRepository` → interface `BookQueryRepositorier` / `BookCommandRepositorier`（CQRS で分離）
- **利用側 interface（handler/usecase 等）**: 内容を表す名詞 + `-er`
  - 例: `ListBooksUsecaser`, `ClaudeRecommender`（想定）
- **標準 Go 慣用の interface**（io 系など）: そのままの慣用を踏襲

## `var _ Interface = (*Impl)(nil)` イディオムの扱い

コンパイル時の interface 適合チェック `var _ Interface = (*Impl)(nil)` は、**本番コードで interface 型として利用される箇所が無い場合のみ**追加する。

- **追加する**: 実装 struct が DI の引数・フィールド・戻り値として interface 型で渡されている箇所がまだ無い（将来のシグネチャ変更を検知できないため明示）
- **追加しない**: `cmd/server/` や usecase 層で既に interface 型として消費されている（コンパイラが利用箇所で検査するため冗長）

本プロジェクトの既存 repository は cmd/server で DI されているため `var _` 不要。

## 判断フローチャート

新しい interface を宣言するとき、以下の順で判断する。

1. **これは repository か？**（ドメインの永続化契約）
   - Yes → `internal/domain/repository/` に配置（原則 B）
   - No → 2 へ
2. **誰が呼ぶか？**
   - handler が呼ぶ → handler パッケージ内で定義
   - usecase が呼ぶ外部サービス → `internal/usecase/port/output/` に定義
   - それ以外のレイヤー内部抽象 → 呼ぶ側のパッケージで定義
3. **`-er` サフィックスを付けたか？**（命名規則）
4. **`var _ ...` は必要か？**（DI 利用箇所の有無で判断）

## 原則違反のサイン

以下のコードを書こうとしたら、このスキルを読み直すこと。

- infrastructure 層に interface を宣言している（repository impl 側に interface を置くのは誤り）
- `internal/usecase/` 直下に `port/input` のような「handler 向け interface 置き場」を作ろうとしている（handler 側で定義すべき）
- repository interface を `internal/usecase/port/` 以下に置こうとしている（原則 B に反する）
- `BookRepositoryInterface` のように `-Interface` サフィックスを使っている（Go 慣用に反する、`-er` を使う）
