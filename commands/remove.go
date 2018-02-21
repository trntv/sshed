package commands

import (
	"github.com/trntv/sshed/ssh"
	"github.com/urfave/cli"
)

func (cmds *Commands) newRemoveCommand() cli.Command {
	return cli.Command{
		Name:      "remove",
		Usage:     "Removes host",
		ArgsUsage: "<key>",
		Action:    cmds.removeAction,
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if c.NArg() > 0 {
				return
			}
			cmds.completeWithServers()
		},
	}
}

func (cmds *Commands) removeAction(c *cli.Context) (err error) {
	var key string
	if c.NArg() == 0 {
		key, err = cmds.askServerKey()
		if err != nil {
			return err
		}
	} else {
		key = c.Args().First()
	}

	ssh.Config.Remove(key)
	return ssh.Config.Save()
}
