deps:
	dep ensure

test:
	go test -coverprofile=c.out ./ ./api/... ./cmd/...
	docker-compose build integration
	docker-compose run --rm integration

build:
	go build

install:
	go install