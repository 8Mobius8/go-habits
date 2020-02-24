APP := go-habits
VERSION := $(shell git describe --tags)
GIT_HASH := $(shell git rev-parse --short HEAD)

LDFLAGS ?= -X main.version=${VERSION}
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build \
	-ldflags "${LDFLAGS}" \
	-o release/${APP}-${VERSION}-$(os)

.PHONY: release
release: lint windows linux darwin

.PHONY: dep
dep:
	go get -u golang.org/x/lint/golint
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -v -u -d ./...

# Install go depedancies using go mods
.PHONY: vendor
vendor:
	go mod vendor -v

.PHONY: install
install: vendor
	go install -ldflags "${LDFLAGS}"

# Will build artifacts: executable, docker images.
.PHONY: build build-images
build: release build-images

# Builds docker images for go-habits to be use in docker-compose
# also can push images to registry. You will need to run this before
# trying to use the docker-compose file.
build-images:
	docker build -f _dockerfiles/go-habits-tester \
		--tag registry.gitlab.com/8mobius8/go-habits/tester:latest .
ifdef api
	docker build -f _dockerfiles/habitica \
		--tag registry.gitlab.com/8mobius8/go-habits/habitica:latest .
endif
ifdef push
	docker push registry.gitlab.com/8mobius8/go-habits/tester:latest
ifdef api
	docker push registry.gitlab.com/8mobius8/go-habits/habitica:latest
endif
endif

PKG_LIST  ?= $(shell go list ./...)
.PHONY: lint
lint:
	@golint ${PKG_LIST}

# Runs full test suite
.PHONY: test test-unit test-unit-randomly test-integration
test:	test-unit test-integration

test-unit:
	@rm -f c.out
	ginkgo \
	-randomizeAllSpecs -randomizeSuites \
	-failOnPending \
	-trace \
	-race \
	-progress \
	-cover -outputdir=. -coverprofile=c.out \
	-skipPackage integration \
	-r \
	.

# Use this to run all unit tests randomly many times
test-unit-randomly: 
	ginkgo \
	-randomizeAllSpecs -randomizeSuites \
	-skipPackage integration \
	-r \
	-untilItFails -p \
	.

SERVER          ?= http://localhost:3000/api
INTEGRATION_ENV ?= BUILD_VERSION=${VERSION} SERVER=${SERVER}
test-integration: install
	${INTEGRATION_ENV} ./integration/wait-for-habitica-api.sh
	${INTEGRATION_ENV} ginkgo -r --trace --progress ./integration

# Uses ginkgo's output which doesn't respect `go tool cover` format.
# `ginkgo` returns mutliple `mode: atomic` lines.
.PHONY: coverage coverage-html code-climate cc-before cc-after
coverage: test-unit
	@cat c.out | grep --max-count=1 ^mode: > c.cov
	@cat c.out | grep -v ^mode: >> c.cov
	@mv c.cov c.out
	go tool cover -func=c.out

coverage-html: test-unit
	@cat c.out | grep --max-count=1 ^mode: > c.cov
	@cat c.out | grep -v ^mode: >> c.cov
	@mv c.cov c.out
	go tool cover -html=c.out -o coverage.html

# Ordered goal to ensure cc-test-reporter is called before tests are run.
# Ment to be run using docker or in CI.
code-climate: | cc-before coverage cc-after
cc-before:
	cc-test-reporter before-build
cc-after:
	cc-test-reporter after-build

# Useful for development on this project. Spin up habitica API server, lint 
# files, uses `gofmt` as linter, continuous run tests when editing go files.
.PHONY: dev dev-format
dev: dev-format install
	${INTEGRATION_ENV} ginkgo watch ./...

# Lints and applies lints to code.
dev-format:
	gofmt -d -w main.go ./api ./cmd ./integration ./log

.PHONY: clean clean-vendor clean-release
clean:
	go clean
	rm -rf **/*c.out

clean-vendor:
	rm -rf ./vendor

clean-release:
	rm -rf ./release

