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