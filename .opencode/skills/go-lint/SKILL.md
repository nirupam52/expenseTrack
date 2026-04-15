---
name: go-lint
description: Runs Go lint, format, and vet checks for code quality
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: ci
---

## What I do

Run Go toolchain checks to ensure code quality and consistency:

1. `go fmt ./...` - Format all Go files
2. `go vet ./...` - Run go vet for static analysis
3. `go mod tidy` - Ensure go.mod and go.sum are consistent
4. `go build ./...` - Verify code compiles without errors

## When to use me

Use this skill after making code changes to ensure:
- Code is properly formatted
- No vet warnings or errors
- Code compiles successfully
- Dependencies are in sync

## How to invoke

```
skill({ name: "go-lint" })
```

Or include in your workflow:
```
After making changes, run go-lint to verify code quality.
```