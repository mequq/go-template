# go lang build image
# FROM golang:1.22.1 AS builder
# WORKDIR /app
# ENV CGO_ENABLED=0 
# ENV GOOS=linux
# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn,direct 
# COPY go.mod .
# COPY go.sum .
# RUN go mod download -x
# RUN go mod verify
# COPY . .
# RUN ls -laR
# RUN go generate ./...
# RUN go build -v  -o ./bin/ ./...
# RUN ls -laR


# 2nd phase
FROM gcr.aban.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY bin/app .


COPY config.example.yaml config.yaml
CMD ["/app/app", "--config", "config.yaml"]




