.PHONY: build
build:
	go mod tidy
	go test ./...
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
	bash create-missing-tests.sh
	go test ./... -coverprofile=cover.out -covermode=count
	bash remove-mock-from-coverage.sh cover.out _mock.go
	go tool cover -html=cover.out -o cover.html
	firefox cover.html