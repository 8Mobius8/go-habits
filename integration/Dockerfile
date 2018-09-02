FROM golang:1.10 AS builder

RUN mkdir -p $GOPATH/src/github.com/8Mobius8/go-habits
WORKDIR $GOPATH/src/github.com/8Mobius8/go-habits

COPY . .

RUN go install

ENTRYPOINT ./integration/wait-for-habitica-api.sh go test -v ./integration/...