package commands

import (
	"errors"
	"github.com/trntv/sshdb/db"
	"github.com/urfave/cli"
)

func (cmds *Commands) newAtCommand() cli.Command {
	return cli.Command{
		Name:      "at",
		Usage:     "executes command on given server",
		ArgsUsage: "[key] [command]",
		Action:    cmds.atAction,
	}
}
func (cmds *Commands) atAction(c *cli.Context) (err error) {
	var key string
	var srv *db.Server

	if c.NArg() < 2 {
		return errors.New("server and command must be set")
	}

	key = c.Args().First()
	srv, err = cmds.database.Get(key)
	if err != nil {
		return err
	}

	err = cmds.database.Close()
	if err != nil {
		return err
	}

	command := c.Args().Get(1)
	cmds.exec(srv, &options{}, command)

	return err
}
