BINARY_NAME=obp

build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}-darwin cmd/obp/*.go
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux cmd/obp/*.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows.exe cmd/obp/*.go

clean-all:
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows.exe

test:
	go mod tidy
	go test ./...

test_coverage:
	go mod tidy
	go test ./... -coverprofile=coverage.out
