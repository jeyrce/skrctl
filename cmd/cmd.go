package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"
)

var C = newConf()

// 执行shell命令，默认不区分stdout和stderr
func run(input string, args ...string) (output string) {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
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

// 拷贝文件
func cp(src, dst string) error {
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
	defer func(source *os.File) {
		_ = source.Close()
	}(source)
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
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
func mv(src, dst string) error {
	return os.Rename(src, dst)
}

func Add(services ...string) {
	for _, s := range services {
		if err := C.Add(s); err != nil {
			fmt.Printf("添加%s失败: %s\n", s, err.Error())
		}
	}
}

func Remove(force bool, services ...string) {
	if force {
		yesOrNo("操作将会停止服务并从systemd移除")
	}
	for _, s := range services {
		C.Remove(s)
		if force {
			Stop(s)
			err := os.Remove(path.Join(serviceDir, newService(s, "").FullName()))
			if err != nil {
				fmt.Printf("移除%s失败: %s\n", s, err.Error())
			}
		}
	}
}

// status 输出表格
type view struct {
	Name      string `table:"Name"`
	PID       string `table:"PID"`
	Port      string `table:"Port"`
	Status    string `table:"Status"`
	Time      string `table:"Time"`
	AutoStart string `table:"Auto"`
	Version   string `table:"Version"`
}

func Status(services ...string) {
	// 如果传入空，则展示所有已托管服务状态
	if len(services) < 1 {
		list := C.List()
		for _, item := range list {
			services = append(services, item.Name)
		}
	}
	var views = make([]view, 0, len(services))
	for _, s := range services {
		if svc := C.Has(s); svc != nil {
			views = append(views, view{
				Name:      svc.Name,
				PID:       svc.PID(),
				Port:      strings.Join(svc.Ports(), ","),
				Status:    svc.Status(),
				Time:      svc.TimeDuration(),
				AutoStart: svc.IsAutoStart(),
				Version:   svc.Version(),
			})
		}
	}
	sort.SliceStable(views, func(i, j int) bool {
		return views[i].Name < views[j].Name
	})
	Output(views)
}

// 尝试部署并启动服务
func Start(services ...string) {
	for _, s := range services {
		svc := C.Has(s)
		if svc == nil {
			fmt.Printf("操作前需要先加入管理(%s)\n", s)
			continue
		}
		svc.Start()
	}
}

func Stop(services ...string) {
	for _, s := range services {
		svc := C.Has(s)
		if svc == nil {
			fmt.Printf("操作前需要先加入管理(%s)\n", s)
			continue
		}
		svc.Stop()
	}
}

func Enable(services ...string) {
	for _, s := range services {
		svc := C.Has(s)
		if svc == nil {
			fmt.Printf("操作前需要先加入管理(%s)\n", s)
			continue
		}
		svc.SetAutoStart()
	}
}

func Disable(services ...string) {
	for _, s := range services {
		svc := C.Has(s)
		if svc == nil {
			fmt.Printf("操作前需要先加入管理(%s)\n", s)
			continue
		}
		svc.CloseAutoStart()
	}
}
