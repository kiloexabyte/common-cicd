package gofunctions

import (
	"context"
	"log"
	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Version() {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "go")

    if err := sh.Exec(ctx, "go", "version"); err != nil {
        log.Fatal(err)
    }
}