.PHONY: audit
audit: 
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet $(go list ./...)
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 $(go list ./...)
	go run golang.org/x/vuln/cmd/govulncheck@latest $(go list ./...)

.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt $(go list ./... | grep -v '^./grammar')

.PHONY: run
run:
	go run ./main.go

.PHONY: test
test:
	gotestsum -f testname

.PHONY: build
build:
	go build -o bin/main ./main.go

.PHONY: build-all
build-all:
	@mkdir -p dist
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o dist/snowflake_linux_amd64 ./main.go; \
	GOOS=linux GOARCH=arm64 go build -o dist/snowflake_linux_arm64 ./main.go; \
	GOOS=darwin GOARCH=amd64 go build -o dist/snowflake_darwin_amd64 ./main.go; \
	GOOS=darwin GOARCH=arm64 go build -o dist/snowflake_darwin_arm64 ./main.go; \
	GOOS=windows GOARCH=amd64 go build -o dist/snowflake_windows_amd64.exe ./main.go;
	@cd dist && sha256sum * > SHA256SUMS
	@echo "Done! Binaries are in the dist/ directory"

.PHONY: versions
versions:
	@git tag -l --sort=-v:refname | head -n 10
