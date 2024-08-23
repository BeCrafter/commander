package jparser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BeCrafter/commander/cmd"
	"github.com/fatih/color"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
)

// 编译时检查 `(*Cmder)(nil)` 是否满足 `cmd.ICmder` 接口的要求
var _ cmd.ICmder = (*Cmder)(nil)

var index int64

type Cmder struct {
	fields    []string
	delimiter string
	debug     bool
	show      bool
}

// 假设日志条目具有如下结构
type LogEntry map[string]interface{}

func NewCmder() *Cmder {
	return &Cmder{}
}

func (c *Cmder) Register() *cli.Command {
	return &cli.Command{
		Name:  "jparser",
		Usage: "JSON parser, accept data input from the pipeline",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "field",
				Aliases:  []string{"f"},
				Usage:    "query field name",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "delimiter",
				Aliases: []string{"d"},
				Usage:   "delimiter string",
				Value:   "\t",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Debug mode",
			},
			&cli.BoolFlag{
				Name:  "show",
				Usage: "structured display of non empty first row data",
			},
		},
	}
}

func (c *Cmder) checker(ctx *cli.Context) error {
	c.fields = ctx.StringSlice("field")
	c.delimiter = ctx.String("delimiter")
	c.debug = ctx.Bool("debug")
	c.show = ctx.Bool("show")

	if len(c.fields) == 0 && !c.show {
		return errors.New("field is not empty")
	}

	return nil
}

func (c *Cmder) Action(ctx *cli.Context) error {
	if err := c.checker(ctx); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				c.processItem(line)
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
			cli.OsExiter(1)
		}
		c.processItem(line)
		index++
	}
	return nil
}

func (c *Cmder) processItem(line string) {
	line = strings.TrimSpace(line)
	if len(line) <= 0 {
		return
	}

	if c.show {
		var entry map[string]interface{}
		json.Unmarshal([]byte(line), &entry)
		body, _ := json.MarshalIndent(entry, "", "    ")
		color.New(color.FgHiCyan).Fprintf(os.Stdout, "\n%s\n\n", string(body))
		cli.OsExiter(0)
	}

	results := gjson.GetMany(line, c.fields...)
	if c.debug {
		fmt.Printf("<<< field_num[%d] \t result_num[%d] \t results[%v]\n", len(results), len(c.fields), results)
	}

	retList := make([]string, 0)
	for k, v := range results {
		if c.debug {
			fmt.Printf(">>> Index[%v] \t Num[%v] \t Raw[%v] \t Str[%v]\n", v.Index, v.Num, v.Raw, v.Str)
		}
		if len(v.Raw) == 0 || v.Str == "-" {
			continue
		}
		retList = append(retList, fmt.Sprintf("%v:%v", c.fields[k], strings.Trim(v.Raw, "\"")))
	}

	if index%2 == 0 {
		color.New(color.FgHiCyan).Fprintf(os.Stdout, "%s\n", strings.Join(retList, c.delimiter))
	} else {
		color.New(color.FgHiGreen).Fprintf(os.Stdout, "%s\n", strings.Join(retList, c.delimiter))
	}
}
