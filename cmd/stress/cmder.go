package stress

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/BeCrafter/commander/cmd"
	"github.com/BeCrafter/commander/cmd/stress/model"
	"github.com/BeCrafter/commander/cmd/stress/server"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// 编译时检查 `(*Cmder)(nil)` 是否满足 `cmd.ICmder` 接口的要求
var _ cmd.ICmder = (*Cmder)(nil)

type Cmder struct {
	params model.Params

	concurrency uint64 // 并发数
	totalNumber uint64 // 请求数(单个并发/协程)
	maxConnect  int    // 单个连接最大请求数
	cpuNumber   int    // CUP 核数，默认为一核，一般场景下单核已经够用了
	timeout     int64  // 接口超时时间
	runtime     int64  // 压测程序最大执行时间，默认不设置
	redirect    bool   // 是否重定向
	debug       bool   // 是否是debug
	http2       bool   // 是否开http2.0
	keepalive   bool   // 是否开启长连接
	path        string // curl文件路径 http接口压测，自定义参数设置
	file        string // 压测源文件
}

func NewCmder() *Cmder {
	return &Cmder{}
}

func (c *Cmder) Register() *cli.Command {
	return &cli.Command{
		Name:  "stress",
		Usage: "Pressure generating tool",
		Flags: []cli.Flag{
			&cli.Uint64Flag{
				Name:    "concurrency",
				Aliases: []string{"c"},
				Usage:   "并发数",
				Value:   1,
			},
			&cli.Uint64Flag{
				Name:    "number",
				Aliases: []string{"n"},
				Usage:   "请求数(单个并发/协程)",
				Value:   1,
			},
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "压测地址",
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "自定义头信息传递给服务器 示例: -H 'Content-Type: application/json'",
			},
			&cli.StringFlag{
				Name:  "data",
				Usage: "HTTP POST 方式传送数据",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "curl文件路径",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "json文件路径",
			},
			&cli.StringFlag{
				Name:  "verify",
				Usage: "验证方法 http 支持:statusCode、json webSocket支持:json",
			},
			&cli.IntFlag{
				Name:  "code",
				Usage: "请求成功的状态码",
				Value: 200,
			},
			&cli.IntFlag{
				Name:  "maxconnect",
				Usage: "单个host最大连接数",
				Value: 1,
			},
			&cli.Int64Flag{
				Name:  "timeout",
				Usage: "接口超时时间 单位 秒，默认为 30s",
				Value: 30,
			},
			&cli.IntFlag{
				Name:  "cpunum",
				Usage: "CUP 核数，默认为一核",
				Value: 1,
			},
			&cli.Int64Flag{
				Name:  "runtime",
				Usage: "压测程序最大执行时间 单位 秒,默认一直压测",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "调试模式",
			},
			&cli.BoolFlag{
				Name:  "redirect",
				Usage: "是否重定向",
			},
			&cli.BoolFlag{
				Name:  "keepalive",
				Usage: "是否开启长连接",
			},
			&cli.BoolFlag{
				Name:  "http2",
				Usage: "是否开 http2.0",
			},
		},
	}
}

func (c *Cmder) checker(ctx *cli.Context) error {
	c.file = ctx.String("file")
	if len(c.file) == 0 {
		c.params = model.Params{
			URL:    ctx.String("url"),
			Header: ctx.StringSlice("header"),
			Body:   ctx.String("data"),
			Verify: ctx.String("verify"),
			Code:   ctx.Int("code"),
		}
	}

	c.path = ctx.String("path")
	c.concurrency = ctx.Uint64("concurrency")
	c.totalNumber = ctx.Uint64("number")
	c.maxConnect = ctx.Int("maxconnect")
	c.cpuNumber = ctx.Int("cpunum")
	c.timeout = ctx.Int64("timeout")
	c.runtime = ctx.Int64("runtime")
	c.debug = ctx.Bool("debug")
	c.redirect = ctx.Bool("redirect")
	c.keepalive = ctx.Bool("keepalive")
	c.http2 = ctx.Bool("http2")

	if c.concurrency == 0 || c.totalNumber == 0 || (c.params.URL == "" && c.path == "" && c.file == "") {
		fmt.Printf("\n示例: \n\n    go run main.go stress -c 1 -n 1 -u https://www.baidu.com/ \n\n")
		fmt.Printf("  1. 压测地址或curl路径必填 \n")
		fmt.Printf("  2. 当前请求参数: -c %d -n %d -d %v -u %s \n\n\n", c.concurrency, c.totalNumber, c.debug, c.params.URL)
		return fmt.Errorf("参数不合法")
	}

	return nil
}

