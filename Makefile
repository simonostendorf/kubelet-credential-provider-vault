GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

default: build

tidy:
	@go mod tidy

deps: dependencies
dependencies: tidy
	@go mod download

build: dependencies
	@GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o ./bin/kubelet-credential-provider-vault_$(GOOS)_$(GOARCH) ./main.go

test:
	@go test -v -coverprofile=coverage.out -covermode=atomic -short ./...
	@go tool cover -func=coverage.out

e2e: build e2e-init
	@cd ./tests/e2e/ && bash e2e_run.sh

e2e-init:
	@cd ./tests/e2e/ && bash e2e_start.sh

e2e-cleanup:
	@cd ./tests/e2e/ && bash e2e_stop.sh

lint:
	@golangci-lint run

lint-security:
	@gosec -tests ./...
