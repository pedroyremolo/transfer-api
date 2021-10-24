TERM=xterm-256color
CLICOLOR_FORCE=true

.PHONY: format
format:
	@echo "### Formatting project... ###"
ifeq (, $(shell command -v goimports 2> /dev/null))
	@echo "Installing goimports..."
	go get golang.org/x/tools/cmd/goimports
endif
	@goimports -l -w ./
	@echo "\n"

.PHONY: run-lint
run-lint: format
	@echo "### Running linter... ###"
ifeq (, $(shell command -v golangci-lint 2> /dev/null))
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1
endif
	@golangci-lint run
	@echo "\n"

.PHONY: test
test:
	@echo "### Running tests... ###"
ifeq (, $(shell command -v richgo 2> /dev/null))
	@echo "richgo wasn't found, then we are installing"
	go get -u github.com/kyoh86/richgo
endif
	richgo test ./...
	@echo "\n"

.PHONY: build
build:
	docker build -t pedroyremolo/transfer-api .

.PHONY: pre-commit
pre-commit: run-lint test
	go mod tidy
