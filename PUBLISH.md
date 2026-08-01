## How to publish

Pick a bump. The script works out the version.

```sh
make release.patch    # bug fixes        v1.105.0 -> v1.105.1
make release.minor    # features         v1.105.0 -> v1.106.0
make release.major    # breaking changes v1.105.0 -> v2.0.0
```

It runs the tests, tags, and pushes. The Release workflow publishes from there.

Before tagging it checks that you are on a clean `master`, that `HEAD` is
already on GitHub, that the tag is a legal successor to the last one, and that
the `go.mod` module path matches the major version. Use `SKIP_TESTS=1` to skip
the local test run when CI is already green on that commit.

Never write a tag by hand — `git tag -l` sorts as text and lists `v0.99.0`
after `v0.103.0`. Run `make versions` to see the real latest.

A major bump also needs the module path renamed to `.../v2` and every internal
import rewritten first, so prefer `minor` unless you mean it.
