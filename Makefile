tidy:
	go fmt ./...
	go mod tidy

build:
	docker build -t unknown .

run:
	docker compose up -d

stop:
	docker compose down