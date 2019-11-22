SHELL = bash

test:
	GO_REDIS_PORT=:9091 go test ./...

github_test:
	GO_REDIS_PORT=:6379 go test ./...

local_test:
	GO_REDIS_PORT=:9090 go test ./... -count=1

lint:
	@go fmt ./... && go vet ./...
