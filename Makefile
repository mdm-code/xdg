GO=go
GOFLAGS=-race
DEV_BIN=bin
COV_PROFILE=cp.out

.DEFAULT_GOAL := build

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet: fmt
	$(GO) vet ./...

.PHONY: lint
lint: vet
	golint -set_exit_status=1 ./...

.PHONY: test
test: lint
	$(GO) clean -testcache
	$(GO) test ./... -v

.PHONY: install
install: test
	$(GO) install ./...

.PHONY: build
build: test
	$(GO) build github.com/mdm-code/xdg/...

.PHONY: cover
cover:
	$(GO) test -coverprofile=$(COV_PROFILE) -covermode=atomic ./...
	$(GO) tool cover -html=$(COV_PROFILE)

.PHONY: clean
clean:
	$(GO) clean github.com/mdm-code/xdg/...
	$(GO) clean -testcache
	rm -f $(COV_PROFILE)
