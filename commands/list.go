package commands

import (
	"fmt"

	"github.com/mgutz/ansi"
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
		fmt.Println(ansi.Color("Servers list is empty", "red"))
		return nil
	}

	blueColorFunc := ansi.ColorFunc("cyan")
	for key := range hosts {
		fmt.Printf("%s\r\n", blueColorFunc(key))
	}

	return nil
}
