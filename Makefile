.PHONY: test race vet fmt-check check

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Run gofmt on:"; gofmt -l .; exit 1)

check: fmt-check vet test race

