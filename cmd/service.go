package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	systemctl  = "/usr/bin/systemctl"
	autoDir    = "/etc/systemd/system/multi-user.target.wants/"
	serviceDir = "/usr/lib/systemd/system/"
	running    = "active (running)"
	dead       = "inactive (dead)"
	unknown    = "unknown"
	null       = "-"
	yes        = "Y"
	no         = "N"
)

type service struct {
	Name         string
	File         string
	statusOutput string // 保存以减少执行次数
}

func newService(name, file string) *service {
	svc := service{
		Name: strings.TrimSuffix(name, ".service"),
		File: file,
	}
	return &svc
}

func (s *service) FullName() string {
	if !strings.HasSuffix(s.Name, ".service") {
		return s.Name + ".service"
	}
	return s.Name
}

// 将service文件md5作为身份特征
func (s *service) ID() string {
	open, err := os.Open(s.File)
	if err != nil {
		return ""
	}
	md5v := md5.New()
	_, err = io.Copy(md5v, open)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(md5v.Sum(nil))
}

func (s *service) IsRunning() string {
	if s.Status() == running {
		return yes
	}
	return no
}

func (s *service) Stop() {
	run(systemctl, "stop", s.Name)
}

func (s *service) Start() {
	worker := path.Join(serviceDir, s.FullName())
	_, err := os.Stat(worker)
	if err != nil && os.IsNotExist(err) {
		_ = cp(s.File, worker)
	}
	run(systemctl, "start", s.Name)
}

func (s *service) output() string {
	if s.statusOutput == "" {
		s.statusOutput = run(systemctl, "status", s.Name)
	}
	return s.statusOutput
}

// 判断是否是一个服务
func (s *service) IsValid() bool {
	features := [5]string{
		"Loaded:",
		"Active:",
		"Docs:",
		"Main PID:",
		"CGroup:",
	}
	for _, feature := range features {
		if !strings.Contains(s.output(), feature) {
			return false
		}
	}
	return true
}

func (s *service) Status() string {
	regex := regexp.MustCompile(`Active: (.+) since*`)
	subMatch := regex.FindAllStringSubmatch(s.output(), 1)
	if len(subMatch) > 0 {
		var state = subMatch[0][1]
		if len(state) > 16 {
			return dead
		}
		return state
	}
	return unknown
}

// 服务状态持续时长
func (s *service) TimeDuration() string {
	regex := regexp.MustCompile(`Active: active \(running\) since .*; (.+ ago)*`)
	subMatch := regex.FindAllStringSubmatch(s.output(), 1)
	if len(subMatch) > 0 {
		return subMatch[0][1]
	}
	return null
}

// 主进程pid
func (s *service) PID() string {
	if s.Status() != running {
		return null
	}
	regex := regexp.MustCompile(`Main PID: (\d+) \(*`)
	subMatch := regex.FindAllStringSubmatch(s.output(), 1)
	if len(subMatch) > 0 {
		return subMatch[0][1]
	}
	return null
}

// 是否开机自启(采用is-enabled查看的不准确)
func (s *service) IsAutoStart() string {
	if _, err := os.Stat(path.Join(autoDir, s.FullName())); err != nil {
		return no
	}
	return yes
}

// 进程所占用的端口: 可能是多个
func (s *service) Ports() []string {
	validPorts := make([]string, 0)
	pid := s.PID()
	if pid == null {
		return validPorts
	}
	output := run(
		"sh",
		"-c",
		fmt.Sprintf(`netstat -nltp|grep %s|awk '{print $4}'|awk -F: '{print $NF}'`, pid),
	)
	ports := strings.Split(output, "\n")
	for _, port := range ports {
		if port != "" {
			validPorts = append(validPorts, port)
		}
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

func (s *service) Version() string {
	// todo: 标识服务版本
	return null
}

// 设置开机自启
func (s *service) SetAutoStart() {
	run(systemctl, "enable", s.Name)
}

// 取消开机自启
func (s *service) CloseAutoStart() {
	run(systemctl, "disable", s.Name)
}
