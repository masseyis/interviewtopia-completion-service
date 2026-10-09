# Interviewtopia Completion Service

## Pre-interview setup

This repository is being shared before the live interview so you can confirm that the starter builds and runs in your development environment.

The exercise brief is intentionally not included here. It will be provided at the start of the interview, and no feature implementation is expected in advance.

## Getting started

Requirements:

- Go 1.23 or newer
- Make, optionally
- No external database. The starter uses embedded `bbolt` storage.

Run the starter:

```sh
go run ./cmd/server
```

The server listens on `http://127.0.0.1:8080` by default. Runtime values are
held in `internal/appconfig.Config` so tests can inject them directly:

| Environment variable | Default |
|---|---|
| `HTTP_ADDR` | `127.0.0.1:8080` |
| `REGISTRY_BASE_URL` | Empty until supplied during the exercise |
| `TRUST_STORE_PATH` | `config/trusted-issuers.json` |
| `COMPLETION_DB_PATH` | `completion.db` |

At startup the executable loads that file and injects the resulting trust
store, registry base URL, embedded database and HTTP client into
`httpapi.Dependencies`. Tests can construct those dependencies directly
without relying on environment variables. Use a path under `t.TempDir()` for
test databases.

`internal/storage.Open` handles opening and closing the embedded database. It
does not choose buckets, records, indexes or transaction boundaries; those are
application-design decisions. `TestDatabaseSurvivesCloseAndReopen` shows the
basic persistence mechanism without implementing the exercise domain.

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

## Included test support

The starter includes mechanics that you may use during the exercise. They do
not define the business rules:

- `internal/testkit.Sign` creates a correctly signed envelope with one line of
  test code using the synthetic issuers in the supplied trust store.
- `internal/testkit.NewRegistryStub` starts a local registry that can return a
  chosen status or body, delay its response, or wait until a request is
  cancelled. This supports deterministic 404, 500, malformed-JSON and timeout
  tests.
- `examples/evidence` contains signed sample messages. Some additional samples
  deliberately use neutral filenames.
- `TestHealthOverHTTP` is a complete HTTP-level test showing a local server,
  client request and JSON response assertion.

For example:

```go
envelope := testkit.Sign(t, "bank-7", payload)

registry := testkit.NewRegistryStub(t, testkit.RegistryResponse{
    StatusCode: http.StatusInternalServerError,
})
```

All keys and signatures in this repository are synthetic test material. Never
reuse them outside this assessment.
