package commands

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli"
)

func (cmds *Commands) newListCommand() cli.Command {
	return cli.Command{
		Name:   "list",
		Usage:  "list all servers from database",
		Action: cmds.listAction,
	}
}

func (cmds *Commands) listAction(ctx *cli.Context) error {
	srvs, err := cmds.database.GetAll()
	if err != nil {
		return err
	}

	if len(srvs) == 0 {
		fmt.Println(ansi.Color("Servers list is empty", "red"))
		return nil
	}

	fmt.Println("")
	for key, srv := range srvs {
		cmds.printServer(key, srv)
		fmt.Println("")
	}

	return err
}
