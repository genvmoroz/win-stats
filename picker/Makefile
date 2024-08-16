include ../common.mk

.PHONY: all
all: tidy gen_api gen_common gci ci build

.PHONY: ci
ci: test vulnerabilities_lookup lint

.PHONY: build
build:
	go build -o ./app.exe ./cmd/service

.PHONY: test
test:
	go test -v ./...

.PHONY: gen_api
gen_api:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		-config api/openapi.cfg.yml api/openapi.yml
