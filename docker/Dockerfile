# Initial stage: download modules
FROM golang:1.24-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download
RUN go build -o main cmd/api/main.go

# Intermediate stage: Build the binary
FROM golang:1.24-alpine as runner

COPY --from=builder ./app/main ./app/main

WORKDIR ./app

ENV config=docker

EXPOSE 8000

ENTRYPOINT ["./main"]
