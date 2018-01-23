package commands

import (
	"github.com/urfave/cli"
)

func (cmds *Commands) newShowCommand() cli.Command {
	return cli.Command{
		Name:      "show",
		Usage:     "show server information",
		ArgsUsage: "<key>",
		Action:    cmds.showAction,
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if c.NArg() > 0 {
				return
			}
			cmds.completeWithServers()
		},
	}
}

func (cmds *Commands) showAction(c *cli.Context) (err error) {
	var key string

	if c.NArg() == 0 {
		key, err = cmds.askServerKey()
		if err != nil {
			return err
		}
	} else {
		key = c.Args().First()
	}

	srv, err := cmds.database.Get(key)
	if err != nil {
		return err
	}

	cmds.printServer(key, srv)

	return nil
}
