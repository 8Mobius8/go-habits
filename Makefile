CI_COMMIT_SHA		?= $(shell git rev-parse HEAD)
BUILD_VERSION 	?= $(shell git describe --tags)
SERVER 					?= http://localhost:3000/api
INTEGRATION_ENV	?= BUILD_VERSION=${BUILD_VERSION} SERVER=${SERVER}
LDFLAGS 				?= -ldflags "-X main.version=${BUILD_VERSION}"

.PHONY: clean

dep:
	dep ensure

dep-clean:
	rm -rf ./vendor

build:
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

test:	install test-unit test-integration

test-unit:
ifdef ccreporter
	./cc-test-reporter before-build
endif
	go test -v -coverprofile=c.out ./api/... ./cmd/...
ifdef ccreporter
	./cc-test-reporter after-build
endif

test-integration: install
	${INTEGRATION_ENV} ./integration/wait-for-habitica-api.sh
	${INTEGRATION_ENV} go test -v ./integration/...

test-clean:
	docker-compose down -v -t 0
	rm c.out

install:
	go install ${LDFLAGS}

dev: install
	docker-compose up -d habitica
	${INTEGRATION_ENV} ginkgo watch ./...

clean: test-clean dep-clean
	go clean