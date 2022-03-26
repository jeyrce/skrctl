package cmd

import (
	"testing"
)

func TestRunCMD(t *testing.T) {
	t.Log(run("whoami"))
	t.Log(run("ls", "-a", "-l"))
	t.Log(run("docker", "ps"))
	t.Log(run("jeyrce", "abc"))
}
