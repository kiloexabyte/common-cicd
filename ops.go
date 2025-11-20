package commoncicd

import gofunctions "github.com/kiloexabyte/common-cicd/go-functions"

type Ops struct {
	gofunctions.Ops
}

var GoOps = &Ops{
	Ops: gofunctions.Ops{},
}
