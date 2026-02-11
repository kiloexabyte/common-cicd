package gofunctions

import (
	"context"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Test() error {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "go")
	return sh.Exec(ctx, "go", "test", "-race", "-shuffle=on", "./...")
}
