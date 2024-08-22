package helper

import (
	"fmt"
	"strconv"
)

type Color int

// Foreground text colors.
const (
	FgBlack   Color = iota + 30 // [常规] 黑色
	FgRed                       // [常规] 红色
	FgGreen                     // [常规] 绿色
	FgYellow                    // [常规] 黄色
	FgBlue                      // [常规] 深蓝色
	FgMagenta                   // [常规] 紫色
	FgCyan                      // [常规] Tiffany蓝
	FgWhite                     // [常规] 白色
)

// Foreground Hi-Intensity text colors.
const (
	FgHiBlack   Color = iota + 90 // [高亮加粗] 黑色
	FgHiRed                       // [高亮加粗] 红色
	FgHiGreen                     // [高亮加粗] 绿色
	FgHiYellow                    // [高亮加粗] 黄色
	FgHiBlue                      // [高亮加粗] 深蓝色
	FgHiMagenta                   // [高亮加粗] 紫色
	FgHiCyan                      // [高亮加粗] Tiffany蓝
	FgHiWhite                     // [高亮加粗] 白色
)

// Colorize a string based on given color.
func ColorSize(s string, c Color) string {
	return fmt.Sprintf("\033[1;%s;40m %s \033[0m", strconv.Itoa((int(c))), s)
}
