binary:
	go build -o govarpkg ./cmd/govarpkg

lint:
	golangci-lint run cmd/... pkg/...

test:
	go test ./pkg/...