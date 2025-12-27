
BINARY := cmd/main
PACKAGE := cmd/main.go

.PHONY: lint test build run_cluster stop_cluster run_client run_seed run_node

lint:
	golangci-lint run

test:
	go test -v internal/db/storage_test.go
	go test -v internal/broadcaster/spreader_test.go
	go test -v internal/cache/bucket_center_test.go
	go test -v internal/cluster/cluster_join_test.go
	go test -v internal/cluster/members_test.go

run_cluster:
	docker-compose up --scale node=5

stop_cluster:
	docker-compose stop

run_client:
	go run ./client/client.go

build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) $(PACKAGE)

run_seed: build
	@echo "Running seed node..."
	./$(BINARY) -host=127.0.0.1 -port=6767 -seed=true

# Usage: make run_node <host> <port> <seed-flag>
# example: make run_node 127.0.0.1 8081 false
run_node: build
	@echo "Running node..."
	./$(BINARY) -host=$(word 2,$(MAKECMDGOALS)) -port=$(word 3,$(MAKECMDGOALS)) -seed=$(word 4,$(MAKECMDGOALS))
