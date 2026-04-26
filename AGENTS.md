# AGENTS.md

> Primary project context lives in `CLAUDE.md`. Read that first.

## Quick Reference

**Run:** `go run .` · Hot reload: `go run scripts/dev.go`

**Stack:** Go 1.25 · `net/http` · SQLite (`modernc.org/sqlite`)

**Layout:** `internal/handlers` → `internal/repository` → `internal/db`

**Env vars:** `APP_PORT` (default `41605`) · `DB_PATH` (default `expense.db`) · `APP_SECRET_KEY`
