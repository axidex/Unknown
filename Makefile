tidy:
	go fmt ./...
	go mod tidy

build:
	docker build -t unknown .

run:
	docker compose up -d --build

run-logs:
	docker compose up --build

run-logs-no-build:
	docker compose up
stop:
	docker compose down

re:
	docker compose down
	docker compose up -d --build