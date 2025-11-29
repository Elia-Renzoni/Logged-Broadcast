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

build:
	go build cmd/main.go
