
ARG DOCKER_IMAGE=golang:1.24.3-bullseye
FROM ${DOCKER_IMAGE} AS builder
RUN mkdir -p /src

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .

RUN go mod download -x


COPY . .
RUN ls -laR /src
RUN mkdir -p /tmp/.bin
RUN go build -v  -o /tmp/.bin/ ./...
RUN ls -laR /tmp/.bin


# 2nd phase
FROM gcr.aban.io/distroless/static-debian12:nonroot

COPY --from=builder /tmp/.bin/ /app/bin/
COPY ./config.example.yaml /app/config.yaml


COPY config.example.yaml /app/config.yaml
CMD ["/app/bin/app", "--config", "/app/config.yaml"]




