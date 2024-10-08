package listdiff

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/BeCrafter/commander/cmd"
	"github.com/BeCrafter/commander/helper"
	"github.com/fatih/color"
	"github.com/spf13/cast"
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
		Name:  "listdiff",
		Usage: "Compare two http requests json data by list",
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
			&cli.IntFlag{
				Name:  "retry",
				Usage: "Retry times (default: 1)",
				Value: 1,
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
			&cli.BoolFlag{
				Name:  "sort",
				Usage: "Sort the result content to ensure the order of data output",
			},
		},
	}
}

type Data struct {
	List  []map[string]interface{} `json:"list"`
	Total int                      `json:"total"`
}
type CommonRes struct {
	Stat    int    `json:"stat"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    Data   `json:"data"`
	TraceID string `json:"traceId"`
}

func (c *Cmder) Action(ctx *cli.Context) error {
	if len(ctx.StringSlice("host")) < 2 {
		color.New(color.FgRed).Fprintf(os.Stderr, "Error: %v\n\n", "至少需要两个请求 or 存在请求失败/空 or 数据不存在可比性")
		return cli.ShowSubcommandHelp(ctx)
	}

	req := &helper.Request{
		Method: strings.ToUpper(ctx.String("method")),
		Host:   ctx.StringSlice("host"),
		Url:    ctx.String("url"),
		Header: ctx.StringSlice("header"),
		Data:   ctx.StringSlice("data"),
		Retry:  ctx.Int("retry"),
		Debug:  ctx.Bool("debug"),
		Quiet:  ctx.Bool("quiet"),
		Sort:   ctx.Bool("sort"),
	}
	str1, str2 := req.Run()

	var ret1 CommonRes
	json.Unmarshal(str1, &ret1)
	ukeyList1, resList1 := doItem(ret1)

	var ret2 CommonRes
	json.Unmarshal(str2, &ret2)
	ukeyList2, resList2 := doItem(ret2)

	fmt.Printf("\n\n")
	color.New(color.FgHiYellow).Println("============================== ## 下面为Diff数据汇总 ## ==============================")
	fmt.Printf("\n\n")

	// 对比结果个数
	fLen := len(ukeyList1)
	sLen := len(ukeyList2)
	if fLen != sLen {
		color.New(color.FgRed).Printf("# 数据个数存异: 第一个[%v] 第二个[%v]\n", fLen, sLen)
	} else {
		if fLen == 0 {
			color.New(color.FgRed).Printf("# 数据个数为空: 第一个[%v] 第二个[%v]\n", fLen, sLen)
		} else {
			color.New(color.FgGreen).Printf("# 数据个数相同: 第一个[%v] 第二个[%v]\n", fLen, sLen)
		}
	}

	if fLen == 0 || sLen == 0 {
		return nil
	}

	fmt.Printf("\n\n\n")

	// 对比结果顺序一致性
	for k, v := range ukeyList1 {
		if k+1 > sLen {
			color.New(color.FgRed).Printf("<-- 数据顺序存异: Pos[%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), v, helper.StrLeftPad("无", len(v)+1, " "))
		} else {
			if v != ukeyList2[k] {
				color.New(color.FgRed).Printf("<-> 数据顺序存异: Pos[%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), v, ukeyList2[k])
			} else {
				color.New(color.FgGreen).Printf("=== 数据顺序相同: Pos[%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), v, ukeyList2[k])
			}
		}
	}
	if sLen > fLen {
		for k, v := range ukeyList2[fLen:] {
			color.New(color.FgRed).Printf("--> 数据顺序存异: Pos[%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k+fLen, 3, " "), helper.StrLeftPad("无", len(v)+1, " "), v)
		}
	}

	fmt.Printf("\n\n\n")

	// 对比结果内容一致性
	for k, v := range resList1 {
		if v2, has := resList2[k]; has {
			if v != v2 {
				color.New(color.FgYellow).Printf("<-> 数据存异: [%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), v, v2)
			} else {
				color.New(color.FgGreen).Printf("=== 数据相同: [%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), v, v2)
			}
			delete(resList2, k)
		} else {
			color.New(color.FgRed).Printf("<-- 数据存异: [%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), v, helper.StrLeftPad("无", len(v)+1, " "))
		}
	}

	if len(resList2) > 0 {
		fmt.Println("")
		for k, v := range resList2 {
			color.New(color.FgRed).Printf("--> 数据存异: [%v] 第一个[%v] 第二个[%v]\n", helper.StrLeftPad(k, 3, " "), helper.StrLeftPad("无", len(v)+1, " "), v)
		}
	}

	return nil
}

func doItem(res CommonRes) ([]string, map[string]string) {
	ukeyList := []string{}
	resList := map[string]string{}
	for _, v := range res.Data.List {
		var ukey string
		if id, has := v["id"]; has {
			ukey = cast.ToString(id)
		}
		if id, has := v["sku_id"]; has {
			ukey = cast.ToString(id)
		}
		if id, has := v["item_id"]; has {
			ukey = cast.ToString(id)
		}

		ukeyList = append(ukeyList, ukey)

		bytes, _ := json.Marshal(v)
		resList[ukey] = helper.Md5(bytes)
	}
	return ukeyList, resList
}
