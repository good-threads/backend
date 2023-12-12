.PHONY: build
build:
	docker build --tag good-threads-backend .

.PHONY: run
run:
	docker run -p 8080:3000 good-threads-backend:latest
