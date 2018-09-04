.PHONY: clean
all: deps test build install
deps:
	dep ensure

test:	test-unit test-integration

test-unit:
	go test -coverprofile=c.out ./api/... ./cmd/...

test-integration:
	docker-compose build
	docker-compose run --rm integration

build:
	go build

install:
	go install

clean:
	go clean