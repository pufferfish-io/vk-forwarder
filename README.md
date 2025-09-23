# vk-forwarder


```
export $(cat .env | xargs) && go run ./cmd/vkforwarder
```

```
go mod tidy
```

```
go build -v -x ./cmd/vkforwarder && rm -f vkforwarder
```

```
docker buildx build --no-cache --progress=plain .
```

```
set -a && source .env && set +a && go run ./cmd/vkforwarder
```

```
git tag v0.1.1
git push origin v0.1.1
```

```
git tag -l
git tag -d vX.Y.Z
git push --delete origin vX.Y.Z
git ls-remote --tags origin | grep 'refs/tags/vX.Y.Z$'
```
