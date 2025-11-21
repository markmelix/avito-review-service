FROM golang:1.24.9-alpine3.22 AS builder

WORKDIR ${GOPATH}/avito-review

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /server ./cmd/server \
    && go clean -cache -modcache


FROM alpine:3.22

COPY --from=builder /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
