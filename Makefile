tidy:
	go fmt ./...
	go mod tidy

build:
	docker build -t unknown .