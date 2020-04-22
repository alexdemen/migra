
.PHONY: build
build:
	go build -o ./bin/migra ./cmd/migra.go

.PHONY: postgres-up
postgres-up:
	sudo docker-compose up postgres

.PHONY: docker-down
docker-down:
	sudo docker-compose down