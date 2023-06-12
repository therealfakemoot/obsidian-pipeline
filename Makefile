BINARY_NAME=obp

$DOCKER_CMD="docker --config ~/.docker/"

.PHONY: docker
docker: docker-image docker-push

.PHONY: docker-push
docker-push:
	$DOCKER_CMD push code.ndumas.com/ndumas/obsidian-pipeline:latest

.PHONY: docker-image
docker-image:
	$DOCKER_CMD build -t code.ndumas.com/ndumas/obsidian-pipeline:latest .

.PHONY: alpine-binary
alpine-binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/${BINARY_NAME}-alpine cmd/obp/cmd/*.go

.PHONY: build-all
build-all: alpine-binary
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin cmd/obp/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux cmd/obp/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows.exe cmd/obp/*.go

.PHONY: clean
clean:
	go clean
	rm -v bin/*

.PHONY: test
test:
	go test ./...

.PHONY: test_coverage
test_coverage:
	go test ./... -coverprofile=coverage.out
