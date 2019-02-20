.PHONY: test
test:
	go test -v ./...

.PHONY: start
start:
	go run cmd/gun/main.go

.PHONY: release
release:
	@rm -rf dist
	goreleaser
