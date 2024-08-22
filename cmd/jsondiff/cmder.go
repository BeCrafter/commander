package jsondiff

import (
	"fmt"
	"os"
	"strings"

	"github.com/BeCrafter/commander/cmd"
	"github.com/BeCrafter/commander/helper"
	"github.com/urfave/cli/v2"
)

// 编译时检查 `(*Cmder)(nil)` 是否满足 `cmd.ICmder` 接口的要求
var _ cmd.ICmder = (*Cmder)(nil)

type Cmder struct{}

func NewCmder() *Cmder {
	return &Cmder{}
}

func (c *Cmder) Register() *cli.Command {
	return &cli.Command{
		Name:  "jsondiff",
		Usage: "Compare two http requests json data by field",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "method",
				Aliases: []string{"X"},
				Usage:   "Request method",
				Value:   "GET",
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "Request header",
			},
			&cli.StringSliceFlag{
				Name:    "host",
				Aliases: []string{"t"},
				Usage:   "Request http host",
			},
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "Request url path",
				Value:   "",
			},
			&cli.StringSliceFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Usage:   `Request data`,
			},
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Diff tools (gojson, cmp)",
				Value: "gojson",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Diff Output Format (ascii, delta)",
				Value:   "ascii",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Debug mode",
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Suppress output, if no differences are found",
			},
		},
	}
}

func (c *Cmder) Action(ctx *cli.Context) error {
	if len(ctx.StringSlice("host")) < 2 {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", helper.ColorSize("至少需要两个请求 or 存在请求失败/空 or 数据不存在可比性", helper.FgRed))
		return cli.ShowSubcommandHelp(ctx)
	}

	req := &helper.Request{
		Method: strings.ToUpper(ctx.String("method")),
		Host:   ctx.StringSlice("host"),
		Url:    ctx.String("url"),
		Header: ctx.StringSlice("header"),
		Data:   ctx.StringSlice("data"),
		Debug:  ctx.Bool("debug"),
		Quiet:  ctx.Bool("quiet"),
	}
	str1, str2 := req.Run()

	var str string
	switch ctx.String("tool") {
	case "cmp":
		str = req.CmpDiff(str1, str2)
	default:
		str = req.JsonDiff(str1, str2, ctx.String("format"))
	}

	fmt.Println(str)
	return nil
}
