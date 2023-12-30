.PHONY: build
build:
	go mod tidy
	go test ./...
	docker build --tag good-threads-backend .

.PHONY: deploy
deploy:
	docker-compose down
	docker-compose up -d

.PHONY: logs
logs:
	docker-compose logs -f backend

.PHONY: from-scratch
from-scratch: build deploy

.PHONY: cover
cover:
	bash create-missing-tests.sh
	go test ./... -coverprofile=cover.out -covermode=count
	bash remove-lines-containing-token.sh cover.out _mock.go
	go tool cover -html=cover.out -o cover.html
	firefox cover.html

.PHONY: clean-db
clean-db:
	docker-compose down
	sudo rm -fr ./.data

.PHONY: e2e-tests
e2e-tests: clean-db from-scratch
	sleep 3
	curl http://localhost:8080/ -b cookies
	curl http://localhost:8080/session -d '{"username":"tom","password":"pepe123"}' -c cookies
	curl http://localhost:8080/user -d '{"username":"tom","password":"pepe123"}'
	curl http://localhost:8080/session -d '{"username":"tom","password":"pepe123"}' -c cookies
	curl http://localhost:8080/ -b cookies
