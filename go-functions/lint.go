package gofunctions

import (
	"log"
	"context"
	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Lint() {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "go")

    if err := sh.Exec(ctx, "golangci-lint", "run"); err != nil {
        log.Fatal(err)
    }
}