.PHONY: dev build clean test lint run

dev:
	@echo "Starting backend..."
	@go run ./cmd/server/ &
	@echo "Starting frontend..."
	@cd web && npm run dev

build:
	@echo "Building backend..."
	@CGO_ENABLED=1 go build -o bin/server ./cmd/server/
	@echo "Building frontend..."
	@cd web && npm run build

clean:
	@rm -rf bin/ web/dist/ data/ uploads/ generated/

test:
	@go test ./... -v
	@cd web && npx vitest run

lint:
	@go vet ./...
	@cd web && npx tsc --noEmit

run:
	@CGO_ENABLED=1 go run ./cmd/server/

docker-build:
	@docker compose build

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down
