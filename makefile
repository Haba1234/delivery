APP_NAME=delivery

.PHONY: build
build: ## Build application
	mkdir -p build
	go build -o build/${APP_NAME} cmd/app/main.go
