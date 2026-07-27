MODULES := serverbase server-test shared-modules/audit shared-modules/booking shared-modules/calendar \
	shared-modules/documents shared-modules/invoice shared-modules/invoice_number \
	shared-modules/organization shared-modules/pdf shared-modules/saas-base \
	shared-modules/settings shared-modules/static

.PHONY: all build-all root-build staticcheck test run-server-test clean help

all: build-all

help:
	@printf "Available targets:\n"
	@printf "  build-all       Build each module listed in go.work\n"
	@printf "  root-build      Run 'go build ./...' from repository root\n"
	@printf "  staticcheck     Run staticcheck for each module\n"
	@printf "  test            Run 'go test ./...' for each module\n"
	@printf "  run-server-test Build and run server-test binary\n"
	@printf "  clean           Remove local build artifacts\n"

build-all:
	@set -e; \
	for d in $(MODULES); do \
			echo "=== Building $$d ==="; \
			(cd "$$d" && go mod tidy && go build ./...) || { echo "BUILD FAILED: $$d"; exit 1; }; \
		done; \
	echo "ALL BUILDS SUCCEEDED"

root-build:
	@$(MAKE) build-all

staticcheck:
	@GOPATH=$(shell go env GOPATH) && SC=$$GOPATH/bin/staticcheck || SC=staticcheck; \
	for d in $(MODULES); do \
		echo "--- staticcheck $$d ---"; \
		(cd "$$d" && $$SC ./...) || true; \
	done

test:
	@set -e; \
	for d in $(MODULES); do \
		echo "=== Testing $$d ==="; \
		(cd "$$d" && go test ./...) || { echo "TEST FAILED: $$d"; exit 1; }; \
	done; \
	echo "ALL TESTS SUCCEEDED"

unit:
	@set -e; \
	export MOCK_EMAIL=true; \
	for d in $(MODULES); do \
		echo "=== Unit testing $$d (short) ==="; \
		(cd "$$d" && go test -short ./...) || { echo "UNIT TEST FAILED: $$d"; exit 1; }; \
	done; \
	echo "ALL UNIT TESTS SUCCEEDED"

run-server-test:
	@echo "Building and running server-test"; \
	(cd server-test && go build -o server-test-bin . && ./server-test-bin)

integration:
	@set -e; \
	for d in $(MODULES); do \
		echo "=== Integration testing $$d ==="; \
		(cd "$$d" && go test ./...) || { echo "INTEGRATION TEST FAILED: $$d"; exit 1; }; \
	done; \
	echo "=== Running HURL integration tests (CI runner) ==="; \
	./server-test/scripts/run-hurl-ci.sh || { echo "HURL TESTS FAILED"; exit 1; }; \
	echo "ALL INTEGRATION TESTS SUCCEEDED"

clean:
	@echo "Cleaning build artifacts..."; \
	-rm -f server-test/server-test-bin
