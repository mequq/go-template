
ARG DOCKER_IMAGE=golang:1.23.2
FROM ${DOCKER_IMAGE} AS builder
RUN mkdir -p /src

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .

RUN go mod download -x


COPY . .
RUN go build -v  -o ./bin/ ./...
RUN ls -laR 


# 2nd phase
FROM gcr.aban.io/distroless/static-debian12:nonroot

COPY --from=builder /src/bin/ /app/bin/
COPY ./config.example.yaml /app/config.yaml


COPY config.example.yaml config.yaml
CMD ["/app/bin/app", "--config", "config.yaml"]




