.PHONY: build
build:
	docker build --tag good-threads-backend .

.PHONY: test-and-build
test-and-build:
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
from-scratch: test-and-build deploy

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
	sudo rm -fr ./.data/mongo
	mkdir -p ./.data/mongo
	sudo chown 999:999 ./.data/mongo

.PHONY: clean-data
clean-data:
	docker-compose down
	sudo rm -fr ./.data
	mkdir -p ./.data/mongo
	sudo chown 999:999 ./.data/mongo
	mkdir -p ./.data/grafana
	sudo chown 472:472 ./.data/grafana
	mkdir -p ./.data/prometheus
	sudo chown 65534:65534 ./.data/prometheus

.PHONY: e2e-tests
e2e-tests: clean-db from-scratch
	sleep 3
	bash -x e2e-tests.sh
