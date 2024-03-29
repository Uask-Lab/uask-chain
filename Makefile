.PHONY: build clean reset

build:
	GOBIN=$(PWD)/build/bin go install ./cmd/uask_node

pull_images:
	docker pull getmeili/meilisearch:v1.2
	docker pull postgres:12-alpine
	# docker pull ipfs/go-ipfs:latest

docker_build:
	docker build -t uask:0.1 .

up: docker_build
	docker-compose pull pg meili
	docker-compose up

stop:
	docker-compose stop

login_db:
	psql postgres://uask:pwd@localhost:5432/uask

run:
	./build/bin/uask_node

clean:
	@rm -rf build/bin

reset:
	@rm -rf ./uask ./meili
	docker-compose rm
