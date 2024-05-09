include ../common.mk

.PHONY: ci
ci: gen gci test vulnerabilities_lookup lint

deps:
	go mod tidy
	go mod verify

.PHONY: build
build:
	go build -o ./app.exe ./cmd/service

.PHONY: test
test:
	go test -v ./...
