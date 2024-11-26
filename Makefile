SHELL:=/bin/sh
GCI_LINT:=v1.62.0


help:      ### Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help


coverage:  ### Collect code coverage
	go test -coverprofile=coverage.out ./...;
	go tool cover -html=coverage.out;
.PHONY: coverage


lint:      ### Run the go formatter and golangci-lint
	go fmt ./...;
	docker run --rm -t -v $(PWD):/mnt -w /mnt golangci/golangci-lint:$(GCI_LINT) \
		golangci-lint run --fix -v;
.PHONY: lint


lint-latest:  ### Run latest golangci-lint with ALL linters enabled
	docker run --rm -t -v $(PWD):/mnt -w /mnt golangci/golangci-lint:latest \
		golangci-lint run --enable-all --no-config -v;
.PHONY: lint-latest


test:      ### Run unit tests
	go test ./...;
.PHONY: test


test-race: ### Run unit tests with the race detector enabled
	go test -race ./...;
.PHONY: test-race