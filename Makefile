
help: ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

install: ## install all go dependencies
	go get \
		github.com/bufbuild/buf/cmd/buf \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		google.golang.org/protobuf/cmd/protoc-gen-go 
		
generate: ## generate grpc proto for golang
	buf generate

test: ## run all layers tests
	@make test.unit
	@make test.integration

test.unit: ## run unit tests
	go test ./internal/...

test.integration: ## run integration tests
	@make test.integration.up
	@make test.integration.start
	@make test.integration.down

test.integration.up: ## build integration test environment
	docker-compose -f deployments/test/docker-compose.yaml up -d

test.integration.start: ## run integration test
	go test ./tests/integration/...
	
test.integration.down: ## shutting down integration environment
	docker-compose -f deployments/test/docker-compose.yaml down --volumes --rmi local

run: ## run app
	go run ./cmd/app/main.go

mock: ## run mock
	go run ./cmd/mock/main.go

dev.up: ## run dev environment
	docker-compose -f deployments/dev/docker-compose.yaml up -d

dev.down: ## shutting down dev environment
	docker-compose -f deployments/dev/docker-compose.yaml down --volumes --rmi local