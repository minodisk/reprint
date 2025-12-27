# reprint

External image uploader CLIs for [deck](https://github.com/k1LoW/deck). Provides tools to upload and delete images to/from external storage services like GCS, S3, and more.

## Why reprint?

[deck](https://github.com/k1LoW/deck) is a tool that converts Markdown to Google Slides. When inserting images into slides, deck needs to upload images to a location with a public URL that Google Slides can access.

By default, deck uses Google Drive for image hosting. However, some organizations have policies that prevent sharing Google Drive files externally. In such environments, you cannot obtain a public URL that Google Slides can access.

For such environments, deck supports external CLI tools for image upload/delete operations (see [PR #2](https://github.com/minodisk/deck/pull/2)). **reprint** provides CLIs that implement this interface, allowing you to use external storage services as temporary image storage.

## Supported Storage

| CLI | Storage |
|-----|---------|
| `mount-gcs` | Google Cloud Storage |
| `mount-s3` | Amazon S3 (coming soon) |

## Installation

```bash
go install github.com/minodisk/mount/cmd/mount-gcs@latest
```

## Usage with deck

```bash
deck apply -u "mount-gcs upload" -d "mount-gcs delete" slide.md
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `MOUNT_BUCKET` | Yes | GCS bucket name |
| `MOUNT_PREFIX` | No | Object prefix (default: empty) |
| `MOUNT_PUBLIC` | No | Generate public URL (`true`/`false`, default: `true`) |

### CLI Flags

You can also use flags instead of environment variables (flags take precedence):

```bash
mount-gcs upload --bucket my-bucket --prefix images/ --public=true
mount-gcs delete --bucket my-bucket
```

## Commands

### upload

Reads image data from stdin and uploads it to GCS.

**Input:**
- stdin: Image binary data
- Environment variables: `DECK_UPLOAD_MIME` (MIME type), `DECK_UPLOAD_FILENAME` (filename)

**Output (stdout):**
```
<public URL>
<resource ID>
```

The resource ID is the GCS object path (`prefix/filename`).

### delete

Deletes the specified object from GCS.

**Input:**
- Environment variable: `DECK_DELETE_ID` (resource ID)

## Authentication

Uses GCP default credentials. Authenticate via:

1. `gcloud auth application-default login`
2. Service account key (`GOOGLE_APPLICATION_CREDENTIALS` environment variable)
3. GCE/Cloud Run metadata server

## Example

```bash
# Configure via environment variables
export MOUNT_BUCKET=my-images-bucket
export MOUNT_PREFIX=deck/

# Use with deck
deck apply -u "mount-gcs upload" -d "mount-gcs delete" presentation.md

# Manual test
export DECK_UPLOAD_MIME=image/png
export DECK_UPLOAD_FILENAME=test.png
cat image.png | mount-gcs upload
# Output:
# https://storage.googleapis.com/my-images-bucket/deck/test.png
# deck/test.png

# Delete
export DECK_DELETE_ID=deck/test.png
mount-gcs delete
```

## License

MIT
