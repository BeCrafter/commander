package main

import (
	"fmt"
	"os"
	"path"

	"github.com/BeCrafter/commander/cmd"
	"github.com/BeCrafter/commander/cmd/jparser"
	"github.com/BeCrafter/commander/cmd/jsondiff"
	"github.com/BeCrafter/commander/cmd/listdiff"
	"github.com/BeCrafter/commander/cmd/stress"
	"github.com/BeCrafter/commander/helper"
	"github.com/urfave/cli/v2"
)

// Commands 所有 Cmder 接口
func RegisterCmder() []cmd.ICmder {
	return []cmd.ICmder{
		jsondiff.NewCmder(),
		listdiff.NewCmder(),
		stress.NewCmder(),
		jparser.NewCmder(),
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Commander"
	app.Authors = []*cli.Author{
		{
			Name:  "kugouming",
			Email: "kugouming@sina.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "link",
			Usage: "Enable subcommand link",
		},
	}

	params := os.Args
	cmdName := path.Base(os.Args[0])
	linkNames := []string{}

	cmders := RegisterCmder()
	for _, cmder := range cmders {
		commond := cmder.Register()
		commond.Action = cmder.Action
		linkNames = append(linkNames, commond.Name)

		if cmdName == commond.Name {
			params = []string{os.Args[0], commond.Name}
			params = append(params, os.Args[1:]...)
		}
		app.Commands = append(app.Commands, commond)
	}

	action := app.Action
	app.Action = func(c *cli.Context) error {
		if c.Bool("link") {
			for _, name := range linkNames {
				if err := helper.LinkCmderBin(name); err != nil {
					return err
				}
			}
			return nil
		}
		return action(c)
	}

	err := app.Run(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "execute failed: %v\n", err)
	}
}
