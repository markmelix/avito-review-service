FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR ${GOPATH}/loadtest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /loadtest ./loadtest.go \
    && go clean -cache -modcache


FROM alpine:3.22

COPY --from=builder /loadtest /loadtest

WORKDIR /

ENTRYPOINT ["/loadtest"]
