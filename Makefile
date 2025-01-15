build:
	go build -o garantex-monitor ./cmd/main.go
lint:
	golangci-lint run
run:
	docker compose up
docker-build:
	docker build -t garantex-monitor .