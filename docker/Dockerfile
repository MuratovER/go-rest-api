# Initial stage: download modules
FROM golang:1.24-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download


# Intermediate stage: Build the binary
FROM golang:1.24-alpine as runner

COPY --from=builder ./app ./app

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR ./app

ENV config=docker

EXPOSE 8000

ENTRYPOINT CompileDaemon --build='go build cmd/api/main.go' --command=./main




