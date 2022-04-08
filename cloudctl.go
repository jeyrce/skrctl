package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"woqutech.com/cloudctl/cmd"
)

var baseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

// 一些前置操作
func init() {
	// 校验cloudctl.yml和service目录
	config := path.Join(baseDir, "cloudctl.yml")
	_, err := os.Stat(config)
	if err != nil {
		cmd.FailedLine("配置文件不存在")
		os.Exit(0)
	}
	var c = new(cmd.Config)
	err = c.Load(config)
	if err != nil {
		cmd.FailedLine("加载配置失败")
		os.Exit(0)
	}
	c.Validate()
	// 加载命令行参数
	flag.Parse()
}

func main() {
	args := flag.Args()
	c := new(cmd.Config)
	conf := path.Join(baseDir, "cloudctl.yml")
	_ = c.Load(conf)
	app := cmd.NewApplication(*c, baseDir)
	if len(args) < 1 {
		app.Usage()
		return
	}

	for _, c := range app.Commands {
		if c.Name() == args[0] {
			if args[0] == "help" {
				c.Execute(args[1:]...)
				return
			}
			services := app.FilterAvailableService(args...)
			c.Execute(*services...)
			return
		}
	}

	app.Error(args[0])
	return
}
