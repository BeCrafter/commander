// Package model 数据模型
package model

import (
	"fmt"
	"testing"
)

// TestCurl 测试函数
func TestCurl(t *testing.T) {
	// ../curl.txt
	cs, err := ParseTheFile("../../../data/test.curl.txt")
	if err != nil {
		return
	}
	c := cs[0]

	fmt.Printf("curl:%s \n", c.String())
	fmt.Printf("url:%s \n", c.GetURL())
	fmt.Printf("method:%s \n", c.GetMethod())
	fmt.Printf("body:%v \n", c.GetBody())
	fmt.Printf("body string:%v \n", c.GetBody())
	fmt.Printf("headers:%s \n", c.GetHeadersStr())
}
