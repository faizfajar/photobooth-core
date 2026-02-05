# Variables
APP_NAME=photobooth-core
MAIN_PATH=cmd/api/main.go

.PHONY: swag run tidy build clean

# 1. Generate Swagger documentation
swag:
	@echo "==> [SWAGGER] Regenerating documentation..."
	@swag init -g $(MAIN_PATH) --parseDependency --parseInternal

# 2. Development mode dengan Nodemon
run:
	@echo "==> [DEV] Starting server with Nodemon..."
	@nodemon --exec "swag init -g $(MAIN_PATH) --parseDependency --parseInternal && go run $(MAIN_PATH)" --signal SIGTERM -e go,html,json

# 3. Merapikan dependencies
tidy:
	@echo "==> [MOD] Tidying up Go modules..."
	@go mod tidy

# 4. Build untuk Production
build: swag
	@echo "==> [BUILD] Building binary for production..."
	@go build -o bin/$(APP_NAME) $(MAIN_PATH)

# 5. Cleanup
clean:
	@echo "==> [CLEAN] Removing docs and binary..."
	@rm -rf docs
	@rm -rf bin/