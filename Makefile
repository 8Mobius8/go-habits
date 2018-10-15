BUILD_VERSION 	?= $(shell git describe --tags)
SERVER 					?= http://localhost:3000/api
INTEGRATION_ENV	?= BUILD_VERSION=${BUILD_VERSION} SERVER=${SERVER}
LDFLAGS 				?= -ldflags "-X main.version=${BUILD_VERSION}"

.PHONY: clean

dep:
	dep ensure

dep-clean:
	rm -rf ./vendor

test:	install test-docker-start test-unit test-integration

test-unit:
	go test -v -coverprofile=c.out ./api/... ./cmd/...

test-integration: install
	${INTEGRATION_ENV} \
	go test -v ./integration/...

test-watch: install test-docker-start
	${INTEGRATION_ENV} \
	ginkgo watch ./...

test-images:
	docker build -t registry.gitlab.com/8mobius8/go-habits/api -f integration/Dockerfile-habitica .
	docker push registry.gitlab.com/8mobius8/go-habits/api 

test-docker:
	docker-compose build
	docker-compose run integration

test-docker-start:
	docker-compose up -d habitica

test-clean:
	docker-compose down -v -t 0
	rm c.out

build:
	go build ${LDFLAGS}

install:
	go install ${LDFLAGS}

clean: test-clean dep-clean
	go clean
