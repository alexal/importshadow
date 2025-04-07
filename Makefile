binary:
	go build -o importshadow ./cmd/importshadow
.PHONY: binary

lint:
	golangci-lint run cmd/... pkg/...
.PHONY: lint

test:
	go test ./pkg/...
.PHONY: test

format:
	@go mod tidy

	@gofmt -w `find . -type f -name '*.go'`
.PHONY: format