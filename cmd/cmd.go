package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// 执行shell命令，默认不区分stdout和stderr
func run(input string, args ...string) (output string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	data, err := exec.CommandContext(timeout, input, args...).CombinedOutput()
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// 用户确认是否继续执行
func yesOrNo(tip string) {
	var ok string
	fmt.Printf("[警告] %s, 是否继续(Y/n):", tip)
	_, err := fmt.Scanln(&ok)
	if err != nil || strings.TrimSpace(strings.ToUpper(ok)) != "Y" {
		os.Exit(0)
	}
}

// 参数处理，简写&&防止重复
func Args(ss *[]string) []string {
	var (
		m    = make(map[string]struct{})
		apps = make([]string, 0, len(*ss))
	)
	for _, s := range *ss {
		s = strings.ToLower(s)
		name := s
		if strings.HasSuffix(s, ".service") {
			name = s[:len(s)-8]
		}
		m[name] = struct{}{}
	}
	for app := range m {
		apps = append(apps, app)
	}
	return apps
}

func Init(dir string) {

}

func Add(services ...string) {

}

func Remove(force bool, services ...string) {
	if force {
		yesOrNo("操作将会停止服务并从systemd移除")
	}
	fmt.Println(services)
}

func Status(services ...string) {

}

func Start(services ...string) {

}

func Stop(services ...string) {

}

func Enable(services ...string) {

}

func Disable(services ...string) {

}
