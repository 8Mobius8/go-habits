BUILD_VERSION   ?= $(shell git describe --tags)
PKG_LIST 				?= $(shell go list ./... | grep -v /vendor/)
SERVER          ?= http://localhost:3000/api
INTEGRATION_ENV ?= BUILD_VERSION=${BUILD_VERSION} SERVER=${SERVER}
LDFLAGS         ?= -ldflags "-X main.version=${BUILD_VERSION}"

# All build so that when `make` is called with
# just build binary
all: dep go-build

# Install go depedancies using dep
# Only `dep`
.PHONY: dep dep-vendor-only
dep:
	dep ensure -v

dep-vendor-only:
	dep ensure -v -vendor-only

# Will build artifacts: executable, docker images.
.PHONY: build go-build build-images
build: go-build build-images

go-build:	
	go build ${LDFLAGS}

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

.PHONY: install
install:
	go install ${LDFLAGS}

# Runs full test suite
.PHONY: test test-unit test-integration
test:	test-unit test-integration

test-unit: dep
	rm -f c.out
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

test-integration: dep install
	${INTEGRATION_ENV} ./integration/wait-for-habitica-api.sh
	${INTEGRATION_ENV} ginkgo -r --trace --progress ./integration

# Uses ginkgo's output which doesn't respect `go tool cover` format.
# `ginkgo` returns mutliple `mode: atomic` lines.
.PHONY: coverage coverage-html code-climate cc-before cc-after
coverage: test-unit
	@cat c.out | grep --max-count=1 ^mode: > c.cov
	@cat c.out | grep -v ^mode: >> c.cov
	@mv c.cov c.out
	@go tool cover -func=c.out

coverage-html: test-unit
	@cat c.out | grep --max-count=1 ^mode: > c.cov
	@cat c.out | grep -v ^mode: >> c.cov
	@mv c.cov c.out
	@go tool cover -html=c.out -o coverage.html

# Ordered goal to ensure cc-test-reporter is called before tests are run.
# Ment to be run using docker or in CI.
code-climate: | cc-before coverage cc-after
cc-before:
	cc-test-reporter before-build
cc-after:
	cc-test-reporter after-build

lint:
	@golint ${PKG_LIST}

# Useful for development on this project. Spin up habitica API server, lint 
# files, uses `gofmt` as linter, continuous run tests when editing go files.
.PHONY: dev dev-format
dev: dev-format install
	${INTEGRATION_ENV} ginkgo watch ./...

# Lints and applies lints to code.
dev-format:
	gofmt -d -w main.go ./api ./cmd ./integration

.PHONY: clean clean-dep
clean:
	go clean
	rm -rf **/*c.out

clean-dep:
	rm -rf ./vendor

