// Package model 请求数据模型package model
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
)

type Params struct {
	URL    string   `json:"url"    validate:"required"`  // 压测的url 目前支持，http/https ws/wss
	Header []string `json:"header" validate:"omitempty"` // 自定义头信息传递给服务器
	Body   string   `json:"body"   validate:"omitempty"` // HTTP POST方式传送数据
	Verify string   `json:"verify" validate:"omitempty"` // verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
	Code   int      `json:"code"   validate:"omitempty"` // 成功状态码
}

// ParseParams 解析参数
func ParseParams(str string, params *Params) error {
	if err := json.NewDecoder(strings.NewReader(str)).Decode(params); err != nil {
		return err
	}
	return validator.New().Struct(params)
}

// 返回 code 码
const (
	// HTTPOk 请求成功
	HTTPOk = 200
	// RequestErr 请求错误
	RequestErr = 509
	// ParseError 解析错误
	ParseError = 510 // 解析错误
)

// 支持协议
const (
	// FormTypeHTTP http 协议
	FormTypeHTTP = "http"
	// FormTypeWebSocket webSocket 协议
	FormTypeWebSocket = "webSocket"
	// FormTypeGRPC grpc 协议
	FormTypeGRPC   = "grpc"
	FormTypeRadius = "radius"
)

// 校验函数
var (
	// verifyMapHTTP http 校验函数
	verifyMapHTTP = make(map[string]VerifyHTTP)
	// verifyMapHTTPMutex http 并发锁
	verifyMapHTTPMutex sync.RWMutex
	// verifyMapWebSocket webSocket 校验函数
	verifyMapWebSocket = make(map[string]VerifyWebSocket)
	// verifyMapWebSocketMutex webSocket 并发锁
	verifyMapWebSocketMutex sync.RWMutex
)

// RegisterVerifyHTTP 注册 http 校验函数
func RegisterVerifyHTTP(verify string, verifyFunc VerifyHTTP) {
	verifyMapHTTPMutex.Lock()
	defer verifyMapHTTPMutex.Unlock()
	key := fmt.Sprintf("%s.%s", FormTypeHTTP, verify)
	verifyMapHTTP[key] = verifyFunc
}

// RegisterVerifyWebSocket 注册 webSocket 校验函数
func RegisterVerifyWebSocket(verify string, verifyFunc VerifyWebSocket) {
	verifyMapWebSocketMutex.Lock()
	defer verifyMapWebSocketMutex.Unlock()
	key := fmt.Sprintf("%s.%s", FormTypeWebSocket, verify)
	verifyMapWebSocket[key] = verifyFunc
}

// Verify 验证器
type Verify interface {
	GetCode() int    // 有一个方法，返回code为200为成功
	GetResult() bool // 返回是否成功
}

// VerifyHTTP http 验证
type VerifyHTTP func(request *Request, response *http.Response, body []byte) (code int, isSucceed bool)

// VerifyWebSocket webSocket 验证
type VerifyWebSocket func(request *Request, seq string, msg []byte) (code int, isSucceed bool)

// Request 请求数据
type Request struct {
	URL       string            // URL
	Form      string            // http/webSocket/tcp
	Method    string            // 方法 GET/POST/PUT
	Headers   map[string]string // Headers
	Body      string            // body
	Verify    string            // 验证的方法
	Timeout   time.Duration     // 请求超时时间
	Debug     bool              // 是否开启Debug模式
	MaxCon    int               // 每个连接的请求数
	HTTP2     bool              // 是否使用http2.0
	Keepalive bool              // 是否开启长连接
	Code      int               // 验证的状态码
	Redirect  bool              // 是否重定向
}

// GetBody 获取请求数据
func (r *Request) GetBody() (body io.Reader) {
	return strings.NewReader(r.Body)
}

// CopyHeaders copy Headers
func (r *Request) CopyHeaders() map[string]string {
	var result = make(map[string]string, len(r.Headers))
	for k, v := range r.Headers {
		result[k] = v
	}
	return result
}

// getVerifyKey 获取校验 key
func (r *Request) getVerifyKey() (key string) {
	return fmt.Sprintf("%s.%s", r.Form, r.Verify)
}

// GetVerifyHTTP 获取数据校验方法
func (r *Request) GetVerifyHTTP() VerifyHTTP {
	verify, ok := verifyMapHTTP[r.getVerifyKey()]
	if !ok {
		panic("GetVerifyHTTP 验证方法不存在:" + r.Verify)
	}
	return verify
}

// GetVerifyWebSocket 获取数据校验方法
func (r *Request) GetVerifyWebSocket() VerifyWebSocket {
	verify, ok := verifyMapWebSocket[r.getVerifyKey()]
	if !ok {
		panic("GetVerifyWebSocket 验证方法不存在:" + r.Verify)
	}
	return verify
}

