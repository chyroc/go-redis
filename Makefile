SHELL = bash

test:
	GO_REDIS_PORT=:9091 go test ./...

github_test:
	GO_REDIS_PORT=:6379 go test ./...

lint:
	@go fmt ./... && go vet ./...
