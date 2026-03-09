# expenseTrack

A personal expense tracker written in Go. Track your own expenses and split them with friends using groups.

## Requirements

- Go 1.25+
- [Air](https://github.com/air-verse/air) (for hot reload, optional)

## Setup

```bash
git clone https://github.com/nirupam52/expenseTrack
cd expenseTrack
go mod download
```

## Run

**Standard:**
```bash
go run .
```

**With hot reload (Air):**
```bash
go run scripts/dev.go
```

The server starts on port `41605` by default. Set the `APP_PORT` environment variable to change it.

## Verify the setup

```bash
curl http://localhost:41605/ping
```

Expected response:
```json
{"success":true,"data":"we good bro :)"}
```
