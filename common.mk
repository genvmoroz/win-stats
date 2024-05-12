.PHONY: gen_common
gen_common:
	go install golang.org/x/tools/cmd/stringer@latest
	go install go.uber.org/mock/mockgen@latest
	go generate ./...

vulnerabilities_lookup:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck -test ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
		run --allow-parallel-runners -c ../.golangci.yml

.PHONY: gci
gci:
	go run github.com/luw2007/gci@latest \
		write . --skip-generated

upgrade_deps:
	go run github.com/genvmoroz/tolatest@latest \
 		./go.mod
