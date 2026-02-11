# common-cicd

Reusable CI/CD operations for Go projects using [lesiw.io/ops](https://lesiw.io/ops).

## Usage

Import in your project's `.ops/main.go`:

```go
import gofunctions "github.com/kiloexabyte/common-cicd/go-functions"

type Ops struct {
    gofunctions.Ops
}
```

## Available Commands

- `build` - Build the project
- `test` - Run tests with race detection
- `lint` - Format code and run golangci-lint
- `version` - Print Go version
