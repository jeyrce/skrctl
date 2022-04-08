package cmd

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// 编译时注入信息
var (
	BuildVersion string // 版本信息
	BuildTime    string // 编译时间
	BuildName    string // 编译用户
	CommitID     string // 编译时代码id
	Signature    string //= "http://www.woqutech.com/" // 签名
)

// yaml 配置文件
type Config struct {
	Systemctl    string   `yaml:"systemctl"`
	ServiceDir   string   `yaml:"serviceDir"`
	Services     []string `yaml:"services"`
	AutoStartDir string   `yaml:"autoStartDir"`
}

func (c *Config) Load(Filename string) error {
	file, err := ioutil.ReadFile(Filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(file, c); err != nil {
		return err
	}
	return nil
}

// 对加载的配置做一些检查
func (c *Config) Validate() {
	switch {
	case c.Systemctl == "":
		FailedLine("未配置systemctl路径")
		os.Exit(0)
	case c.ServiceDir == "":
		FailedLine("未配置serviceDir路径")
		os.Exit(0)
	case len(c.Services) < 1:
		FailedLine("至少需要一个service")
		os.Exit(0)
	case c.AutoStartDir == "":
		FailedLine("未配置开机自启文件目录")
		os.Exit(0)
	}
}
