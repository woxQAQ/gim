.phony: module
module: 
	go mod tidy -compat=1.22
	go mod verify

GO_FILES := $(shell git ls-files | grep "\.go$$")

.phony: imports
imports:
	goimports -local github.com/woxQAQ/gim -w $(GO_FILES)

.phony: fmt
fmt:
	go fmt ./... 

.phony: test
test: fmt module
	@ echo "\033[1;32mtest gim...\033[0m"
	go test ./... -coverprofile cover.out

.phony: lint
lint: module
	@ echo -e "\033[1;32mgolangci-lint...\033[0m"
	golangci-lint run

.phony: install-hooks
install-hooks:
	@ echo -e "\033[1;32mInstalling git hooks...\033[0m"
	chmod +x .git/hooks/pre-commit