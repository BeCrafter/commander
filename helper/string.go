package helper

import "github.com/spf13/cast"

// StrLeftPad
// input     string 原字符串
// length    int    规定补齐后的字符串位数
// padString string 自定义填充字符串
func StrLeftPad(input any, length int, padString string) string {

	output := ""
	inpStr := cast.ToString(input)
	inputLen := len(inpStr)

	if inputLen >= length {
		return inpStr
	}

	padStringLen := len(padString)
	needFillLen := length - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}

	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	return output + inpStr
}

// StrRightPad
// input     string 原字符串
// length    int    规定补齐后的字符串位数
// padString string 自定义填充字符串
func StrRightPad(input any, length int, padString string) string {

	output := ""
	inpStr := cast.ToString(input)
	inputLen := len(inpStr)

	if inputLen >= length {
		return inpStr
	}

	padStringLen := len(padString)
	needFillLen := length - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}

	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}

	return inpStr + output
}
