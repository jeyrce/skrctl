package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"woqutech.com/cloudctl/table"
)

// 工具的帮助提示信息
type Help struct {
	conf    *Config
	baseDir string
}

func (h Help) Name() string {
	return "help"
}

func (h Help) Help() {
	help(h.Name(), "show usage")
}

func (h Help) Execute(args ...string) {
	if len(args) < 1 || args[0] == h.Name() {
		h.Help()
		return
	}
	app := NewApplication(*h.conf, h.baseDir)
	for _, c := range app.Commands {
		if c.Name() == args[0] {
			c.Help()
			return
		}
	}
	app.Error(args[0])
}

// 启动服务
type Start struct {
	conf    *Config
	baseDir string
}

func (s Start) Name() string {
	return "start"
}

func (s Start) Help() {
	help(s.Name(), s.Name())
}

func (s Start) Execute(args ...string) {
	for idx, arg := range args {
		TipLine(fmt.Sprintf("[%d]Try to %s the service: %s", idx, s.Name(), arg), "-> ")
		cmdArgs := append([]string{s.Name()}, arg)
		_, _ = runCmd(s.conf.Systemctl, cmdArgs...)
	}
}

// 停止服务
type Stop struct {
	baseDir string
	conf    *Config
}

func (s Stop) Name() string {
	return "stop"
}

func (s Stop) Help() {
	help(s.Name(), s.Name())
}

func (s Stop) Execute(args ...string) {
	for idx, arg := range args {
		TipLine(fmt.Sprintf("[%d]Try to %s the service: %s", idx, s.Name(), arg), "-> ")
		cmdArgs := append([]string{s.Name()}, arg)
		_, _ = runCmd(s.conf.Systemctl, cmdArgs...)
	}
}

// 重启服务
type Restart struct {
	baseDir string
	conf    *Config
}

func (r Restart) Name() string {
	return "restart"
}

func (r Restart) Help() {
	help(r.Name(), r.Name())
}

func (r Restart) Execute(args ...string) {
	for idx, arg := range args {
		TipLine(fmt.Sprintf("[%d]Try to %s the service: %s", idx, r.Name(), arg), "-> ")
		cmdArgs := append([]string{r.Name()}, arg)
		_, _ = runCmd(r.conf.Systemctl, cmdArgs...)
	}
}

// 查询运行状态、时长、开机自启设置
type Status struct {
	baseDir string
	conf    *Config
}

func (s Status) Name() string {
	return "status"
}

func (s Status) Help() {
	help(s.Name(), fmt.Sprintf("show %s", s.Name()))
}

func (s Status) Execute(args ...string) {
	// 展示service、pid、状态、开机自启状态等
	var StatusData []*TableService
	for _, arg := range args {
		statusOut, _ := runCmd(s.conf.Systemctl, []string{s.Name(), arg}...)
		service := Service{
			serviceName:   arg,
			statusOutput:  statusOut,
			autoStartFile: path.Join(s.conf.AutoStartDir, arg),
		}
		for _, port := range service.Ports() {
			StatusData = append(StatusData, &TableService{
				Name:      service.Name(),
				PID:       service.PID(),
				Port:      port,
				Status:    service.Status(),
				Time:      service.TimeDuration(),
				AutoStart: service.AutoStart(),
				Version:   service.Version(),
			})
		}
	}
	// 按照名称排序(多端口占用的情况会放一起)
	sort.SliceStable(StatusData, func(i, j int) bool { return StatusData[i].Name <= StatusData[j].Name })
	sort.SliceStable(StatusData, func(i, j int) bool {
		// 如果两个服务名字一样，则端口从小到大排序，否则顺序不变
		if StatusData[i].Name == StatusData[j].Name {
			iPort, errI := strconv.Atoi(StatusData[i].Port)
			jPort, errJ := strconv.Atoi(StatusData[j].Port)
			if errI != nil || errJ != nil {
				return StatusData[i].Port < StatusData[j].Port
			}
			return iPort < jPort
		}
		return false
	})
	// 将同一服务的不同端口简单标记
	var previousName string
	for _, s := range StatusData {
		if previousName == s.Name {
			s.Name = theSame
			s.PID = theSame
			s.Time = theSame
			s.Version = theSame
			s.Status = theSame
			s.AutoStart = theSame
			continue
		}
		previousName = s.Name
	}
	table.Output(StatusData)
}

