# reprint

[deck](https://github.com/k1LoW/deck) 用の外部画像アップローダーCLI群。GCS、S3などの外部ストレージサービスへの画像のアップロード・削除を行います。

## なぜ reprint が必要か？

[deck](https://github.com/k1LoW/deck) は Markdown を Google Slides に変換するツールです。スライドに画像を挿入する際、deck は Google Slides がアクセスできる公開URLを持つ場所に画像をアップロードする必要があります。

デフォルトでは、deck は画像のホスティングに Google Drive を使用します。しかし、組織によっては Google Drive にアップロードしたファイルを外部共有できないポリシーに設定している場合があります。そのような環境では、Google Slides がアクセスできる公開URLを取得できません。

そのような環境向けに、deck は画像のアップロード・削除操作に外部CLIツールを使用できます（[PR #2](https://github.com/minodisk/deck/pull/2) 参照）。**reprint** はこのインターフェースを実装したCLI群で、外部ストレージサービスを一時的な画像ストレージとして使用できるようにします。

## 対応ストレージ

| CLI | ストレージ |
|-----|-----------|
| `reprint-gcs` | Google Cloud Storage |
| `reprint-s3` | Amazon S3（近日対応予定） |

## インストール

```bash
go install github.com/minodisk/reprint/cmd/reprint-gcs@latest
```

## deck での使用方法

```bash
deck apply -u "reprint-gcs upload" -d "reprint-gcs delete --filename {{filename}}" slide.md
```

## 設定

設定はCLIフラグ、環境変数、設定ファイルで指定できます。

**優先順位（高い順）:** CLIフラグ > 環境変数 > 設定ファイル

| CLIフラグ | 環境変数 | 設定ファイル | 必須 | 説明 |
|----------|----------|-------------|------|------|
| `--bucket` | `REPRINT_BUCKET` | `bucket` | Yes | GCSバケット名 |
| `--prefix` | `REPRINT_PREFIX` | `prefix` | No | オブジェクトのプレフィックス（デフォルト: 空） |
| `--credentials` | `REPRINT_CREDENTIALS` | `credentials` | No | サービスアカウントキーファイルのパス |

### 認証

**優先順位（高い順）:**
1. `--credentials` / `REPRINT_CREDENTIALS` / `credentials`（サービスアカウントキーファイル）
2. `GOOGLE_APPLICATION_CREDENTIALS` 環境変数
3. `gcloud auth application-default login`
4. GCE/Cloud Run のメタデータサーバー

## コマンド

### upload

stdin から画像データを読み取り、GCS にアップロードします。

**入力:**
- stdin: 画像バイナリデータ

| CLIフラグ | 環境変数 | 必須 | 説明 |
|----------|----------|------|------|
| `--mime` | `DECK_UPLOAD_MIME` | Yes | 画像のMIMEタイプ |

**出力 (stdout):**
```
<公開URL>
<ファイル名>
```

ファイル名は拡張子なしのUUID（例: `a1b2c3d4-5678-90ab-cdef-1234567890ab`）。

### delete

指定されたオブジェクトをGCSから削除します。

**入力:**

| CLIフラグ | 環境変数 | 必須 | 説明 |
|----------|----------|------|------|
| `--filename` | `DECK_DELETE_FILENAME` | Yes | 削除するファイル名 |

## 使用例

### 設定ファイル

`~/.config/reprint/config.yaml` を作成:

```yaml
bucket: my-images-bucket
prefix: deck/
```

### 環境変数

```bash
export REPRINT_BUCKET=my-images-bucket
export REPRINT_PREFIX=deck/
```

### 使用方法

```bash
# deck から使用
deck apply -u "reprint-gcs upload" -d "reprint-gcs delete --filename {{filename}}" presentation.md

# 手動アップロードテスト
cat image.png | reprint-gcs upload
# 出力:
# https://storage.googleapis.com/my-images-bucket/deck/a1b2c3d4-5678-90ab-cdef-1234567890ab
# a1b2c3d4-5678-90ab-cdef-1234567890ab

# 手動削除テスト
reprint-gcs delete --filename a1b2c3d4-5678-90ab-cdef-1234567890ab
```

## ライセンス

[MIT](LICENSE)
