package gofunctions

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Lint() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PKGNAME": "cmdio",
	})
	defer rnr.Close()

	err := rnr.Run("echo", "hello from", rnr.Env("PKGNAME"))
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("golangci-lint", "run")
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("echo", "goodbye from", rnr.Env("PKGNAME"))
	if err != nil {
		log.Fatal(err)
	}
}