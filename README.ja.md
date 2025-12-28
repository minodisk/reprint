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
deck apply -u "reprint-gcs upload" -d "reprint-gcs delete" slide.md
```

### 環境変数

| 変数名 | 必須 | 説明 |
|--------|------|------|
| `REPRINT_BUCKET` | Yes | GCSバケット名 |
| `REPRINT_PREFIX` | No | オブジェクトのプレフィックス（デフォルト: 空） |
| `REPRINT_PUBLIC` | No | 公開URLを生成するか（`true`/`false`、デフォルト: `true`） |

### CLIフラグ

環境変数の代わりにフラグでも設定できます（フラグが優先）:

```bash
reprint-gcs upload --bucket my-bucket --prefix images/ --public=true
reprint-gcs delete --bucket my-bucket
```

## コマンド

### upload

stdin から画像データを読み取り、GCS にアップロードします。

**入力:**
- stdin: 画像バイナリデータ
- 環境変数: `DECK_UPLOAD_MIME`（MIMEタイプ）、`DECK_UPLOAD_FILENAME`（ファイル名）

**出力 (stdout):**
```
<公開URL>
<リソースID>
```

リソースIDはGCSオブジェクトのパス（`prefix/filename`）です。

### delete

指定されたオブジェクトをGCSから削除します。

**入力:**
- 環境変数: `DECK_DELETE_ID`（リソースID）

## 認証

GCPのデフォルト認証情報を使用します。以下のいずれかで認証できます:

1. `gcloud auth application-default login`
2. サービスアカウントキー（`GOOGLE_APPLICATION_CREDENTIALS` 環境変数）
3. GCE/Cloud Run のメタデータサーバー

## 使用例

```bash
# 環境変数で設定
export REPRINT_BUCKET=my-images-bucket
export REPRINT_PREFIX=deck/

# deck から使用
deck apply -u "reprint-gcs upload" -d "reprint-gcs delete" presentation.md

# 手動テスト
export DECK_UPLOAD_MIME=image/png
export DECK_UPLOAD_FILENAME=test.png
cat image.png | reprint-gcs upload
# 出力:
# https://storage.googleapis.com/my-images-bucket/deck/test.png
# deck/test.png

# 削除
export DECK_DELETE_ID=deck/test.png
reprint-gcs delete
```

## ライセンス

MIT
