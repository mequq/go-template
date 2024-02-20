


.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make generate;



.DEFAULT_GOAL := generate
