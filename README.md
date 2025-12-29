# reprint

External image uploader CLIs for [deck](https://github.com/k1LoW/deck). Provides tools to upload and delete images to/from external storage services like GCS, S3, and more.

## Why reprint?

[deck](https://github.com/k1LoW/deck) is a tool that converts Markdown to Google Slides. When inserting images into slides, deck needs to upload images to a location with a public URL that Google Slides can access.

By default, deck uses Google Drive for image hosting. However, some organizations have policies that prevent sharing Google Drive files externally. In such environments, you cannot obtain a public URL that Google Slides can access.

For such environments, deck supports external CLI tools for image upload/delete operations (see [PR #2](https://github.com/minodisk/deck/pull/2)). **reprint** provides CLIs that implement this interface, allowing you to use external storage services as temporary image storage.

## Supported Storage

| CLI | Storage |
|-----|---------|
| `reprint-gcs` | Google Cloud Storage |
| `reprint-s3` | Amazon S3 (coming soon) |

## Installation

```bash
go install github.com/minodisk/reprint/cmd/reprint-gcs@latest
```

## Usage with deck

```bash
deck apply -u "reprint-gcs upload" -d "reprint-gcs delete --filename {{filename}}" slide.md
```

## Configuration

Configuration can be set via CLI flags, environment variables, or config file.

**Priority (highest to lowest):** CLI flag > Environment variable > Config file

| CLI flag | Environment variable | Config file | Required | Description |
|----------|---------------------|-------------|----------|-------------|
| `--bucket` | `REPRINT_BUCKET` | `bucket` | Yes | GCS bucket name |
| `--prefix` | `REPRINT_PREFIX` | `prefix` | No | Object prefix (default: empty) |
| `--credentials` | `REPRINT_CREDENTIALS` | `credentials` | No | Service account key file path |

### Authentication

**Priority (highest to lowest):**
1. `--credentials` / `REPRINT_CREDENTIALS` / `credentials` (service account key file)
2. `GOOGLE_APPLICATION_CREDENTIALS` environment variable
3. `gcloud auth application-default login`
4. GCE/Cloud Run metadata server

## Commands

### upload

Reads image data from stdin and uploads it to GCS.

**Input:**
- stdin: Image binary data

| CLI flag | Environment variable | Required | Description |
|----------|---------------------|----------|-------------|
| `--mime` | `DECK_UPLOAD_MIME` | Yes | Image MIME type |

**Output (stdout):**
```
<public URL>
<filename>
```

Filename is auto-generated UUID without extension (e.g., `a1b2c3d4-5678-90ab-cdef-1234567890ab`).

### delete

Deletes the specified object from GCS.

**Input:**

| CLI flag | Environment variable | Required | Description |
|----------|---------------------|----------|-------------|
| `--filename` | `DECK_DELETE_FILENAME` | Yes | Filename to delete |

## Example

### Config file

Create `~/.config/reprint/config.yaml`:

```yaml
bucket: my-images-bucket
prefix: deck/
```

### Environment variables

```bash
export REPRINT_BUCKET=my-images-bucket
export REPRINT_PREFIX=deck/
```

### Usage

```bash
# Use with deck
deck apply -u "reprint-gcs upload" -d "reprint-gcs delete --filename {{filename}}" presentation.md

# Manual upload test
cat image.png | reprint-gcs upload
# Output:
# https://storage.googleapis.com/my-images-bucket/deck/a1b2c3d4-5678-90ab-cdef-1234567890ab
# a1b2c3d4-5678-90ab-cdef-1234567890ab

# Manual delete test
reprint-gcs delete --filename a1b2c3d4-5678-90ab-cdef-1234567890ab
```

## License

[MIT](LICENSE)
