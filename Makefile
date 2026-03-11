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



recommendation:
	cd recommedation-system-microservice && go run cmd/rec-system/main.go --config=./config/local.yaml

recommendation-migrator:
	cd recommedation-system-microservice && go run cmd/migrator/main.go --migrations-path=./migrations

notifications:
	cd notification-microservice && go run cmd/app/main.go --config=./config/local.yaml

apigateway:
	cd api-gateway && go run cmd/gateway/main.go