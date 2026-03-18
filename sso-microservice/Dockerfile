FROM golang:1.25.3-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct
ENV GONOSUMDB=*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /app/bin/app ./cmd/sso/main.go

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /app/bin/migrator ./cmd/migrator/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/bin/migrator /app/migrator
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

COPY Docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]