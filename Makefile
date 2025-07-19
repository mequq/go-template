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

all_tests:
	go test -v ./internal/biz/... -v  ./internal/handler/...  -bench=. -cover  -coverprofile=coverage.out -benchmem -cpu=1 -timeout=500ms -json


bench_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -bench=. -benchmem -cpu=1,2,3,4 -timeout=500ms

unit_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/...

coverage_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -cover  -coverprofile=coverage.out



devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/daixiang0/gci@v0.11.2



swagger:
	swag i  -d internal/service/handler/ -g wire.go  -pd



check:
	golangci-lint run  \
		--build-tags "${BUILD_TAG}" \
		--timeout=20m0s \
		--enable=gofmt \
		--enable=unconvert \
		--enable=unparam \
		--enable=asciicheck \
		--enable=misspell \
		--enable=revive \
		--enable=decorder \
		--enable=reassign \
		--enable=usestdlibvars \
		--enable=nilerr \
		--enable=gosec \
		--enable=exportloopref \
		--enable=whitespace \
		--enable=gocyclo \
		--enable=nestif \
		--enable=gochecknoinits \
		--enable=gocognit \
		--enable=funlen \
		--enable=forbidigo \
		--enable=godox \
		--enable=gocritic \
		--enable=gci \
		--enable=lll \
		--config=issues.exclude.yaml




build:
	docker build -t buildf .

