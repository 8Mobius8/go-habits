CI_COMMIT_SHA   ?= $(shell git rev-parse HEAD)
BUILD_VERSION   ?= $(shell git describe --tags)
PKG_LIST 				?= $(shell go list ./... | grep -v /vendor/)
SERVER          ?= http://localhost:3000/api
INTEGRATION_ENV ?= BUILD_VERSION=${BUILD_VERSION} SERVER=${SERVER}
LDFLAGS         ?= -ldflags "-X main.version=${BUILD_VERSION}"


# Setup to show what needs to be installed
.PHONY: setup setup-go setup-docker setup-dep setup-test
setup:	setup-go setup-docker setup-dep setup-test

setup-go:
ifeq (, $(shell which go))
	@echo "Please install a verison of golang."
	@echo "https://golang.org/doc/install"
endif

setup-docker:
ifeq (, $(shell which docker-compose))
	@echo "Please install a verison docker and docker-compose."
	@echo "https://docs.docker.com/install"
endif

setup-dep:
ifeq (, $(shell which dep))
	@echo "Please install dep"
	@echo "https://golang.github.io/dep/docs/installation.html"
endif

setup-test:
ifeq (, $(shell which ginkgo))
	@echo "No ginkgo in PATH. Will attempt to install"
	go get github.com/onsi/ginkgo/ginkgo
endif

# Install all depedancies for building and testings
# Only `dep`
.PHONY: dep 
dep: setup
	dep ensure -v

# Will build artifacts: executable, docker images.
.PHONY: build go-build build-images
build: go-build build-images

go-build:	
	go build ${LDFLAGS}

build-images:
	CI_COMMIT_SHA=${CI_COMMIT_SHA} docker-compose build tests
ifdef push
	CI_COMMIT_SHA=${CI_COMMIT_SHA} docker-compose push tests
endif
ifdef api
	docker-compose build habitica
ifdef push
	docker-compose push habitica
endif
endif

.PHONY: install
install:
	go install ${LDFLAGS}

# Run full test suite
.PHONY: test test-unit test-integration
test:	install test-unit test-integration

test-unit:
	rm c.out
ifdef ccreporter
	./cc-test-reporter before-build
endif
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
ifdef ccreporter
	./cc-test-reporter after-build
endif

test-integration: install docker-services
	${INTEGRATION_ENV} ./integration/wait-for-habitica-api.sh
	${INTEGRATION_ENV} ginkgo -r --trace --progress ./integration

.PHONY: docker-services
docker-services:
	docker-compose up -d habitica


# Uses ginkgo's output which doesn't respect `go tool cover` format.
# `ginkgo` returns mutliple `mode: atomic` lines.
.PHONY: coverage
coverage: test-unit
	cat c.out | grep --max-count=1 ^mode: > c.cov
	cat c.out | grep -v ^mode: >> c.cov 
	go tool cover -func=c.cov

# Useful for development on this project. Spin up habitica API server, lint 
# files, uses `gofmt` as linter, continuous run tests when editing go files.
.PHONY: dev dev-format
dev: dev-format install docker-services
	${INTEGRATION_ENV} ginkgo watch ./...

# Lints and applies lints to code.
dev-format:
	gofmt -d -w main.go ./api ./cmd ./integration

.PHONY: clean clean-dep
clean:
	go clean
	docker-compose down -v -t 0
	rm -rf **/*c.out

clean-dep:
	rm -rf ./vendor

