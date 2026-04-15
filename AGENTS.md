# AGENTS.md

## Repo Context
- This is a go project for building a simple expense tracker to track personal expenses and also split them with friends
- This project will provide REST APIs to add, modify expenses. 
- It will also provide reports and CRON jobs which will summarize spending for the user.
- This project is aimed to be a self hosted solution which anybody can deploy on their own infrastructure or on the cloud

## Run Commands
- `go run .` - standard run
- `go run scripts/dev.go` - hot reload with Air (uses `.air.toml` on Windows, `.air.unix.toml` otherwise)

## Environment Variables
- `APP_PORT` - server port (default: 41605, not 8080)
- `DB_PATH` - SQLite database path (default: `expense.db`)
- `APP_SECRET_KEY` - application secret 

## Architecture
- Go 1.25 with standard library `http.ServeMux` (no web framework like Gin/Echo)
- SQLite database with auto-migration on startup (`db.InitDB()` runs `internal/db/schema.sql` via `//go:embed`)
- Clean architecture: handlers -> repository -> database
- No test files exist in this repo

## Key Files
- `main.go` - entrypoint, sets up routes and database
- `internal/config/config.go` - environment variable loading
- `internal/db/db.go` - database connection and schema initialization
- `internal/handlers/` - HTTP handlers for expenses and users
- `internal/repository/` - data access layer

## Routes
- `GET /ping` - health check
- `POST /expenses`, `GET /expenses/{id}`, `GET /expenses?user_id=X`, `PUT /expenses/{id}`, `DELETE /expenses/{id}`
- Similar CRUD routes for users under `/users`