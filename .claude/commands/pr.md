---
allowed-tools: Bash(git checkout:*), Bash(git add:*), Bash(git commit:*), Bash(git push:*), Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(gh pr create:*)
description: Create a pull request from current changes
---

Create a pull request from the current changes.

1. Check `git status` and `git diff` to understand current changes
2. Create a new branch with a descriptive name based on the changes
3. Stage and commit the changes with an appropriate message
4. Push the branch to origin
5. Create a PR using `gh pr create` with a summary and test plan
