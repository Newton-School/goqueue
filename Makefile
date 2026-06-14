.PHONY: fmt test vet verify

GOFILES := $(shell find . -name '*.go' -not -path './.git/*')

fmt:
	gofmt -w $(GOFILES)

test:
	go test ./...

vet:
	go vet ./...

verify:
	@test -z "$$(gofmt -l $(GOFILES))"
	go vet ./...
	go test ./...
