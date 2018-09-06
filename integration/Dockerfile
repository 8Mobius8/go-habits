FROM golang:1.10 AS builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR $GOPATH/src/github.com/8Mobius8/go-habits
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v -vendor-only

COPY . ./

CMD ./integration/wait-for-habitica-api.sh make test-integration