.PHONY: default check test build image

IMAGE_NAME := traefik/whoamimcp

default: check test build

build:
	CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o whoamimcp

test:
	go test -v -cover ./...

check:
	golangci-lint run

image:
	docker build -t $(IMAGE_NAME) .
