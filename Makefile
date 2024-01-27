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
	docker-compose up -d mongo
	sleep 3
	docker exec -it backend_mongo_1 mongo mongodb://root:example@localhost:27017 --eval 'rs.initiate()'
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
	bash remove-lines-containing-token.sh cover.out _mock.go # TODO(thomasmarlow): there must a way to go without this line
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

.PHONY: e2e-tests
e2e-tests:
	sleep 3
	bash -x e2e-tests.sh

.PHONY: setup-config-ownership
setup-config-ownership:
	sudo chown 999:999 ./static-mounted/mongo-keyfile
	sudo chmod 400 ./static-mounted/mongo-keyfile

.PHONY: setup-logging
setup:
	docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions || true
	docker plugin ls
	sleep 3
	echo '{"log-driver": "loki", "log-opts": {"loki-url": "http://localhost:9005/loki/api/v1/push", "loki-batch-size": "400", "loki-retries": "2", "loki-max-backoff": "800ms", "loki-timeout": "1s", "keep-file": "true"}}' > /tmp/daemon.json
	sudo mv /tmp/daemon.json /etc/docker/daemon.json
	sudo systemctl restart docker

.PHONY: git-add-all
git-add-all:
	sudo chown $$(whoami) ./static-mounted/mongo-keyfile
	git add .