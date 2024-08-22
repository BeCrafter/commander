package package_name

import (
	"github.com/BeCrafter/commander/cmd"
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
		Name:  "{CmderName}",
		Usage: "{CmderUsage}",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "{FlagName}",
				Usage: "{FlagUsage}",
			},
		},
	}
}

func (c *Cmder) Action(ctx *cli.Context) error {
	return nil
}
