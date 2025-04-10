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
	go build -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$(shell git describe --tags --always)'" -o bin/main ./main.go

.PHONY: build-all
build-all:
	@mkdir -p dist
	@echo "Building for multiple platforms..."
	@echo "Version from git describe: $(shell git describe --tags --always)"
	@echo "VERSION env var (if set): $(VERSION)"
	@VERSION_VALUE=$${VERSION:-$(shell git describe --tags --always)}; echo "Using version: $$VERSION_VALUE"; \
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$$VERSION_VALUE'" -o dist/snowflake_linux_amd64 ./main.go; \
	GOOS=linux GOARCH=arm64 go build -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$$VERSION_VALUE'" -o dist/snowflake_linux_arm64 ./main.go; \
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$$VERSION_VALUE'" -o dist/snowflake_darwin_amd64 ./main.go; \
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$$VERSION_VALUE'" -o dist/snowflake_darwin_arm64 ./main.go; \
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$$VERSION_VALUE'" -o dist/snowflake_windows_amd64.exe ./main.go;
	@cd dist && sha256sum * > SHA256SUMS
	@echo "Done! Binaries are in the dist/ directory"

.PHONY: install
install:
	go install -ldflags="-X 'github.com/gitkumi/snowflake/cmd/cli.version=$(shell git describe --tags --always)'" ./...

.PHONY: release
release:
	@read -p "Enter the new version (e.g., v0.20.0): " version; \
	git tag $$version && \
	git push origin $$version && \
	echo "Pushed tag $$version. GitHub Actions will handle the release process."

.PHONY: list-versions
list-versions:
	@git tag -l --sort=-v:refname | head -n 10
