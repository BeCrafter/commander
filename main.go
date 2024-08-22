package main

import (
	"fmt"
	"os"

	"github.com/BeCrafter/commander/cmd"
	"github.com/BeCrafter/commander/cmd/jsondiff"
	"github.com/BeCrafter/commander/cmd/listdiff"
	"github.com/urfave/cli/v2"
)

// Commands 所有 Cmder 接口
func RegisterCmder() []cmd.ICmder {
	return []cmd.ICmder{
		jsondiff.NewCmder(),
		listdiff.NewCmder(),
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

	cmders := RegisterCmder()
	for _, cmder := range cmders {
		commond := cmder.Register()
		commond.Action = cmder.Action
		app.Commands = append(app.Commands, commond)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "execute failed: %v\n", err)
	}
}
