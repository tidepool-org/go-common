SHELL = /bin/sh

TOOLS_BIN = tools/bin
NPM_BIN = node_modules/.bin

OAPI_CODEGEN = $(TOOLS_BIN)/oapi-codegen
REDOCLY_CLI = $(NPM_BIN)/redocly

NPM_PKG_SPECS = \
	@redocly/cli@1.34.1

ifeq ($(CI),)
GO_BUILD_FLAGS =
GO_TEST_FLAGS =
else
GO_BUILD_FLAGS = -v
GO_TEST_FLAGS = -v
endif

.PHONY: build
build:
	go build $(GO_BUILD_FLAGS) ./...

.PHONY: test
test:
	go test $(GO_TEST_FLAGS) ./...

.PHONY: test-cover
test-cover:
	go test -coverprofile cover.out $(GO_TEST_FLAGS) ./...

.PHONY: generate
# Generates client api
generate: $(REDOCLY_CLI) $(OAPI_CODEGEN)
	$(REDOCLY_CLI) bundle ../TidepoolApi/reference/summary.v1.yaml -o ./spec/summary.v1.yaml
	$(OAPI_CODEGEN) -package=api -generate=types spec/summary.v1.yaml > clients/summary/types.go
	$(OAPI_CODEGEN) -package=api -generate=client spec/summary.v1.yaml > clients/summary/client.go
	cd clients/summary && go generate ./...

$(OAPI_CODEGEN):
	GOBIN=$(shell pwd)/$(TOOLS_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

$(REDOCLY_CLI): npm-tools

.PHONY: npm-tools
npm-tools:
# When using --no-save, any dependencies not included will be deleted, so one
# has to install all the packages all at the same time. But it saves us from
# having to muck with packages.json.
	npm install --no-save --local $(NPM_PKG_SPECS)

.PHONY: clean
clean:
	rm -rf node_modules tools

.PHONY: ci-generate
ci-generate: generate

.PHONY: ci-build
ci-build: build

.PHONY: ci-test
ci-test: test
