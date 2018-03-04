package commands

import (
	"github.com/fatih/color"
	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
)

func (cmds *Commands) newListCommand() cli.Command {
	return cli.Command{
		Name:   "list",
		Usage:  "Lists all hosts",
		Action: cmds.listAction,
	}
}

func (cmds *Commands) listAction(ctx *cli.Context) error {
	hosts := ssh.Config.GetAll()
	if len(hosts) == 0 {
		color.New(color.FgRed).Println("Servers list is empty")
		return nil
	}

	cyan := color.New(color.FgCyan).PrintlnFunc()
	for key := range hosts {
		cyan(key)
	}

	return nil
}
