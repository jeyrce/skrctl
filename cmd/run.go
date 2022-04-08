package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	yes     = "Y"
	no      = "N"
	timeout = 10
	running = "active (running)"
	dead    = "inactive (dead)"
	null    = "-"
	theSame = "↑" // 代表和上一个相同
)

// 一些重要操作需要输入进行确认
func YesOrNo(tip string) {
	var input string
	yellow := color.New(color.FgYellow)
	for {
		_, _ = yellow.Printf("%s, 是否确认(Y/n):", tip)
		_, _ = fmt.Scanln(&input)
		switch strings.TrimSpace(strings.ToUpper(input)) {
		case yes:
			fmt.Println("已确认操作...")
			goto confirmed
		case no:
			fmt.Println("已取消操作...")
			os.Exit(0)
		default:
			continue
		}
	}
confirmed:
	return
}

// 执行操作系统命令
func runCmd(cmd string, args ...string) (output string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	var Cmd = exec.CommandContext(ctx, cmd, args...)
	stdout, _ := Cmd.StdoutPipe()
	defer func() {
		_ = stdout.Close()
	}()
	if err := Cmd.Start(); err != nil {
		return "", err
	}
	outByte, _ := ioutil.ReadAll(stdout)
	return string(outByte), nil
}

// 命令提示信息
func help(op, desc string) {
	// 自动补全空格输出
	blankLen := 32
	switch op {
	case "version":
		fmt.Printf("\t%s%sto %s\n", op, strings.Repeat(" ", blankLen-7), desc)
	case "help":
		fmt.Printf("\t%s [cmd_name]%sto %s the cmd named\n", op, strings.Repeat(" ", blankLen-len(op)-11), desc)
	default:
		fmt.Printf("\t%s%sto %s all of service in conf\n", op, strings.Repeat(" ", blankLen-len(op)), desc)
		fmt.Printf("\t%s [service_name]%sto %s the service named\n", op, strings.Repeat(" ", blankLen-len(op)-15), desc)
	}
}

// service 对象
type Service struct {
	serviceName   string // 待解析的service名称
	statusOutput  string // systemctl status 的输出
	autoStartFile string // 开机自启文件
}

// 判断是否是一个服务
func (s *Service) IsValid() bool {
	features := [5]string{
		"Loaded:",
		"Active:",
		"Docs:",
		"Main PID:",
		"CGroup:",
	}
	for _, feature := range features {
		if !strings.Contains(s.statusOutput, feature) {
			return false
		}
	}
	return true
}

// 服务名称
func (s *Service) Name() string {
	regex := regexp.MustCompile(`loaded \(/usr/lib/systemd/system/(.+\.service);*`)
	subMatch := regex.FindAllStringSubmatch(s.statusOutput, 1)
	if len(subMatch) > 0 {
		return strings.TrimSuffix(subMatch[0][1], ".service")
	}
	return strings.TrimSuffix(s.serviceName, ".service")
}

// 服务状态
func (s *Service) Status() string {
	defaultStatus := "unknown"
	regex := regexp.MustCompile(`Active: (.+) since*`)
	subMatch := regex.FindAllStringSubmatch(s.statusOutput, 1)
	if len(subMatch) > 0 {
		var state = subMatch[0][1]
		if len(state) > 16 {
			return dead
		}
		return state
	}
	return defaultStatus
}

// 服务状态持续时长
func (s *Service) TimeDuration() string {
	if s.Status() != running {
		return null
	}
	regex := regexp.MustCompile(`Active: active \(running\) since .*; (.+ ago)*`)
	subMatch := regex.FindAllStringSubmatch(s.statusOutput, 1)
	if len(subMatch) > 0 {
		return subMatch[0][1]
	}
	return null
}

// 主进程pid
func (s *Service) PID() string {
	if s.Status() != running {
		return null
	}
	regex := regexp.MustCompile(`Main PID: (\d+) \(*`)
	subMatch := regex.FindAllStringSubmatch(s.statusOutput, 1)
	if len(subMatch) > 0 {
		return subMatch[0][1]
	}
	return null
}

// 是否开机自启(采用is-enabled查看的不准确)
func (s *Service) AutoStart() string {
	if _, err := os.Stat(s.autoStartFile); err != nil {
		return no
	}
	return yes
}

// TODO: 服务版本信息
func (s *Service) Version() string {
	return null
}

// 进程所占用的端口: 可能是多个
func (s *Service) Ports() []string {
	validPorts := make([]string, 0)
	empty := []string{null}
	pid := s.PID()
	if pid == null {
		return empty
	}
	output, err := runCmd(
		"sh",
		[]string{
			"-c",
			fmt.Sprintf(`netstat -nltp|grep %s|awk '{print $4}'|awk -F: '{print $NF}'`, pid),
		}...,
	)
	if (err != nil) || (output == "") {
		return empty
	}
	ports := strings.Split(output, "\n")
	for _, port := range ports {
		if port != "" {
			validPorts = append(validPorts, port)
		}
	}
	if len(validPorts) < 1 {
		return empty
	}
	sort.SliceStable(validPorts, func(i, j int) bool {
		port1, err1 := strconv.Atoi(validPorts[i])
		port2, err2 := strconv.Atoi(validPorts[j])
		if err1 != nil || err2 != nil {
			return validPorts[i] < validPorts[j]
		}
		return port1 < port2
	})
	return validPorts
}

type TableService struct {
	Name      string `table:"Name"`
	PID       string `table:"PID"`
	Port      string `table:"Port"`
	Status    string `table:"Status"`
	Time      string `table:"Time"`
	AutoStart string `table:"Auto"`
	Version   string `table:"Version"`
}

// 拷贝文件
func CopyFile(src, dst string) error {
	var buf = make([]byte, 5*2^20)
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("invalid file: %s", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	for {
		Byte, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if Byte == 0 {
			break
		}
		_, err = destination.Write(buf[:Byte])
		if err != nil {
			return err
		}
	}
	return nil
}

// 移动文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}
