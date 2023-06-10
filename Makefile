BINARY_NAME=obp

docker: docker-image docker-push

docker-push:
	docker push code.ndumas.com/ndumas/obsidian-pipeline:latest

docker-image:
	docker build -t code.ndumas.com/ndumas/obsidian-pipeline:latest .

alpine-binary:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/${BINARY_NAME}-alpine cmd/obp/cmd/*.go

build-all: alpine-binary
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin cmd/obp/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux cmd/obp/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows.exe cmd/obp/*.go

clean-all:
	go clean
	rm -v bin/*

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out
