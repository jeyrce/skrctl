package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/skrbox/skrctl/cmd"
)

var (
	// init
	cmdInit = kingpin.Command("init", "初始化skrctl配置")
	// add
	cmdAdd     = kingpin.Command("add", "将服务进行接管")
	cmdAddList = cmdAdd.Arg("service files", "服务描述文件").ExistingFiles()
	// rm
	cmdRemove     = kingpin.Command("remove", "移除已接管服务")
	removeForce   = cmdRemove.Flag("force", "移除接管并卸载").Short('f').Default("0").Bool()
	cmdRm         = kingpin.Command("rm", "alias of remove")
	rmForce       = cmdRm.Flag("force", "移除接管并卸载").Short('f').Default("0").Bool()
	cmdRmList     = cmdRm.Arg("service", "服务名称").Strings()
	cmdRemoveList = cmdRemove.Arg("service", "服务名称").Strings()
	// ps
	cmdStatus     = kingpin.Command("status", "查看服务状态")
	cmdPs         = kingpin.Command("ps", "alias of status")
	cmdPsList     = cmdPs.Arg("service", "服务名").Strings()
	cmdStatusList = cmdStatus.Arg("service", "服务名").Strings()
	// start
	cmdStart     = kingpin.Command("start", "尝试启动服务")
	cmdRun       = kingpin.Command("run", "alias of start")
	cmdUp        = kingpin.Command("up", "alias of start")
	cmdUpList    = cmdUp.Arg("service", "服务名").Strings()
	cmdRunList   = cmdRun.Arg("service", "服务名").Strings()
	cmdStartList = cmdStart.Arg("service", "服务名").Strings()
	// stop
	cmdStop     = kingpin.Command("stop", "尝试停止服务")
	cmdDown     = kingpin.Command("down", "alias of stop")
	cmdDownList = cmdDown.Arg("service", "服务名").Strings()
	cmdStopList = cmdStop.Arg("service", "服务名").Strings()
	// restart
	cmdRestart     = kingpin.Command("restart", "尝试重启服务")
	cmdReload      = kingpin.Command("reload", "alias of restart")
	cmdReloadList  = cmdReload.Arg("service", "服务名").Strings()
	cmdRestartList = cmdRestart.Arg("service", "服务名").Strings()
	// set/unset auto start
	cmdEnable      = kingpin.Command("enable", "设置为开机自启")
	cmdEnableList  = cmdEnable.Arg("service", "服务名").Strings()
	cmdDisable     = kingpin.Command("disable", "关闭开机自启")
	cmdDisableList = cmdDisable.Arg("service", "服务名").Strings()
)

func init() {
	kingpin.Version("v0.1.0-beta").VersionFlag.Short('v')
	kingpin.HelpFlag.Short('h')
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	switch kingpin.Parse() {
	case cmdInit.FullCommand():
		cmd.Init(wd)
	case cmdAdd.FullCommand():
		cmd.Add(cmd.Args(cmdAddList)...)
	case cmdRemove.FullCommand():
		cmd.Remove(*removeForce, cmd.Args(cmdRemoveList)...)
	case cmdRm.FullCommand():
		cmd.Remove(*rmForce, cmd.Args(cmdRmList)...)
	case cmdStatus.FullCommand():
		cmd.Status(cmd.Args(cmdStatusList)...)
	case cmdPs.FullCommand():
		cmd.Status(cmd.Args(cmdPsList)...)
	case cmdStart.FullCommand():
		cmd.Start(cmd.Args(cmdStartList)...)
	case cmdUp.FullCommand():
		cmd.Start(cmd.Args(cmdUpList)...)
	case cmdRun.FullCommand():
		cmd.Start(cmd.Args(cmdRunList)...)
	case cmdStop.FullCommand():
		cmd.Stop(cmd.Args(cmdStopList)...)
	case cmdDown.FullCommand():
		cmd.Stop(cmd.Args(cmdDownList)...)
	case cmdRestart.FullCommand():
		cmd.Stop(cmd.Args(cmdRestartList)...)
		cmd.Start(cmd.Args(cmdRestartList)...)
	case cmdReload.FullCommand():
		cmd.Stop(cmd.Args(cmdReloadList)...)
		cmd.Start(cmd.Args(cmdReloadList)...)
	case cmdEnable.FullCommand():
		cmd.Enable(cmd.Args(cmdEnableList)...)
	case cmdDisable.FullCommand():
		cmd.Disable(cmd.Args(cmdDisableList)...)
	default:
		kingpin.Usage()
	}
}
