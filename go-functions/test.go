package gofunctions

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Test() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PKGNAME": "cmdio",
	})
	defer rnr.Close()

	err := rnr.Run("echo", "hello from", rnr.Env("PKGNAME"))
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("go", "test", "-shuffle=on", "./...")
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("echo", "goodbye from", rnr.Env("PKGNAME"))
	if err != nil {
		log.Fatal(err)
	}
}