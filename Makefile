.PHONY: audit
audit: 
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

.PHONY: run
run:
	go run ./main.go

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o bin/main ./main.go
