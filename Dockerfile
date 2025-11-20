ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN mkdir -p bin && go build -o bin/api ./cmd/api

FROM debian:bookworm

COPY --from=builder /usr/src/app/bin/api /usr/local/bin/api
CMD ["/usr/local/bin/api"]
