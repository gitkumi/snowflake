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
