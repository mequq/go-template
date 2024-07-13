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
	go test -v ./internal/http/handler/... ./internal/biz/... -bench=. -cover  -coverprofile=coverage.out -benchmem -cpu=1,2,3,4 -timeout=500ms


bench_tests:
	go test -v ./internal/http/handler/... ./internal/biz/... -bench=. -benchmem -cpu=1,2,3,4 -timeout=500ms

unit_tests:
	go test -v ./internal/http/handler/... ./internal/biz/...

coverage_tests:
	go test -v ./internal/http/handler/... ./internal/biz/... -cover  -coverprofile=coverage.out

fmt:
	gofumpt -l -w .

devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install go.uber.org/mock/mockgen@
	go install github.com/swaggo/swag/cmd/swag@latest



swagger:
	swag init --parseDependency -g ./cmd/main.go -o ./docs

check:
	golangci-lint run \
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
		--enable=goimports \
		--enable=gocyclo \
		--enable=nestif \
		--enable=gochecknoinits \
		--enable=gocognit \
		--enable=funlen \
		--enable=forbidigo \
		--enable=godox \
		--enable=gocritic \
		--enable=gci \
		--enable=lll



build:
	docker build -t buildf .

