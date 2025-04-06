## How to publish

1. Run test

```sh
make test
```

2. Add git tag

```sh 
git tag v0.X.0
```

3. Push the tag. GitHub action should trigger and publish the package.

```sh
git push --tags
```
