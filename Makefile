.PHONY: test
test:
	go build -o testhelp/fatal testhelp/fatal.go
	go test -v

.PHONY: get-deps
get-deps:
	go get -t -v ./...
