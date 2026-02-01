FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/bin/app ./cmd/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/bin/migrator ./cmd/migrator


FROM alpine:3.19

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/bin/migrator /app/migrator

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV GRPC_PORT=50051
EXPOSE 50051

ENTRYPOINT ["/entrypoint.sh"]
