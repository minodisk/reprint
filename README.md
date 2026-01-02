# reprint

External image uploader CLIs for [deck](https://github.com/k1LoW/deck). Provides tools to upload and delete images to/from external storage services like GCS, S3, and more.

## Why reprint?

[deck](https://github.com/k1LoW/deck) is a tool that converts Markdown to Google Slides. When inserting images into slides, deck needs to upload images to a location accessible by Google Slides.

By default, deck uses Google Drive for image hosting. However, some organizations have policies that prevent sharing Google Drive files externally.

For such environments, deck supports external CLI tools for image upload/delete operations (see [PR #2](https://github.com/minodisk/deck/pull/2)). **reprint** provides CLIs that implement this interface, allowing you to use external storage services as temporary image storage.

reprint uses [Signed URLs](https://cloud.google.com/storage/docs/access-control/signed-urls) for temporary access. The storage bucket does **not** need to be public.

## Supported Storage

| CLI | Storage | Documentation |
|-----|---------|---------------|
| [`reprint-gcs`](cmd/reprint-gcs) | Google Cloud Storage | [README](cmd/reprint-gcs/README.md) |
| `reprint-s3` | Amazon S3 | Not yet implemented |

## License

[MIT](LICENSE)
