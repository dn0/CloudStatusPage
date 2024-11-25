# Development

## Protocol Buffer

### Install protoc

```sh
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### Run protoc

```sh
make proto
```

## Test & lint

```sh
make lint
make test
```

## Build & run

```sh
make build/dev

bin/...
```

## Templ

```sh
go install github.com/a-h/templ/cmd/templ@latest

make templ
```

## Tailwind

```sh
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 ~/bin/tailwindcss

make tailwind
```

## Air

```sh
go install github.com/air-verse/air@latest

make run/mon-web
```


## Copy of production DB

```sh
# In production:
time kubectl exec -n postgres -c postgres core-2 -- pg_dump -Fc -d app > app.dump
# On localhost:
time kubectl exec -n postgres -c postgres -i core-1 -- pg_restore -d app1 --verbose < app.dump
```
