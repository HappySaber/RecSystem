.PHONY: generate sso sso-migrator sso-migrator-test sso-test

generate:
	buf generate


sso:
	cd sso-microservice && go run cmd/sso/main.go --config=./config/local.yaml

sso-migrator:
	cd sso-microservice && go run cmd/migrator/main.go --migrations-path=./migrations

sso-migrator-test:
	cd sso-microservice && go run cmd/migrator/main.go --migrations-path=./tests/migrations --test

sso-test:
	cd sso-microservice && go run cmd/migrator/main.go --migrations-path=./tests/migrations


catalog:
	cd catalog-system-microservice && go run cmd/catalog/main.go --config=./config/local.yaml

catalog-migrator:
	cd catalog-system-microservice && go run cmd/migrator/main.go --migrations-path=./migrations

catalog-importer:
	cd catalog-system-microservice && go run cmd/importer/main.go --config=./config/local.yaml

catalog-test-integration:
	cd catalog-system-microservice && go test ./tests/... -v -timeout 120s

catalog-test-unit:
	cd catalog-system-microservice && go test ./internal/services/... -v

catalog-test-all:
	cd catalog-system-microservice && go test ./... -v -timeout 120s


recommendation:
	cd recommedation-system-microservice && go run cmd/rec-system/main.go --config=./config/local.yaml

recommendation-migrator:
	cd recommedation-system-microservice && go run cmd/migrator/main.go --migrations-path=./migrations

notifications:
	cd notification-microservice && go run cmd/app/main.go --config=./config/local.yaml

notifications-test-unit-mail:
	cd notification-microservice && go test ./internal/services/... -v

notifications-test-unit-consumer:
	cd notification-microservice && go test ./internal/consumer/... -v

apigateway:
	cd api-gateway && go run cmd/gateway/main.go