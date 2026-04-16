---
name: schema-apply
description: Atlasスキーマのdry-run確認と適用を実行するワークフロー。スキーマ変更の適用、テーブル追加・変更のDB反映、atlas schema applyに関するタスクで使用する。
---

# Schema Apply Workflow

Atlas宣言的スキーマ管理のdry-run→承認→適用→検証ワークフローを実行する。

## ステップ

### 1. Dry-run（差分SQL確認）

以下のコマンドで生成されるSQLをプレビューする。

```bash
docker compose run --rm atlas schema apply --env local --dry-run
```

出力された差分SQLをユーザーに提示し、意図通りの変更であるかを確認する。「Schema is synced, no changes to be made」と表示された場合は、スキーマ定義ファイルに変更がないことを伝えて終了する。

### 2. 承認確認

dry-runの結果を確認した上で、ユーザーに適用の承認を求める。承認が得られない場合はここで終了する。

### 3. 適用

承認後、以下のコマンドでスキーマを適用する。

```bash
docker compose run --rm atlas schema apply --env local --auto-approve
```

### 4. 検証

適用後、以下のコマンドでテーブル一覧を表示し、変更が正しく反映されたことを確認する。

```bash
docker compose exec db psql -U user -d mare_scientiae -c '\dt'
```

特定テーブルの詳細構造を確認する場合は、`\d テーブル名` を使用する。