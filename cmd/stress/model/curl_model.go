// Package model 数据模型
package model

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/BeCrafter/commander/helper"
	"github.com/mattn/go-shellwords"
)

// CURL curl参数解析
type CURL struct {
	Data map[string][]string
}

// getDataValue 获取数据
func (c *CURL) getDataValue(keys []string) []string {
	var (
		value = make([]string, 0)
	)
	for _, key := range keys {
		var (
			ok bool
		)
		value, ok = c.Data[key]
		if ok {
			break
		}
	}
	return value
}

// ParseTheFile 从文件中解析curl
func ParseTheFile(path string) (curls []*CURL, err error) {
	if path == "" {
		err = errors.New("路径不能为空")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())
		return
	}
	defer func() {
		_ = file.Close()
	}()

	itemstr := ""
	itemsList := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		itemstr += " " + line

		if strings.HasPrefix(line, "---") {
			itemsList = append(itemsList, itemstr)
			itemstr = ""
		}
	}

	if len(itemstr) > 0 {
		itemsList = append(itemsList, itemstr)
	}

	for _, item := range itemsList {
		curl := &CURL{
			Data: make(map[string][]string),
		}
		args, err := shellwords.Parse(item)
		if err != nil {
			return nil, errors.New("解析文件失败:" + err.Error())
		}
		args = argsTrim(args)
		var key string
		for _, arg := range args {
			arg = removeSpaces(arg)
			if arg == "" {
				continue
			}
			if isURL(arg) {
				curl.Data[keyCurl] = append(curl.Data[keyCurl], arg)
				key = ""
				continue
			}
			if isKey(arg) {
				key = arg
				continue
			}
			curl.Data[key] = append(curl.Data[key], arg)
		}
		curls = append(curls, curl)
	}

	return curls, nil
}

func argsTrim(args []string) []string {
	result := make([]string, 0)
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "\n" {
			continue
		}
		if strings.Contains(arg, "\n") {
			arg = strings.ReplaceAll(arg, "\n", "")
		}
		if strings.Index(arg, "-X") == 0 {
			result = append(result, arg[0:2])
			result = append(result, arg[2:])
		} else {
			result = append(result, arg)
		}
	}
	return result
}

func removeSpaces(data string) string {
	data = strings.TrimFunc(data, func(r rune) bool {
		if r == ' ' || r == '\\' || r == '\n' {
			return true
		}
		return false
	})
	return data
}

func isKey(data string) bool {
	return strings.HasPrefix(data, "-") || strings.HasPrefix(data, keyCurl)
}

func isURL(data string) bool {
	return strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://")
}

// String string
func (c *CURL) String() (url string) {
	curlByte, _ := json.Marshal(c)
	return string(curlByte)
}

const (
	keyCurl = "curl"
)

// GetURL 获取url
func (c *CURL) GetURL() (url string) {
	keys := []string{keyCurl, "--url", "--location"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return
	}
	url = value[0]
	return
}

// GetMethod 获取 请求方式
func (c *CURL) GetMethod() (method string) {
	keys := []string{"-X", "--request"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return c.defaultMethod()
	}
	method = strings.ToUpper(value[0])
	if helper.InArrayStr(method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
		return method
	}
	return c.defaultMethod()
}

// defaultMethod 获取默认方法
func (c *CURL) defaultMethod() (method string) {
	method = http.MethodGet
	body := c.GetBody()
	if len(body) > 0 {
		return http.MethodPost
	}
	return
}

// GetHeaders 获取请求头
func (c *CURL) GetHeaders() (headers map[string]string) {
	headers = make(map[string]string)
	keys := []string{"-H", "--header"}
	value := c.getDataValue(keys)
	for _, v := range value {
		getHeaderValue(v, headers)
	}
	return
}

// GetHeadersStr 获取请求头string
func (c *CURL) GetHeadersStr() string {
	headers := c.GetHeaders()
	bytes, _ := json.Marshal(&headers)
	return string(bytes)
}

// GetBody 获取body
func (c *CURL) GetBody() (body string) {
	keys := []string{"--data", "-d", "--data-urlencode", "--data-raw", "--data-binary"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		body = c.getPostForm()
		return
	}
	body = value[0]
	return
}

// getPostForm get post form
func (c *CURL) getPostForm() (body string) {
	keys := []string{"--form", "-F", "--form-string"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return
	}
	body = strings.Join(value, "&")
	return
}
