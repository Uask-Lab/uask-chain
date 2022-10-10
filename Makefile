.PHONY: build clean reset

build:
	GOBIN=$(PWD)/build/bin go install ./cmd/uask_node

run:
	./build/bin/uask_node -k=./yu.toml

clean:
	@rm -rf build/bin

reset:
	@rm -f chain.db yu.db