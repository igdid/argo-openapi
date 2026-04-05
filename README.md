# Argo-OpenAPI

Inspired by [Argo Events](https://argoproj.github.io/events/) and [Cobra CLI](https://github.com/spf13/cobra-cli).

## Installation
1. Install dependencies (Protobuf + GRPC)

```
brew install protobuf
protobuf --version
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/spf13/cobra-cli@latest
```

2. Run it

```
go run main.go generate client -f examples/petstore/store.yaml
```

## Usage
1. **Create a Git repo anywhere.** Argo-OpenAPI uses **origin** from `git remote -v` for getting go_package address. So this tool must be run in a Git repo.
2. 
