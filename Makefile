.PHONY: build clean reset

build:
	GOBIN=$(PWD)/build/bin go install ./cmd/uask_node

pull_images:
	docker pull getmeili/meilisearch:v1.2
	docker pull ipfs/go-ipfs:latest

docker_build:
	docker build -t uask:0.1 .

run:
	./build/bin/uask_node

clean:
	@rm -rf build/bin

reset:
	@rm -f ./uask/chain.db ./uask/yu.db