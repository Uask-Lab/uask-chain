.PHONY: build clean reset

build:
	GOBIN=$(PWD)/build/bin go install ./cmd/uask_node

run:
	./build/bin/uask_node

clean:
	@rm -rf build/bin

reset:
	@rm -f ./uask/chain.db ./uask/yu.db