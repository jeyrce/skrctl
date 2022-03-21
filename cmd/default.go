package cmd

import (
	"fmt"
	"path"
)

var (
	ConfDir = ".skrctl/"
	Conf    = path.Join(ConfDir, "skrctl.yaml")
)

func Start(args ...string) {
	// todo: ...
	fmt.Println(args)
}
