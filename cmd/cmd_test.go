package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestRunCMD(t *testing.T) {
	t.Log(run("whoami"))
	t.Log(run("ls", "-a", "-l"))
	t.Log(run("docker", "ps"))
	t.Log(run("jeyrce", "abc"))
}

func TestFileWalk(t *testing.T) {
	getwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	files, err := ioutil.ReadDir(getwd)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func TestStat(t *testing.T) {
	file := "conf.go"
	stat, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stat.Name())
}
