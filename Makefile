# Convenience wrappers for Mac/Linux. On Windows use: go run scripts/<name>.go
.PHONY: dev debug build start

dev:
	go run scripts/dev.go

debug:
	go run scripts/debug.go

build:
	go run scripts/build.go

start:
	go run scripts/start.go
