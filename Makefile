.PHONY: fmt test vet race cover staticcheck vulncheck integration-test verify audit docs-install docs-build docs-serve

GOFILES := $(shell find . -name '*.go' -not -path './.git/*')

fmt:
	gofmt -w $(GOFILES)

test:
	go test ./...

vet:
	go vet ./...

race:
	go test -race ./...

cover:
	go test -cover ./...

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

integration-test:
	GOQUEUE_RUN_INTEGRATION_TESTS=true go test -count=1 ./redisbackend

verify:
	@test -z "$$(gofmt -l $(GOFILES))"
	go vet ./...
	go test ./...

audit: verify staticcheck vulncheck race

docs-install:
\tcd docs && npm run docs-install

docs-build:
\tcd docs && npm run docs-build

docs-serve:
\tcd docs && npm run docs-serve
