# Publishing 

https://go.dev/doc/modules/publishing

1. Test 

```bash
go test ./...
```

1. Tag

```bash
git tag v0.18.0
```

2. Publish

```bash
GOPROXY=proxy.golang.org go list -m github.com/gitkumi/snowflake@v0.18.0
```
