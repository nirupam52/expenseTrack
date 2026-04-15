---
name: code-standards
description: Enforces project-specific code conventions and best practices
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: development
---

## What I do

Enforce consistent code patterns across the expenseTrack project:

### Response Pattern (REQUIRED)
Use these functions from `internal/response/response.go`:
- `response.WriteSuccess(w, status, data)` - for successful responses
- `response.WriteError(w, status, message)` - for errors
- `response.WriteList(w, data)` - for list responses

Response struct:
```go
type Response[T any] struct {
    Success bool `json:"success"`
    Data    T    `json:"data"`
}
```

### Handler Pattern (REQUIRED)
- Method validation: Check `r.Method` at start of every handler
- Input validation: Validate all required fields before DB operations
- Context: Use `r.Context()` instead of `context.Background()` for request-scoped operations
- Error handling: Return early with appropriate HTTP status codes

```go
func (h *Handler) HandlerName(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    // validation, business logic
}
```

### Repository Pattern (REQUIRED)
- Always use `context.Context` as first parameter
- Wrap errors with `fmt.Errorf("action: %w", err)` pattern
- Check `sql.ErrNoRows` for not found errors
- Use `defer` for resource cleanup

```go
func (r *Repo) GetByID(ctx context.Context, id int64) (*Model, error) {
    // query...
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("entity not found")
    }
    if err != nil {
        return nil, fmt.Errorf("get by id: %w", err)
    }
    return &model, nil
}
```

### HTTP Status Codes (REQUIRED)
- 200 - OK / Success
- 201 - Created
- 400 - Bad Request (validation errors)
- 401 - Unauthorized
- 403 - Forbidden
- 404 - Not Found
- 409 - Conflict (duplicate resource)
- 500 - Internal Server Error

### Naming Conventions
- Handlers: `UserHandler`, `ExpenseHandler` - singular, CamelCase
- Repositories: `UserRepository`, `ExpenseRepository` - singular, CamelCase
- Routes: RESTful - `/expenses`, `/users`, `/expenses/{id}`
- Query params: `user_id`, `group_id` - snake_case

## When to use me

Before committing code or creating a PR, verify:
1. All handlers use response helper functions
2. Context is properly passed through the call stack
3. Errors are wrapped with context
4. Method validation is present
5. Input validation is complete

## How to invoke

```
skill({ name: "code-standards" })
```

This loads the full standards for the agent to verify code against.