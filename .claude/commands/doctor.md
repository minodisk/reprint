---
allowed-tools: Bash(go build:*), Bash(./tmp/reprint-gcs doctor:*)
description: Build and run reprint-gcs doctor
---

Build reprint-gcs and run doctor command to check configuration.

1. Build: `go build -o tmp/reprint-gcs ./cmd/reprint-gcs`
2. Run: `./tmp/reprint-gcs doctor`
