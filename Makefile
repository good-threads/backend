.PHONY: build
build:
	go mod tidy
	docker build --tag good-threads-backend .

.PHONY: reset
reset:
	docker-compose down
	docker-compose up -d
	docker-compose logs -f

.PHONY: from-scratch
from-scratch: build reset

.PHONY: test
test:
	go test ./... -coverprofile=cover.out -covermode=count
	go tool cover -html=cover.out -o cover.html
	firefox cover.html