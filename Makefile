
build:
	GOBIN=$(PWD)/build/bin go install ./cmd/uask_node

run:
	./build/bin/uask_node -k=./yu.toml