// NewRequest 生成请求结构体
// url 压测的url
// verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
// timeout 请求超时时间
// debug 是否开启debug
// path curl文件路径 http接口压测，自定义参数设置
func NewRequest(params Params, timeout time.Duration, maxCon int,
	http2, keepalive, redirect, debug bool) (request *Request, err error) {
	var (
		method  = "GET"
		headers = make(map[string]string)
		body    string
	)

	if params.Body != "" {
		method = "POST"
		body = params.Body
	}
	for _, v := range params.Header {
		getHeaderValue(v, headers)
	}
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	}

	var form string
	form, params.URL = getForm(params.URL)
	if form == "" {
		err = fmt.Errorf("url:%s 不合法,必须是完整http、webSocket连接", params.URL)
		return
	}
	var ok bool
	switch form {
	case FormTypeHTTP:
		// verify
		if params.Verify == "" {
			params.Verify = "statusCode"
		}
		key := fmt.Sprintf("%s.%s", form, params.Verify)
		_, ok = verifyMapHTTP[key]
		if !ok {
			err = errors.New("验证器不存在:" + key)
			return
		}
	case FormTypeWebSocket:
		// verify
		if params.Verify == "" {
			params.Verify = "json"
		}
		key := fmt.Sprintf("%s.%s", form, params.Verify)
		_, ok = verifyMapWebSocket[key]
		if !ok {
			err = errors.New("验证器不存在:" + key)
			return
		}
	}
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	request = &Request{
		URL:       params.URL,
		Form:      form,
		Method:    strings.ToUpper(method),
		Headers:   headers,
		Body:      body,
		Verify:    params.Verify,
		Timeout:   timeout,
		Debug:     debug,
		MaxCon:    maxCon,
		HTTP2:     http2,
		Keepalive: keepalive,
		Code:      params.Code,
		Redirect:  redirect,
	}
	return
}

func getForm(url string) (string, string) {
	form := ""
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		form = FormTypeHTTP
	} else if strings.HasPrefix(url, "ws://") || strings.HasPrefix(url, "wss://") {
		form = FormTypeWebSocket
	} else if strings.HasPrefix(url, "grpc://") || strings.HasPrefix(url, "rpc://") {
		form = FormTypeGRPC
	} else if strings.HasPrefix(url, "radius://") {
		form = FormTypeRadius
		url = url[9:]
	} else {
		form = FormTypeHTTP
		url = fmt.Sprintf("http://%s", url)
	}
	return form, url
}

// getHeaderValue 获取 header
func getHeaderValue(v string, headers map[string]string) {
	index := strings.Index(v, ":")
	if index < 0 {
		return
	}
	vIndex := index + 1
	if len(v) >= vIndex {
		value := strings.TrimPrefix(v[vIndex:], " ")
		if _, ok := headers[v[:index]]; ok {
			headers[v[:index]] = fmt.Sprintf("%s; %s", headers[v[:index]], value)
		} else {
			headers[v[:index]] = value
		}
	}
}

// Print 格式化打印
func (r *Request) Print() {
	if r == nil {
		return
	}
	result := fmt.Sprintf("Request:\n  form:%s \n  url:%s \n  method:%s \n  headers:%v \n", r.Form, r.URL, r.Method,
		r.Headers)
	result = fmt.Sprintf("%s  data:%v \n", result, r.Body)
	result = fmt.Sprintf("%s  verify:%s \n  timeout:%s \n  debug:%v \n", result, r.Verify, r.Timeout, r.Debug)
	result = fmt.Sprintf("%s  http2.0：%v \n  keepalive：%v \n  maxCon:%v ", result, r.HTTP2, r.Keepalive, r.MaxCon)
	color.New(color.FgBlue).Println(result)
	return
}

// GetDebug 获取 debug 参数
func (r *Request) GetDebug() bool {
	return r.Debug
}

// IsParameterLegal 参数是否合法
func (r *Request) IsParameterLegal() (err error) {
	r.Form = "http"
	// statusCode json
	r.Verify = "json"
	key := fmt.Sprintf("%s.%s", r.Form, r.Verify)
	_, ok := verifyMapHTTP[key]
	if !ok {
		return errors.New("验证器不存在:" + key)
	}
	return
}

// RequestResults 请求结果
type RequestResults struct {
	ID            string // 消息ID
	ChanID        uint64 // 消息ID
	Time          uint64 // 请求时间 纳秒
	IsSucceed     bool   // 是否请求成功
	ErrCode       int    // 错误码
	ReceivedBytes int64
}

// SetID 设置请求唯一ID
func (r *RequestResults) SetID(chanID uint64, number uint64) {
	id := fmt.Sprintf("%d_%d", chanID, number)
	r.ID = id
	r.ChanID = chanID
}
