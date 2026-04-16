# Go Code Style Guide

Google Go Style Guide (https://google.github.io/styleguide/go/) に準拠する。

## 設計原則

優先順位: **Clarity > Simplicity > Concision > Maintainability > Consistency**

コードの目的と意図が読み手に明確であることを最優先とする。

## 命名規則

- MixedCaps を使用する。アンダースコアは使わない
- 定数も MixedCaps とする（`ALL_CAPS` 禁止）
- パッケージ名は小文字の単一単語（アンダースコア不可）
- レシーバ名は1-2文字の短い省略形とし、型内で一貫させる（例: `Book` → `b`）
- パッケージ名をエクスポート名に繰り返さない（`domain.Book` ✓、`domain.DomainBook` ✗）
- initialisms はケースを統一する（`XMLAPI`、`GRPC`）
- 関数名: 値を返す関数は名詞的、アクションを実行する関数は動詞的に命名する

## パッケージ設計

- `util`、`helper`、`common` などの汎用パッケージ名は禁止
- パッケージごとに明確な単一責務を持たせる
- パッケージ名から内容が推測できる命名にする

## import

- グループ順序: 標準ライブラリ → 外部パッケージ → 内部パッケージ
- 各グループは空行で区切る
- blank import（`_ "pkg"`）は `main` パッケージとテストファイルのみで使用する

```go
import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"helloworld/internal/domain"
)
```

## エラーハンドリング

- `error` は常に最後の戻り値として返す
- sentinel error は `var Err... = errors.New()` で定義する
- エラーラップにはプログラムによる検査が必要な場合 `%w` を使用し、システム境界では `%v` を使用する
- `%w` はエラー文字列の末尾に配置する
- エラーの判定には `errors.Is()` を使用する
- blank identifier（`_`）でのエラー破棄は禁止
- in-band エラーシグナル（-1、空文字列等の戻り値）は避け、追加の戻り値を使う

## 関数シグネチャ

- `context.Context` は第1引数として渡す（struct に埋め込まない）
- 大きな関数は option struct または variadic options で分割する

## フォーマット

- `gofmt` 準拠
- 固定行長の上限は設けない（長い行はリファクタリングで対応する）
- インデント変更前や URL 中での行分割は回避する

## 並行処理

- goroutine のライフタイムを明示的かつ明確にする
- context cancellation または `sync.WaitGroup` で終了を管理する
- 結果を直接返す同期関数を優先し、非同期パターンは必要な場合のみ使用する
