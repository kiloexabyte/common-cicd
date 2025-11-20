package commands

import (

	"lesiw.io/cmdio/sys"
)

func (Ops) Build() {
	var rnr = sys.Runner().WithEnv(map[string]string{
		"PKGNAME": "cmdio",
	})
	defer rnr.Close()

}