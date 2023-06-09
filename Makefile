BINARY_NAME=obp

build-all:
	go version
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin cmd/obp/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux cmd/obp/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows.exe cmd/obp/*.go

clean-all:
	go version
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows.exe

test:
	go version
	go mod tidy
	go test ./...

test_coverage:
	go version
	go mod tidy
	go test ./... -coverprofile=coverage.out
