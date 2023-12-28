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
	@echo "Running..."
	@echo "${PORT}"
	@go run main.go

test:
	@echo "Testing..."
	@go test -v -cover ./...

tidy:
	@echo "Tidying..."
	@go mod tidy

.PHONY: upd downd run tidy migrateup migratedown
