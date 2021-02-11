.PHONY: test
test:
	go test -cover

.PHONY: vet
vet:
	go vet ./...
