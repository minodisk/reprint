# Contributing

## Prerequisites

- Go 1.24.6+
- Docker (for integration tests)
- [mise](https://mise.jdx.dev/) (optional, for task runner)

## Setup

```bash
mise install
```

## Build

```bash
mise run build        # Build all CLIs
mise run build:gcs    # Build reprint-gcs only
```

## Testing

### Unit Tests

```bash
mise run test
```

### Integration Tests

Integration tests require [fake-gcs-server](https://github.com/fsouza/fake-gcs-server) emulator.

**Terminal 1: Start emulator**
```bash
mise run emulator:gcs
```

**Terminal 2: Run tests**
```bash
mise run test:integration
```

**Stop emulator**
```bash
mise run emulator:gcs:stop
# or Ctrl+C in terminal 1
```

### Manual Testing

#### Upload

```bash
cat image.png | ./tmp/reprint-gcs upload --mime image/png
# Output:
# https://storage.googleapis.com/my-bucket/a1b2c3d4-...?X-Goog-Algorithm=...&X-Goog-Expires=900&X-Goog-Signature=...
# a1b2c3d4-5678-90ab-cdef-1234567890ab
```

#### Delete

```bash
./tmp/reprint-gcs delete --object-id a1b2c3d4-5678-90ab-cdef-1234567890ab
```

## Project Structure

```
.
├── cmd/
│   ├── reprint-gcs/       # GCS CLI
│   └── reprint-s3/        # S3 CLI (not yet implemented)
├── internal/
│   ├── config/            # Configuration loading
│   └── gcs/               # GCS client wrapper
└── docs/
    └── adr/               # Architecture Decision Records
```

## Release

Releases are automated with [GoReleaser](https://goreleaser.com/).

```bash
git tag v1.0.0
git push origin v1.0.0
```

This triggers GitHub Actions to:
1. Build binaries for all platforms (linux, darwin, windows × amd64, arm64)
2. Create GitHub release with changelog
3. Upload pre-built binaries
4. Update Homebrew tap formula

See [ADR 0002](docs/adr/0002-release-process.md) for details.

## Architecture Decision Records

Design decisions are documented in `docs/adr/`.
