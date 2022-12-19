# go lang build image
FROM golang:1.19 AS builder
WORKDIR /app
ENV CGO_ENABLED=0 
ENV GOOS=linux
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct 
COPY go.mod .
COPY go.sum .
RUN go mod download -x 
COPY . .
RUN go build -a -installsuffix cgo -o ./bin ./...


# 2nd phase
FROM scratch
COPY --from=builder /app/bin/cmd .
COPY config.yaml .
CMD ["./cmd"]




