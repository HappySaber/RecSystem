FROM golang:1.25.3-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct
ENV GONOSUMDB=*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /app/bin/app ./cmd/gateway/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/config /app/config

EXPOSE 8080

CMD ["/app/app", "--config=/app/config/local.yaml"]