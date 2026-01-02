# ADR 0002: Release Process

## Status

Accepted

## Context

We need a way to distribute reprint CLIs to users. Options considered:

1. **Manual releases** - Build and upload binaries manually
2. **GoReleaser** - Automated releases with GitHub Actions

## Decision

We chose **GoReleaser** for the following reasons:

### 1. Automation

GoReleaser automates the entire release process:
- Build binaries for multiple platforms (linux, darwin, windows × amd64, arm64)
- Create GitHub releases with changelogs
- Upload pre-built binaries

### 2. Unified Versioning

All CLIs (reprint-gcs, reprint-s3) share the same version number:
- Single git tag (e.g., `v1.0.0`) releases all CLIs
- Simpler version management
- Consistent internal package versions

## Release Process

1. Create and push a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions automatically:
   - Builds binaries for all platforms
   - Creates GitHub release with changelog
   - Uploads pre-built binaries to GitHub Releases
   - Updates Homebrew tap formula

## Homebrew Tap Setup (for maintainers)

1. Create repository: `minodisk/homebrew-tap`
2. Create GitHub Personal Access Token with `repo` scope
3. Add token as `HOMEBREW_TAP_GITHUB_TOKEN` secret in reprint repo

## Installation Methods

### 1. Homebrew (macOS/Linux)
```bash
brew tap minodisk/tap
brew install reprint-gcs
```

### 2. Download Binary

Download from [GitHub Releases](https://github.com/minodisk/reprint/releases):
```bash
# Example for macOS arm64
curl -LO https://github.com/minodisk/reprint/releases/download/v1.0.0/reprint-gcs_1.0.0_darwin_arm64.tar.gz
tar xzf reprint-gcs_1.0.0_darwin_arm64.tar.gz
mv reprint-gcs /usr/local/bin/
```

### 2. Go Install
```bash
go install github.com/minodisk/reprint/cmd/reprint-gcs@latest
```

## Consequences

- Automated releases reduce manual work
- Users can download pre-built binaries for their platform
- Changelog is automatically generated from commit messages
