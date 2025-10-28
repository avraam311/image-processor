.PHONY: lint, up, down, build

lint:
	go vet ./...
	golangci-lint run ./...

up:
	docker compose up

buildup:
	docker compose up --build

down:
	docker compose down -v