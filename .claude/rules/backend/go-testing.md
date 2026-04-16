# Go Testing Guide

Google Go Style Guide (https://google.github.io/styleguide/go/) に準拠する。

## テスト命名

- テスト関数名は `Test` プレフィックス + 大文字始まり
- `TestFunctionName_Scenario` パターンで命名する

## テーブル駆動テスト

- 無名構造体スライスでテストケースを定義する（`name`, 入力値, `want`, `wantErr` フィールド）
- `for _, tt := range tests` でイテレーションする
- `t.Run(tt.name, func(t *testing.T) { ... })` でサブテスト化する
- テストケースの追加が容易な構造を維持する

```go
tests := []struct {
	name    string
	input   int
	want    int
	wantErr error
}{
	{"positive", 10, 5, nil},
	{"zero division", 0, 0, ErrDivideByZero},
}
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		got, err := Divide(tt.input, 2)
		// assertions...
	})
}
```

## アサーション

- assertion library は使用しない。標準の比較演算を使う
- 複雑な構造体比較には `cmp` パッケージ（`github.com/google/go-cmp/cmp`）を使用する

## `t.Error` vs `t.Fatal`

- `t.Error` / `t.Errorf`: テスト継続可能な失敗に使用する
- `t.Fatal` / `t.Fatalf`: セットアップの失敗や、後続チェックが無意味な場合のみ使用する
- 別 goroutine からの fatal 呼び出しは禁止

## 失敗メッセージ

- 形式: `"FuncName(%v) = %v, want %v"`
- actual（実際の結果）を先、expected（期待値）を後に記述する
- 関数名と入力値を含めることで失敗箇所を特定しやすくする

```go
if got != tt.want {
	t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
}
```

## ヘルパー関数

- テストヘルパー関数の冒頭で `t.Helper()` を呼ぶ
- これによりエラー発生箇所がヘルパー内ではなく呼び出し元として報告される
- `*testing.T`、`testing.TB`、カスタムインターフェースを引数に取れる

## クリーンアップ

- `t.Cleanup()` でティアダウン処理を登録する（`defer` より優先）
- テスト終了時（skip、fatal、panic 含む）に確実に実行される
- サブテスト内でもネスト可能（逆順で実行される）

## ベンチマーク

- `BenchmarkXxx(b *testing.B)` パターンで命名する
- `b.N` ループが必須
- `b.Run()` でサブベンチマークを作成する

## テスト構造

- テストロジックは Test 関数内に保持し、過度な抽象化を避ける
- バリデーション関数からはエラーを返すことを優先する（assertion ヘルパーより）
