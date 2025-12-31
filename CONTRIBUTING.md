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

## Architecture Decision Records

Design decisions are documented in `docs/adr/`. See [ADR 0001](docs/adr/0001-cli-framework.md) for an example.
