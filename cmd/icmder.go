package cmd

import "github.com/urfave/cli/v2"

type ICmder interface {
	Register() *cli.Command
	Action(ctx *cli.Context) error
}
