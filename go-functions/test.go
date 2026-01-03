package gofunctions

import (
	"context"
	"log"
	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Test() {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "go")

    if err := sh.Exec(ctx, "go", 
		"test", 
		"-race","-shuffle=on", 
	"./..."); err != nil {
        log.Fatal(err)
    }
}