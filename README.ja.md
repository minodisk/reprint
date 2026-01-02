# reprint

[deck](https://github.com/k1LoW/deck) 用の外部画像アップローダーCLI群。GCS、S3などの外部ストレージサービスへの画像のアップロード・削除を行います。

## なぜ reprint が必要か？

[deck](https://github.com/k1LoW/deck) は Markdown を Google Slides に変換するツールです。スライドに画像を挿入する際、deck は Google Slides がアクセスできる公開URLを持つ場所に画像をアップロードする必要があります。

デフォルトでは、deck は画像のホスティングに Google Drive を使用します。しかし、組織によっては Google Drive にアップロードしたファイルを外部共有できないポリシーに設定している場合があります。そのような環境では、Google Slides がアクセスできる公開URLを取得できません。

そのような環境向けに、deck は画像のアップロード・削除操作に外部CLIツールを使用できます（[PR #2](https://github.com/minodisk/deck/pull/2) 参照）。**reprint** はこのインターフェースを実装したCLI群で、外部ストレージサービスを一時的な画像ストレージとして使用できるようにします。

## 対応ストレージ

| CLI                              | ストレージ           | ドキュメント                        |
| -------------------------------- | -------------------- | ----------------------------------- |
| [`reprint-gcs`](cmd/reprint-gcs) | Google Cloud Storage | [README](cmd/reprint-gcs/README.md) |
| `reprint-s3`                     | Amazon S3            | 未実装                              |

## ライセンス

[MIT](LICENSE)
