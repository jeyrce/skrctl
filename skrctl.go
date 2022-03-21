package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/skrbox/skrctl/cmd"
)

var (
	skrctl = kingpin.CommandLine
	// init
	cmdInit = skrctl.Command("init", "初始化skrctl配置")
	// add
	cmdAdd     = skrctl.Command("add", "将服务进行接管")
	cmdAddList = cmdAdd.Arg("service file", "服务描述文件").Strings()
	// rm
	cmdRemove     = skrctl.Command("remove", "移除已接管服务")
	_             = cmdRemove.Flag("force", "移除接管并卸载").Short('f').Default("0").Bool()
	cmdRm         = skrctl.Command("rm", "alias of remove")
	_             = cmdRm.Flag("force", "移除接管并卸载").Short('f').Default("0").Bool()
	cmdRmList     = cmdRm.Arg("service", "服务名称").Strings()
	cmdRemoveList = cmdRemove.Arg("service", "服务名称").Strings()
	// ps
	cmdStatus     = skrctl.Command("status", "查看服务状态")
	cmdPs         = skrctl.Command("ps", "alias of status")
	cmdPsList     = cmdPs.Arg("service", "服务名").Strings()
	cmdStatusList = cmdStatus.Arg("service", "服务名").Strings()
	// start
	cmdStart     = skrctl.Command("start", "尝试启动服务")
	cmdRun       = skrctl.Command("run", "alias of start")
	cmdUp        = skrctl.Command("up", "alias of start")
	cmdUpList    = cmdUp.Arg("service", "服务名").Strings()
	cmdRunList   = cmdRun.Arg("service", "服务名").Strings()
	cmdStartList = cmdStart.Arg("service", "服务名").Strings()
	// stop
	cmdStop     = skrctl.Command("stop", "尝试停止服务")
	cmdDown     = skrctl.Command("down", "alias of stop")
	cmdDownList = cmdDown.Arg("service", "服务名").Strings()
	cmdStopList = cmdStop.Arg("service", "服务名").Strings()
	// restart
	cmdRestart     = skrctl.Command("restart", "尝试重启服务")
	cmdReload      = skrctl.Command("reload", "alias of restart")
	cmdReloadList  = cmdReload.Arg("service", "服务名").Strings()
	cmdRestartList = cmdRestart.Arg("service", "服务名").Strings()
	// set/unset auto start
	cmdEnable      = skrctl.Command("enable", "设置为开机自启")
	cmdEnableList  = cmdEnable.Arg("service", "服务名").Strings()
	cmdDisable     = skrctl.Command("disable", "关闭开机自启")
	cmdDisableList = cmdDisable.Arg("service", "服务名").Strings()
)

func init() {
	kingpin.Version("v0.1.0-beta").VersionFlag.Short('v')
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
}

func main() {
	switch kingpin.Parse() {
	case cmdInit.FullCommand():
	case cmdAdd.FullCommand():
		fmt.Println(cmdAddList)
	case cmdRemove.FullCommand(), cmdRm.FullCommand():
	case cmdStatus.FullCommand(), cmdPs.FullCommand():
	case cmdStart.FullCommand(), cmdUp.FullCommand(), cmdRun.FullCommand():
		cmd.Start(*cmdStartList...)
	case cmdStop.FullCommand(), cmdDown.FullCommand():
	case cmdRestart.FullCommand(), cmdReload.FullCommand():
	case cmdEnable.FullCommand():
	case cmdDisable.FullCommand():
	default:
		kingpin.Usage()
	}
	os.Exit(0)
}
