package cmd

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	cmds := [][]string{
		{"whoami"},
		{"ping", "-c", "1", "-w", "1", "192.168.1.99"},
		{"whereis", "systemctl"},
	}
	for _, cmd := range cmds {
		out, err := runCmd(cmd[0], cmd[1:]...)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(cmd, ">>>", out)
	}
}

func TestLoadConfig(t *testing.T) {
	c := new(Config)
	err := c.Load("../conf.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(c)
	c.Validate()
}

func TestCmd(t *testing.T) {
	cmds := []Command{
		new(Version),
	}
	for _, cmd := range cmds {
		cmd.Execute([]string{}...)
	}
}
