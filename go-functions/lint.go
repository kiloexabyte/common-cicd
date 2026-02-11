package gofunctions

import (
	"context"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Lint() error {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "golangci-lint", "go")
	if err := sh.Exec(ctx, "go", "fmt", "./..."); err != nil {
		return err
	}
	return sh.Exec(ctx, "golangci-lint", "run")
}