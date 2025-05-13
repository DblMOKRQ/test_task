.PHONY: vet test run-server run-worker docker-up docker-down

vet:
	go vet ./...

test:
	go test -race ./...

run-server:
	go run ./cmd/server

run-worker:
	go run ./cmd/worker

docker-up:
	docker-compose up --build --detach --wait

docker-down:
	docker-compose down