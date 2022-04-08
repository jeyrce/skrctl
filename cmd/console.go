package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

/**
控制台颜色输出
*/

const (
	SucceedPrefix = "[OK]"
	FailedPrefix  = "[ER]"
	WarnPrefix    = "[WN]"
)

// 将字符串居中对齐后返回
func AlignCenter(info string, letter string, maxLength int) string {
	var fill string
	Len := len(info)
	if Len < maxLength {
		fill = strings.Repeat(letter, (maxLength-Len)/2)
	} else {
		fill = ""
	}
	return fmt.Sprintf("%s%s%s", fill, info, fill)
}

// 标题使用字符填充后居中对齐输出
// example: =========xx操作=========
func TitleWithFixLetter(info string, letter string) {
	newString := AlignCenter(info, letter, 60)
	HighLight(newString, color.FgHiBlue)
	fmt.Println()
}

// 分割线
func SplitLine(length int) {
	line := color.New(color.FgHiWhite)
	line.Println(strings.Repeat("-", length))
}

// 标题居中，带分割线
func TitleWithSplitLine(info string) {
	newString := AlignCenter(info, " ", 60)
	HighLight(newString, color.FgHiBlue)
	fmt.Println()
	SplitLine(60)
}

// 一个成功执行的提示
func SucceedLine(info string) {
	fmt.Print("\t")
	fmt.Fprint(color.Output, HighLight(SucceedPrefix, color.FgGreen), " ", info)
	fmt.Println()
}

// 一个错误输出的提示
func FailedLine(info string) {
	fmt.Print("\t")
	fmt.Fprint(color.Output, HighLight(FailedPrefix, color.FgRed), " ", info)
	fmt.Println()
}

// 一个警告提示
func WarnLine(info string) {
	fmt.Print("\t")
	fmt.Fprint(color.Output, HighLight(WarnPrefix, color.FgYellow), " ", info)
	fmt.Println()
}

// 一条普通的结果输出
func TipLine(info string, startWith string) {
	fmt.Println(startWith, info)
}

// 高亮显示一段字符，不换行
func HighLight(info string, clr color.Attribute) string {
	var c *color.Color
	c = color.New(clr)
	return c.Sprintf(info)
}
