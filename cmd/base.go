package cmd

import (
	"fmt"
	"strings"
)

type Command interface {
	Name() string
	Help()
	Execute(...string)
}

type Application struct {
	BaseDir  string
	Conf     Config
	Commands []Command
}

func NewApplication(config Config, baseDir string) *Application {
	return &Application{
		BaseDir: baseDir,
		Conf:    config,
		Commands: []Command{
			Help{conf: &config, baseDir: baseDir},
			Version{conf: &config, baseDir: baseDir},
			Install{conf: &config, baseDir: baseDir},
			Uninstall{conf: &config, baseDir: baseDir},
			Start{conf: &config, baseDir: baseDir},
			Stop{conf: &config, baseDir: baseDir},
			Restart{conf: &config, baseDir: baseDir},
			Status{conf: &config, baseDir: baseDir},
			Enable{conf: &config, baseDir: baseDir},
			Disable{conf: &config, baseDir: baseDir},
		},
	}
}

func (app *Application) Usage() {
	// 命令工具的用法
	TipLine("Usage: ./cloudctl [options] [service_name]", "")
	SplitLine(80)
	for _, cmd := range app.Commands {
		cmd.Help()
	}
}

func (app *Application) Error(cmd string) {
	// 未匹配任何命令
	FailedLine(fmt.Sprintf("cmd [%s] not supported\n", cmd))
}

// 依据配置文件过滤用户输入中可管理的service
func (app *Application) FilterAvailableService(args ...string) *[]string {
	if len(args[1:]) > 0 {
		// 过滤可管理service
		var services []string
		for _, arg := range args {
			if !strings.HasSuffix(arg, ".service") {
				arg = fmt.Sprintf("%s%s", arg, ".service")
			}
			for _, s := range app.Conf.Services {
				if arg == s {
					services = append(services, arg)
				}
			}
		}
		return &services
	}
	return &app.Conf.Services
}
