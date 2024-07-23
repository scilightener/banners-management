.PHONY: docker

migrate-up:
	@echo "Migrating database up..."
	@CONFIG_PATH=./config/local.json go run ./cmd/migrator -migrations-path ./migrations -direction up

migrate-down:
	@echo "Migrating database down..."
	@CONFIG_PATH=./config/local.json go run ./cmd/migrator -migrations-path ./migrations -direction down

run:
	@echo "Running server..."
	@CONFIG_PATH=./config/local.json go run ./cmd/app

jwt-user:
	@echo "Generating JWT token for user..."
	@CONFIG_PATH=./config/local.json go run ./cmd/jwt-generator -role user

jwt-admin:
	@echo "Generating JWT token for admin..."
	@CONFIG_PATH=./config/local.json go run ./cmd/jwt-generator -role admin

lint:
	@echo "Running linter..."
	@golangci-lint run ./... -c ./config/.golangci.yml

docker:
	@echo "Running docker-compose..."
	@docker-compose -f ./docker/local.docker-compose.yaml up -d

docker-down:
	@echo "Stopping docker-compose..."
	@docker-compose -f ./docker/local.docker-compose.yaml down --remove-orphans

docker-deps:
	@echo "Running dependencies in docker..."
	@docker-compose -f ./docker/deps.docker-compose.yaml up -d
	@CONFIG_PATH=./config/local.docker.deps.json go run ./cmd/migrator -migrations-path ./migrations -direction up

docker-deps-down:
	@echo "Stopping dependencies in docker..."
	@docker-compose -f ./docker/deps.docker-compose.yaml down --remove-orphans

run-docker-deps:
	@echo "Running server (docker-deps)..."
	@CONFIG_PATH=./config/local.docker.deps.json go run ./cmd/app

test: docker-deps
	@echo "Running tests..."
	-@go test ./tests
	@make docker-deps-down
