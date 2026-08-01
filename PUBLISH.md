## How to publish

Never hand-type a tag. Pick a bump and let the script work out the version:

```sh
make release.patch    # bug fixes        v1.104.0 -> v1.104.1
make release.minor    # features         v1.104.0 -> v1.105.0
make release.major    # breaking changes v1.104.0 -> v2.0.0
```

The script refuses to tag unless all of the following hold, then runs the tests
and pushes. The Release workflow takes over from there.

- you are on `master`, with a clean working tree
- your `HEAD` matches `origin/master` (the commit is already on GitHub)
- the tag does not already exist locally or on the remote
- the version is a **legal successor** to the previous tag
- the `go.mod` module path matches the major version

`SKIP_TESTS=1 make release.minor` skips the local test run when CI has already
gone green on that exact commit.

## Why the guardrails exist

**Lexicographic sorting lies.** `git tag -l` sorts as text, so it lists
`v0.99.0` *after* `v0.103.0`. Use `make versions` (which sorts with `-v:refname`)
whenever you need to see the real latest version.

**A bump resets everything below it.** After `v0.103.0` the only legal next
versions are `v0.103.1`, `v0.104.0`, or `v1.0.0`. Carrying the minor across a
major bump is how `v1.104.0` happened.

**v2+ needs a `/v2` module path.** Go refuses to serve a `v2.0.0` tag on a
module named `github.com/gitkumi/snowflake`. The tag pushes fine and the failure
is quiet: `go install ...@latest` just keeps resolving to the last v1. Before
tagging a major:

1. set `module github.com/gitkumi/snowflake/v2` in `go.mod`
2. rewrite internal imports to `github.com/gitkumi/snowflake/v2/...`
3. commit and push, then `make release.major`

Because of this, `major` is rarely what you want for a CLI — the semver major
buys you nothing that users see, and costs a module path migration. Prefer
`minor`.

## Recovering from a bad tag

A tag is immutable *once the Go proxy has served it* — after that, the version
is burned and you must move forward, never reuse. But a tag the proxy has not
yet fetched can still be deleted:

```sh
curl -s https://proxy.golang.org/github.com/gitkumi/snowflake/@v/vX.Y.Z.info
# 404 means it was never published and is safe to delete:
git push origin :refs/tags/vX.Y.Z
git tag -d vX.Y.Z
```

Check what the proxy actually believes is current — this, not your local tag
list, is what users get:

```sh
curl -s https://proxy.golang.org/github.com/gitkumi/snowflake/@latest
```
