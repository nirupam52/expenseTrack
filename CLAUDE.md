# CLAUDE.md

## Project Overview
Self-hosted personal expense tracker built in Go with minimal dependencies. Goal: deploy on personal or cloud infrastructure to track personal expenses and split them with friends. Build incrementally, backend first, then frontend.

## Tech Stack
- Go 1.25, `net/http` standard library — no web framework
- SQLite via `modernc.org/sqlite` (pure Go, no CGo)
- `golang.org/x/crypto/bcrypt` for password hashing

## Commands
- `go run .` — start server
- `go run scripts/dev.go` — hot reload (`.air.toml` on Windows, `.air.unix.toml` on Unix)
- `go build ./...` — verify compilation
- `go vet ./...` — static analysis
- `go fmt ./...` — format

## Environment Variables
| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `41605` | Server port |
| `DB_PATH` | `expense.db` | SQLite file path |
| `APP_SECRET_KEY` | — | Session token signing key (required for auth) |

## Architecture
```
main.go
  └─ internal/
       ├─ config/     env var loading
       ├─ db/         connection + schema migration (schema.sql embedded via //go:embed)
       ├─ models/     shared data types (User, Expense, ...)
       ├─ handlers/   HTTP layer: parse input, validate, call repo, write response
       ├─ repository/ DB access: parameterized queries, error wrapping
       └─ response/   JSON helpers: WriteSuccess, WriteList, WriteError
```

## Code Conventions

### Context — critical
Always use `r.Context()` in handlers. Never `context.Background()`.

### Response helpers (`internal/response/response.go`)
```go
response.WriteSuccess(w, http.StatusOK, data)     // wraps in {success:true, data:...}
response.WriteList(w, items)                       // wraps in {success:true, data:[...], meta:{count:N}}
response.WriteError(w, http.StatusBadRequest, msg) // wraps in {success:false, error:...}
```

### Handler pattern
```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var input InputType
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        response.WriteError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    // validate required fields, then:
    result, err := h.repo.Action(r.Context(), ...)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            response.WriteError(w, http.StatusNotFound, "resource not found")
            return
        }
        response.WriteError(w, http.StatusInternalServerError, "failed to ...")
        return
    }
    response.WriteSuccess(w, http.StatusCreated, result)
}
```

### Repository pattern
```go
var ErrNotFound = errors.New("not found")  // defined once in the package

func (r *Repo) GetByID(ctx context.Context, id int64) (*models.X, error) {
    err := r.db.QueryRowContext(ctx, query, id).Scan(...)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("get x by id: %w", err)
    }
    return &x, nil
}
```

### Error messages
- Lowercase, no trailing punctuation: `"user not found"` not `"User Not Found"`
- Always wrap DB errors: `fmt.Errorf("create user: %w", err)`
- Never expose internal error details in HTTP responses

### Route registration
Each handler has a `RegisterRoutes(mux *http.ServeMux) error` method. Call it from `main.go`.

### HTTP Status Codes
`200` OK · `201` Created · `400` Bad Request · `401` Unauthorized · `403` Forbidden · `404` Not Found · `409` Conflict · `500` Internal Server Error

## Agents (`.claude/agents/`)

| Agent | When to use |
|---|---|
| `planner` | Before any feature — reads codebase, produces a concrete plan for approval |
| `go-builder` | After plan is approved — implements following conventions, verifies build |
| `go-reviewer` | After builder finishes — strict read-only review with file:line feedback |
| `security-auditor` | After reviewer approves — scans for SQL injection, secrets, auth issues |

Typical cycle: **planner → (approval) → go-builder → go-reviewer → security-auditor**

## Checkpoint Roadmap

### CP1: Foundation ✅
Users (register, get, list) + Expenses (CRUD, list by user), SQLite schema, response helpers.

### CP2: Authentication
- `POST /auth/login` — verify password, create session, return token
- `POST /auth/logout` — delete session
- Auth middleware: validate `Authorization: Bearer <token>` header
- Protect all routes except `POST /users/register` and `POST /auth/login`

### CP3: Groups
- `POST /groups` — create group
- `GET /groups/{id}` — get group with members
- `POST /groups/{id}/members` — add member
- `DELETE /groups/{id}/members/{uid}` — remove member
- `GET /groups/{id}/expenses` — list group expenses

### CP4: Expense Splits
- Create expense with optional `splits` array (who owes what)
- `GET /expenses/{id}/splits` — view splits for an expense
- `POST /splits/{id}/settle` — mark a split as settled
- `GET /users/{id}/balances` — net balance per person

### CP5: Insights
- Add `category` field to expenses (requires schema migration)
- `GET /reports/summary?period=month&user_id=X`
- `GET /reports/by-category?user_id=X`

### CP6: CRON / Notifications
- Scheduled weekly/monthly spending summaries
- Delivery via email or webhook (configurable)

## Database Schema
See `internal/db/schema.sql`.
Tables: `users`, `sessions`, `groups`, `group_members`, `expenses`, `expense_splits`.

## Current Routes
```
GET  /ping                              (public)
POST /users/register                    (public)
POST /auth/login                        (public)
POST /auth/logout                       (protected)
GET  /users        GET  /users/{id}     (protected)
POST /expenses     GET  /expenses       (protected)
GET  /expenses/{id}                     (protected)
PUT  /expenses/{id}                     (protected)
DELETE /expenses/{id}                   (protected)
```

Protected routes require: `Authorization: Bearer <token>`
