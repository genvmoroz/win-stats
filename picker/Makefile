include ../common.mk

.PHONY: all
all: gen_api gen_common gci ci

.PHONY: ci
ci: test vulnerabilities_lookup lint

deps:
	go mod tidy
	go mod verify

.PHONY: build
build:
	go build -o ./app.exe ./cmd/service

.PHONY: test
test:
	go test -v ./...

.PHONY: gen_api
gen_api:
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-config api/openapi.cfg.yml api/openapi.yml