// 设置开机自启
type Enable struct {
	baseDir string
	conf    *Config
}

func (e Enable) Name() string {
	return "enable"
}

func (e Enable) Help() {
	help(e.Name(), e.Name())
}

func (e Enable) Execute(args ...string) {
	for idx, arg := range args {
		TipLine(fmt.Sprintf("[%d]Try to %s the service: %s", idx, e.Name(), arg), "-> ")
		cmdArgs := append([]string{e.Name()}, arg)
		_, _ = runCmd(e.conf.Systemctl, cmdArgs...)
	}
}

// 禁用开机自启
type Disable struct {
	baseDir string
	conf    *Config
}

func (d Disable) Name() string {
	return "disable"
}

func (d Disable) Help() {
	help(d.Name(), d.Name())
}

func (d Disable) Execute(args ...string) {
	for idx, arg := range args {
		TipLine(fmt.Sprintf("[%d]Try to %s the service: %s", idx, d.Name(), arg), "-> ")
		cmdArgs := append([]string{d.Name()}, arg)
		_, _ = runCmd(d.conf.Systemctl, cmdArgs...)
	}
}

// 显示本工具版本信息
type Version struct {
	baseDir string
	conf    *Config
}

func (v Version) Name() string {
	return "version"
}

func (v Version) Help() {
	help(v.Name(), "show version of the tool")
}

func (v Version) Execute(_ ...string) {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	goVersion := runtime.Version()
	var versionInfoTmpl = `
	cloudctl, version {{.buildVersion}} (revision: {{.commitID}})
      build user:       {{.buildName}}
      build date:       {{.buildTime}}
      go version:       {{.goVersion}}
      platform  :       {{.platform}}
	{{.signature}}
	`
	var m = map[string]string{
		"buildVersion": BuildVersion,
		"commitID":     CommitID,
		"buildName":    BuildName,
		"buildTime":    BuildTime,
		"platform":     platform,
		"goVersion":    goVersion,
		"signature":    Signature,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	fmt.Println(strings.TrimSpace(buf.String()))
}

// 安装conf中的service
type Install struct {
	baseDir string
	conf    *Config
}

func (i Install) Name() string {
	return "install"
}

func (i Install) Help() {
	help(i.Name(), i.Name())
}

func (i Install) Execute(args ...string) {
	YesOrNo(fmt.Sprintf("注意！此操作命令将部署:\n  %s\n共%d个服务", strings.Join(args, "\n  "), len(args)))
	// 将services目录中的文件放到/usr/lib/systemd/system/下
	for _, arg := range args {
		src := filepath.Join(i.baseDir, "services", arg)
		dst := filepath.Join(i.conf.ServiceDir, arg)
		if err := CopyFile(src, dst); err != nil {
			continue
		}
	}
	// 重载配置
	TipLine("[0]Try to reload systemd", "-> ")
	_, _ = runCmd(i.conf.Systemctl, "daemon-reload")
	// 启动服务
	Start{conf: i.conf, baseDir: i.baseDir}.Execute(args...)
	// 设置开机自启
	Enable{baseDir: i.baseDir, conf: i.conf}.Execute(args...)
}

// 服务卸载
type Uninstall struct {
	baseDir string
	conf    *Config
}

func (u Uninstall) Name() string {
	return "uninstall"
}

func (u Uninstall) Help() {
	help(u.Name(), u.Name())
}

func (u Uninstall) Execute(args ...string) {
	YesOrNo(fmt.Sprintf("注意！此操作命令将卸载:\n  %s\n共%d个服务", strings.Join(args, "\n  "), len(args)))
	// 关掉服务
	Stop{baseDir: u.baseDir, conf: u.conf}.Execute(args...)
	// 取消开机自启
	Disable{baseDir: u.baseDir, conf: u.conf}.Execute(args...)
	// 移除service文件(并非真正移除而是放到当前services目录)
	for _, arg := range args {
		dst := filepath.Join(u.baseDir, "services", arg)
		src := filepath.Join(u.conf.ServiceDir, arg)
		if err := MoveFile(src, dst); err != nil {
			continue
		}
	}
	TipLine("[0]Try to reload systemd", "-> ")
	// 重载配置
	_, _ = runCmd(u.conf.Systemctl, "daemon-reload")
}
