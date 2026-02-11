package main

import (
	commands "ops/commands"

	"lesiw.io/ops"
)

// run touch ./.ops/mod.go && op -l to recompile
func main() {
	ops.Handle(commands.Ops{})
}
