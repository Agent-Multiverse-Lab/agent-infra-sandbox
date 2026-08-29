.PHONY: fmt test vet check

fmt:
	gofmt -w $$(find . -name '*.go' -type f)

test:
	go test ./...

vet:
	go vet ./...

check:
	test -z "$$(gofmt -l .)"
	go test ./...
	go vet ./...
