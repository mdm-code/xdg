GO=go
GOFLAGS=-race
DEV_BIN=bin
COV_PROFILE=cp.out

.DEFAULT_GOAL := test

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...

.PHONY: test
test: vet
	$(GO) clean -testcache
	$(GO) test ./... -v

.PHONY: install
install: test
	$(GO) install ./...

.PHONY: build
build:
	$(GO) build github.com/mdm-code/xdg/...

.PHONY: cover
cover:
	$(GO) test -coverprofile=$(COV_PROFILE) ./...
	$(GO) tool cover -html=$(COV_PROFILE)

.PHONY: clean
clean:
	$(GO) clean github.com/mdm-code/xdg/...
	$(GO) clean -testcache
	rm -f $(COV_PROFILE)
