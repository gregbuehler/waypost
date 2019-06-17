SRC_DIR = ./
GO_PKGS	:= $(shell go list ./... | grep -v "/vendor/")
GO_FILES := ./
CGO_ENABLED = 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= $(shell uname -m | sed -e 's/x86_64/amd64/')
VALIDATE_DEPS    = golang.org/x/lint/golint
TEST_DEPS        = github.com/axw/gocov/gocov github.com/AlekSi/gocov-xml


all: build

build: clean validate-deps validate compile test

clean:
	@echo "=== [ clean ]: purging binaries and coverage file"
	@rm -rfv bin coverage.xml

validate-deps:
	@echo "=== [ validate-deps ]: installing validation dependencies"
	@go get -v $(VALIDATE_DEPS)

validate-source:
ifeq ($(strip $(GO_FILES)),)
	@echo "=== [ validate ]: no Go files found. Skipping validation."
else
	@printf "=== [ validate ]: running gofmt... "
	# `gofmt` expects files instead of packages. `go fmt` works with
	# packages, but forces -l -w flags.
	@OUTPUT="$(shell gofmt -l $(GO_FILES))" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Incorrect syntax in the following files:" ;\
		echo "$$OUTPUT" ;\
		exit 1 ;\
	fi
	@printf "=== [ validate ]: running golint... "
	@OUTPUT="$(shell golint $(SRC_DIR)...)" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Issues found:" ;\
		echo "$$OUTPUT" ;\
		exit 1 ;\
	fi
	@printf "=== [ validate ]: running go vet... "
	@OUTPUT="$(shell go vet $(SRC_DIR)...)" ;\
	if [ -z "$$OUTPUT" ]; then \
		echo "passed." ;\
	else \
		echo "failed. Issues found:" ;\
		echo "$$OUTPUT" ;\
		exit 1;\
	fi
endif

validate: validate-deps validate-source

compile-deps:
	@echo "=== [ compile-deps ]: installing build dependencies"
	@go get -v -d -t ./...

bin/waypost:
	@echo "=== [ compile ]: building waypost"
	@env CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo -v -o bin/waypost $(GO_FILES)

compile: compile-deps bin/waypost

test-deps:
	@echo "=== [ test-deps ]: installing testing dependencies"
	@go get -v $(TEST_DEPS)

test-source:
	@echo "=== [ test ]: ..."
	@gocov test ./... | gocov-xml > coverage.xml


test: test-deps test-source


# TODO: add integration goals

install: bin/waypost
	@echo "=== [ install ]: installing bin/waypost"
	@install -D --mode=755 --owner=root $(ROOT)bin/waypost /usr/local/bin/waypost
	@install -D --mode=644 --owner=root $(ROOT)/etc/waypost/waypost.toml /etc/waypost/waypost.toml

.PHONY: all build clean validate-deps validate-source validate compile-deps compile install
