include ../common.mk

.PHONY: all
all: deps gen_api gen_common gci ci build

.PHONY: ci
ci: test vulnerabilities_lookup lint

deps:
	go mod tidy
	go mod verify

.PHONY: build
build:
	go build -o ./svc ./cmd/service

.PHONY: test
test:
	go test -v ./...

.PHONY: gen_api
gen_api:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-config internal/openapi/config.yml internal/openapi/openapi.yml
