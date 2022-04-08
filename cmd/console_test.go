package cmd

import (
	"github.com/fatih/color"
	"testing"
)

func TestColor(t *testing.T) {
	red := color.New(color.FgRed, color.BgBlue)
	red.Println("红色")
}

func TestTitle(t *testing.T) {
	title := "你好，李焕英"
	TitleWithFixLetter(title, "=")
	TitleWithSplitLine(title)
	SucceedLine("完成")
	FailedLine("失败")
	WarnLine("警告")
}

func TestYesOrNo(t *testing.T) {
	YesOrNo("新消息")
}