func (c *Cmder) Action(ctx *cli.Context) error {
	if err := c.checker(ctx); err != nil {
		return cli.ShowSubcommandHelp(ctx)
	}

	runtime.GOMAXPROCS(c.cpuNumber)

	color.New(color.FgGreen).Printf("\n开始启动  并发数:%d 请求数:%d 请求参数: \n\n", c.concurrency, c.totalNumber)

	var requests [](*model.Request)
	if c.path != "" { // curl 文件路径
		curls, err := model.ParseTheFile(c.path)
		if err != nil {
			return fmt.Errorf("文件读取失败, Err: %v", err)
		}
		for _, curl := range curls {
			var params model.Params
			params.URL = curl.GetURL()
			for k, v := range curl.GetHeaders() {
				params.Header = append(params.Header, k+":"+v)
			}
			params.Body = curl.GetBody()
			request, err := model.NewRequest(params, time.Duration(c.timeout)*time.Second,
				c.maxConnect, c.http2, c.keepalive, c.redirect, c.debug)
			if err != nil {
				return fmt.Errorf("参数不合法, Err: %v", err)
			}
			requests = append(requests, request)
		}
	} else if c.file != "" { // 自定义json文件路径
		// 此处读取文件，逐行解析并追加到 requests 数组中
		file, err := os.Open(c.file)
		if err != nil {
			return fmt.Errorf("文件读取失败, Err: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 || strings.HasPrefix(line, "#") {
				continue
			}

			var params model.Params
			err := model.ParseParams(line, &params)

			if err != nil {
				return fmt.Errorf("参数不合法, Err: %v", err)
			}
			request, err := model.NewRequest(params, time.Duration(c.timeout)*time.Second,
				c.maxConnect, c.http2, c.keepalive, c.redirect, c.debug)
			if err != nil {
				return fmt.Errorf("参数不合法, Err: %v", err)
			}
			requests = append(requests, request)
		}
	} else { // 命令行参数
		request, err := model.NewRequest(c.params, time.Duration(c.timeout)*time.Second,
			c.maxConnect, c.http2, c.keepalive, c.redirect, c.debug)
		if err != nil {
			return fmt.Errorf("参数不合法, Err: %v", err)
		}
		request.Print()
		requests = append(requests, request)
	}

	// 开始处理
	cctx := context.Background()
	if c.runtime > 0 {
		var cancel context.CancelFunc
		cctx, cancel = context.WithTimeout(cctx, time.Duration(c.runtime)*time.Second)
		defer cancel()
		deadline, ok := ctx.Deadline()
		if ok {
			fmt.Printf(" deadline %s", deadline)
		}
	}

	// 处理 ctrl+c 信号
	cctx, cancelFunc := context.WithCancel(cctx)
	sigChan := make(chan os.Signal, 1) // 使用带缓冲的通道来避免阻塞
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		cancelFunc()
	}()

	if c.debug && len(requests) > 0 {
		color.New(color.FgYellow).Printf("\nDebug 模式开启: \n\n")
		for _, r := range requests {
			r.Print()
			fmt.Printf("\n")
		}
	}
	server.Dispose(cctx, c.concurrency, c.totalNumber, requests)

	return nil
}
