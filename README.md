# Golang Boilerplate (template)
This simple boilerplate, based on clean architecture principles, serves as an excellent foundation for your REST API microservices.
.
## Setup 
```
cp config.example.yaml config.yaml 
```
### Commands

#### Generate

Generates the necessary files for the project.

```sh
make generate
```

This command performs the following steps:
- Tidies the `go.mod` file.
- Installs the latest version of `wire`.
- Runs `go generate` on the project.
- Tidies the `go.mod` file again.

#### Generate All

Runs the `generate` command.

```sh
make all
```

#### All Tests

Runs all tests with benchmarking, code coverage, and memory allocation statistics.

```sh
make all_tests
```

#### Benchmark Tests

Runs benchmark tests with memory allocation statistics.

```sh
make bench_tests
```

#### Unit Tests

Runs unit tests.

```sh
make unit_tests
```

#### Coverage Tests

Runs tests with code coverage and outputs the coverage profile to `coverage.out`.

```sh
make coverage_tests
```

#### Format Code

Formats the code using `gofumpt` and `gci`.

```sh
make fmt
```

#### Install Development Tools

Installs the necessary development tools.

```sh
make devtools
```

This command installs the following tools:
- `golangci-lint`
- `gofumpt`
- `mockgen`
- `swag`

#### Generate Swagger Documentation (v1)

Generates Swagger documentation for version 1 of the API.

```sh
make swagger-v1
```

#### Lint Check

Runs `golangci-lint` with various enabled linters.
<br>
[Go Cli Dcos](https://golangci-lint.run)

```sh
make check
```

#### Build Docker Image

Builds the Docker image for the project.

```sh
make build
```

This will create a Docker image tagged as `buildf`.

### Notes

- Ensure that you have Docker installed and running for the `build` command.
- The `devtools` command installs tools needed for development and should be run once initially or whenever you need to update these tools.
