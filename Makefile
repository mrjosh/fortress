export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOBIN?=$(BIN)
export GO=$(shell which go)
export BUILD=cd $(ROOT) && $(GO) install -v -ldflags "-s"
export CGO_ENABLED=1

# Linter configurations
export LINTER=$(GOBIN)/golangci-lint
export LINTERCMD=run --no-config -v \
	--print-linter-name \
	--skip-files ".*.gen.go" \
	--skip-files ".*_test.go" \
	--sort-results \
	--disable-all \
	--enable=structcheck \
	--enable=deadcode \
	--enable=gocyclo \
	--enable=ineffassign \
	--enable=revive \
	--enable=goimports \
	--enable=errcheck \
	--enable=varcheck \
	--enable=goconst \
	--enable=megacheck \
	--enable=misspell \
	--enable=unused \
	--enable=typecheck \
	--enable=staticcheck \
	--enable=govet \
	--enable=gosimple

all:
	@$(BUILD) ./...

# lint runs vet plus a number of other checkers, it is more comprehensive, but louder
lint:
	@LINTER_BIN=$$(command -v $(LINTER)) || { echo "golangci-lint command not found! Installing..." && $(MAKE) install-metalinter; };
	@$(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| xargs $(LINTER) $(LINTERCMD) ./...; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Lint found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# for ci jobs, runs lint against the changed packages in the commit
ci-lint:
	@$(LINTER) $(LINTERCMD) --deadline 10m --new-from-rev=HEAD~ ./...

test:
	@$(GO) test -v -race ./...

# Check if golangci-lint not exists, then install it
install-metalinter:
	@$(GO) get -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
