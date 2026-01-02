# reprint-gcs

External image uploader CLI for [deck](https://github.com/k1LoW/deck) using Google Cloud Storage.

## Installation

```bash
go install github.com/minodisk/reprint/cmd/reprint-gcs@latest
```

## Usage with deck

```bash
deck apply -u "reprint-gcs upload --mime {{mime}}" -d "reprint-gcs delete --object-id {{id}}" slide.md
```

## Configuration

Configuration can be set via CLI flags, environment variables, or config file.

| CLI flag        | Environment variable  | Config file   | Required | Description                                                                       |
| --------------- | --------------------- | ------------- | -------- | --------------------------------------------------------------------------------- |
| `--bucket`      | `REPRINT_BUCKET`      | `bucket`      | Yes      | GCS bucket name                                                                   |
| `--prefix`      | `REPRINT_PREFIX`      | `prefix`      | No       | Object prefix (default: empty)                                                    |
| `--credentials` | `REPRINT_CREDENTIALS` | `credentials` | No       | Service account key file path (default: `~/.config/reprint-gcs/credentials.json`) |

**Priority:** CLI flag > Environment variable > Config file > Default path

### Authentication

A service account key file is required. User credentials (`gcloud auth application-default login`) are not supported because signed URLs require a private key for signing.

**Setup:**

1. Create a service account in GCP Console
2. Download the key file (JSON)
3. Place at `~/.config/reprint-gcs/credentials.json`

## Commands

### upload

Reads image data from stdin and uploads it to GCS.

**Input:**

- stdin: Image binary data

| CLI flag | Environment variable | Required | Description     |
| -------- | -------------------- | -------- | --------------- |
| `--mime` | `DECK_UPLOAD_MIME`   | Yes      | Image MIME type |

**Priority:** CLI flag > Environment variable

**Output (stdout):**

```
<Signed URL>
<id>
```

- **Signed URL**: Temporary URL with expiration (default: 15 minutes). The bucket does not need to be public.
- **id**: Auto-generated UUID (e.g., `a1b2c3d4-5678-90ab-cdef-1234567890ab`). Used as GCS object name.

### delete

Deletes the specified object from GCS.

**Input:**

| CLI flag      | Environment variable | Required | Description         |
| ------------- | -------------------- | -------- | ------------------- |
| `--object-id` | `DECK_DELETE_ID`     | Yes      | Object ID to delete |

**Priority:** CLI flag > Environment variable

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

The service account needs the following permissions on the bucket:

| Permission               | Purpose                                    |
| ------------------------ | ------------------------------------------ |
| `storage.objects.create` | Upload objects                             |
| `storage.objects.delete` | Delete objects                             |
| `storage.objects.get`    | Generate Signed URLs                       |
| `storage.buckets.get`    | Check bucket access (for `doctor` command) |

These permissions can be granted with the following roles:

| Role                         | Purpose                                     |
| ---------------------------- | ------------------------------------------- |
| `roles/storage.objectAdmin`  | Upload/delete objects, generate Signed URLs |
| `roles/storage.bucketViewer` | Check bucket access (for `doctor` command)  |

```bash
# Grant permissions to a service account
gcloud storage buckets add-iam-policy-binding gs://your-bucket-name \
  --member=serviceAccount:your-sa@project.iam.gserviceaccount.com \
  --role=roles/storage.objectAdmin

gcloud storage buckets add-iam-policy-binding gs://your-bucket-name \
  --member=serviceAccount:your-sa@project.iam.gserviceaccount.com \
  --role=roles/storage.bucketViewer
```

## Example

### Minimal setup (using default credentials path)

1. Place credentials at `~/.config/reprint-gcs/credentials.json`
2. Create `~/.config/reprint/config.yaml`:

```yaml
bucket: my-images-bucket
```

### Custom credentials path

If you want to use a different credentials path:

```yaml
# ~/.config/reprint/config.yaml
bucket: my-images-bucket
prefix: deck/
credentials: /path/to/service-account-key.json
```

Or via environment variables:

```bash
export REPRINT_BUCKET=my-images-bucket
export REPRINT_PREFIX=deck/
export REPRINT_CREDENTIALS=/path/to/service-account-key.json
```

### Usage

```bash
deck apply -u "reprint-gcs upload --mime {{mime}}" -d "reprint-gcs delete --object-id {{id}}" presentation.md
```
