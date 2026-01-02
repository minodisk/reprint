---
allowed-tools: Bash(gh pr merge:*), Bash(gh pr view:*), Bash(git checkout:*), Bash(git pull:*)
description: Merge the current pull request
---

Merge the current pull request using GitHub CLI.

1. Run `gh pr merge --squash --delete-branch` to squash and merge
2. After merge, checkout main and pull: `git checkout main && git pull`
