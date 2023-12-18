.PHONY: build
build:
	go mod tidy
	docker build --tag good-threads-backend .

.PHONY: run
run:
	docker-compose down
	docker-compose up -d
	docker-compose logs -f

.PHONY: reset
reset: build run
