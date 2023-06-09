build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pod *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t pod:latest
