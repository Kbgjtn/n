upd:
	@echo "Initializing database..."
	@docker-compose up -d

downd:
	@echo "Stopping database..."
	@docker-compose down -v

migrateup:
	@echo "Running migration schema"
	@migrate

migratedown:
	@echo "Running migration schema"
	@migrate

run: 
	@echo "Running Server..."
	@go run main.go

test:
	@echo "Testing..."
	@go test -v -cover ./...

tidy:
	@echo "Tidying..."
	@go mod tidy

setup:
	@echo "Setting up server...\n"
	@echo "Copying .env.example to .env.local...\n"
	@cp -pf .env.example .env.local
	@echo "Installing dependencies...\n"
	@go mod download
	@echo "tidying...\n"
	@go mod tidy
	@echo "create database...\n"
	@docker-compose up -d
	@echo " ===== Setting up is done  ===== "
	@echo " ===== Run 'make run' to start server ===== "

.PHONY: upd downd run tidy migrateup migratedown test setup
