# ADR 0001: CLI Framework Selection

## Status

Accepted

## Context

reprint needs a CLI framework to handle subcommands (upload, delete) and flags. We considered the following options:

1. **spf13/cobra** - Full-featured CLI framework
2. **No framework** - Standard library only (os.Args, flag package)

## Decision

We chose **spf13/cobra** for the following reasons:

### 1. Consistency with deck

[deck](https://github.com/k1LoW/deck), the primary consumer of reprint, uses cobra. Using the same framework maintains consistency in the ecosystem.

### 2. Subcommand Support

reprint has multiple subcommands (upload, delete). Cobra provides clean subcommand handling out of the box:

```go
rootCmd.AddCommand(uploadCmd)
rootCmd.AddCommand(deleteCmd)
```

### 3. Flag Management

Cobra integrates well with spf13/pflag for POSIX-compliant flag parsing:
- Persistent flags (shared across subcommands)
- Command-specific flags
- Environment variable binding via viper

### Alternatives Considered

**laminate** (another deck-related tool) does not use a CLI framework. However, laminate has simpler requirements without subcommands. For reprint's use case with multiple commands, cobra provides clearer structure.

## Consequences

- Adds cobra as a dependency
- Familiar patterns for Go developers
- Easy to extend with new subcommands (e.g., reprint-s3)
