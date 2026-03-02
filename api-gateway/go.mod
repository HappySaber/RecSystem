module api-gateway

go 1.25.3

require (
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/gorilla/mux v1.8.1
	github.com/segmentio/kafka-go v0.4.50
	google.golang.org/protobuf v1.36.10
)

replace rec-system-microservice => ../recommedation-system-microservice
