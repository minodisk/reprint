# reprint-gcs

External image uploader CLI for [deck](https://github.com/k1LoW/deck) using Google Cloud Storage.

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
<Signed URL>
<filename>
```

- **Signed URL**: Temporary URL with expiration (default: 15 minutes). The bucket does not need to be public.
- **filename**: Auto-generated UUID without extension (e.g., `a1b2c3d4-5678-90ab-cdef-1234567890ab`)

### delete

Deletes the specified object from GCS.

**Input:**

| CLI flag | Environment variable | Required | Description |
|----------|---------------------|----------|-------------|
| `--filename` | `DECK_DELETE_FILENAME` | Yes | Filename to delete |

## GCS Bucket Setup

### Creating a Bucket

```bash
gcloud storage buckets create gs://your-bucket-name --location=REGION
```

Choose a region close to your users. See [available locations](https://cloud.google.com/storage/docs/locations).

### Security

**Do NOT make the bucket public.** reprint-gcs uses [Signed URLs](https://cloud.google.com/storage/docs/access-control/signed-urls) for temporary access. deck only needs temporary access to embed images in Google Slides, then deletes the files.

Making the bucket public is a security risk and unnecessary for this use case.

### Required IAM Permissions

The service account or user needs the following permissions on the bucket:

- `storage.objects.create` - Upload objects
- `storage.objects.delete` - Delete objects
- `storage.objects.get` - Generate Signed URLs
- `storage.buckets.get` - Check bucket access (for `doctor` command)

**Recommended role:** `roles/storage.objectAdmin` on the specific bucket.

```bash
# Grant permissions to a service account
gcloud storage buckets add-iam-policy-binding gs://your-bucket-name \
  --member=serviceAccount:your-sa@project.iam.gserviceaccount.com \
  --role=roles/storage.objectAdmin
```

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
# https://storage.googleapis.com/my-images-bucket/deck/a1b2c3d4-...?X-Goog-Algorithm=...&X-Goog-Expires=900&X-Goog-Signature=...
# a1b2c3d4-5678-90ab-cdef-1234567890ab

# Manual delete test
reprint-gcs delete --filename a1b2c3d4-5678-90ab-cdef-1234567890ab
```
