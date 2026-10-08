# Interviewtopia Completion Service

## Pre-interview setup

This repository is being shared before the live interview so you can confirm that the starter builds and runs in your development environment.

The exercise brief is intentionally not included here. It will be provided at the start of the interview, and no feature implementation is expected in advance.

## Getting started

Requirements:

- Go 1.23 or newer
- Make, optionally

Run the starter:

```sh
go run ./cmd/server
```

The server listens on `http://127.0.0.1:8080` by default. Override it with `HTTP_ADDR`.

Check the starter:

```sh
go test ./...
go test -race ./...
go vet ./...
```

Or run all checks with:

```sh
make check
```

The starter health endpoint is:

```sh
curl -i http://127.0.0.1:8080/healthz
```

If these commands succeed, you are ready for the interview. Please raise any environment or toolchain problem before the session so that interview time is not spent on setup.
